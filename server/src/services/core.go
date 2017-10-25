package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	ini "github.com/vaughan0/go-ini"

	"jbrodriguez/controlr/plugin/server/src/dto"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/model"
	"jbrodriguez/controlr/plugin/server/src/specific"
)

var iniPrefs string = "/boot/config/plugins/dynamix/dynamix.cfg"

// Core service
type Core struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	actor *actor.Actor

	client *http.Client
	state  *model.State

	manager     specific.Manager
	logLocation map[string]string

	info    dto.Info
	watcher *fsnotify.Watcher

	ups model.Ups
}

// NewCore - constructor
func NewCore(bus *pubsub.PubSub, settings *lib.Settings, state *model.State) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
		state:    state,
		manager:  specific.NewManager(state.Version),
		logLocation: map[string]string{
			"system": "/var/log/syslog",
			"docker": "/var/log/docker.log",
			"vm":     "/var/log/libvirt/libvirtd.log",
		},
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	core.client = &http.Client{Timeout: time.Second * 10, Transport: tr}

	return core
}

// Start service
func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	c.actor.Register("model/REFRESH", c.refresh)
	c.actor.Register("model/UPDATE_USER", c.updateUser)
	c.actor.Register("api/GET_LOG", c.getLog)
	c.actor.Register("api/GET_INFO", c.getInfo)
	c.actor.Register("api/GET_MAC", c.getMac)
	c.actor.Register("api/GET_PREFS", c.getPrefs)

	ups, err := model.IdentifyUps()
	if err != nil {
		mlog.Warning("Error identifying UPS: %s", err)
		c.ups = model.NewNoUps()
	} else {
		switch ups {
		case model.APC:
			c.ups = model.NewApc()
		case model.NUT:
			c.ups = model.NewNut()
		default:
			c.ups = model.NewNoUps()
			break
		}
	}

	wake := _getMac()
	prefs, err := _getPrefs()
	if err != nil {
		mlog.Warning("Unable to load/parse prefs file (%s): %s", iniPrefs, err)
	}
	samples := c.ups.GetStatus()

	c.info = dto.Info{
		Version: 1,
		Wake:    wake,
		Prefs:   prefs,
		Samples: samples,
	}

	c.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		mlog.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-c.watcher.Events:
				mlog.Info("event: %s", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					mlog.Info("modified file: %s", event.Name)

					prefs, err := _getPrefs()
					if err != nil {
						mlog.Warning("Unable to load/parse prefs file (%s): %s", iniPrefs, err)
					}

					c.info.Prefs = prefs
				}
			case err := <-c.watcher.Errors:
				mlog.Warning("Error:", err)
			}
		}
	}()

	err = c.watcher.Add(iniPrefs)
	if err != nil {
		mlog.Fatal(err)
	}

	go c.actor.React()

	return nil
}

// Stop service
func (c *Core) Stop() {
	if c.watcher != nil {
		c.watcher.Close()
	}

	mlog.Info("stopped service Core ...")
}

// PLUGIN APP HANDLERS
func (c *Core) refresh(msg *pubsub.Message) {
	go func() {
		_dockers, err := lib.Get(c.client, c.state.Host, "/Docker")
		if err != nil {
			mlog.Warning("Unable to get dockers: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (dockers): %s", err)}
			c.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
			return
		}

		_vms, err := lib.Get(c.client, c.state.Host, "/VMs")
		if err != nil {
			mlog.Warning("Unable to get vms: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (vms): %s", err)}
			c.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
			return
		}

		_users, err := lib.Get(c.client, c.state.Host, "/state/users.ini")
		if err != nil {
			mlog.Warning("Unable to get users.ini: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (users): %s", err)}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			return
		}

		mlog.Info("Getting users ...")
		users := c.manager.GetUsers(_users)
		mlog.Info("Got %d users", len(users))
		mlog.Info("Getting apps ...")
		apps := c.manager.GetApps(_dockers, _vms)
		mlog.Info("Got %d apps", len(apps))

		outbound := &dto.Packet{Topic: "model/REFRESHED", Payload: map[string]interface{}{"users": users, "apps": apps}}
		c.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
	}()
}

func (c *Core) updateUser(msg *pubsub.Message) {
	args := msg.Payload.(map[string]interface{})

	data := map[string]string{
		"userName":    args["name"].(string),
		"userDesc":    args["perms"].(string),
		"cmdUserEdit": "Apply",
	}
	if c.state.CsrfToken != "" {
		data["csrf_token"] = c.state.CsrfToken
	}

	_, err := lib.Post(c.client, c.state.Host, "/update.htm", data)
	if err != nil {
		mlog.Warning("Unable to post updateUser: %s", err)
		outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: "Unable to update User"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	outbound := &dto.Packet{Topic: "model/USER_UPDATED", Payload: map[string]interface{}{"status": "ok"}}
	c.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
}

// API HANDLERS
func (c *Core) getLog(msg *pubsub.Message) {
	logType := msg.Payload.(string)

	log := make([]string, 0)

	exists, err := lib.Exists(c.logLocation[logType])
	if err != nil {
		mlog.Warning("Unable to check for log existence: %s", err)
		msg.Reply <- log
		return
	}

	if !exists {
		mlog.Warning("Log %s is not present in the system", logType)
		msg.Reply <- log
		return
	}

	cmd := "tail -n 40 " + c.logLocation[logType]

	lib.Shell(cmd, func(line string) {
		log = append(log, line)
	})

	msg.Reply <- log
}

func (c *Core) getInfo(msg *pubsub.Message) {
	c.info.Samples = c.ups.GetStatus()
	msg.Reply <- c.info
}

func (c *Core) getMac(msg *pubsub.Message) {
	wake := _getMac()

	msg.Reply <- wake.Mac
}

func (c *Core) getPrefs(msg *pubsub.Message) {
	prefs, err := _getPrefs()
	if err != nil {
		mlog.Warning("Unable to load/parse prefs file (%s): %s", iniPrefs, err)
	}

	msg.Reply <- prefs
}

func _getMac() dto.Wake {
	wake := dto.Wake{
		Mac:       "",
		Broadcast: "255.255.255.255",
	}

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// mlog.Info("[%s] = %s", iface.Name, iface.HardwareAddr)
		if iface.Name == "eth0" {
			wake.Mac = fmt.Sprintf("%s", iface.HardwareAddr)
			break
		}
	}

	return wake
}

func _getPrefs() (dto.Prefs, error) {
	prefs := dto.Prefs{
		Number: ".,",
		Unit:   "C",
	}

	file, err := ini.LoadFile(iniPrefs)
	if err != nil {
		return prefs, err
	}

	for key, value := range file["display"] {
		if key == "number" {
			prefs.Number = strings.Replace(value, "\"", "", -1)
		}

		if key == "unit" {
			prefs.Unit = strings.Replace(value, "\"", "", -1)
		}
	}

	return prefs, nil
}
