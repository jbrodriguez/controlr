package services

import (
	// "encoding/base64"
	"errors"
	"fmt"
	// "github.com/ddliu/go-httpclient"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "github.com/vaughan0/go-ini"
	"io/ioutil"
	"jbrodriguez/controlr/plugin/server/src/dto"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/specific"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// Unraid service
type Unraid struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings

	mailbox chan *pubsub.Mailbox

	client *http.Client
	data   map[string]string

	manager specific.Manager
}

// NewUnraid - constructor
func NewUnraid(bus *pubsub.PubSub, settings *lib.Settings, data map[string]string) *Unraid {
	core := &Unraid{
		bus:      bus,
		settings: settings,
		client:   &http.Client{Timeout: time.Second * 10},
		data:     data,
	}

	core.manager = specific.NewManager(data["version"])

	core.init()

	return core
}

// Start service
func (u *Unraid) Start() (err error) {
	mlog.Info("starting service Unraid ...")

	u.mailbox = u.register(u.bus, "model/REFRESH", u.refresh)
	u.registerAdditional(u.bus, "model/UPDATE_USER", u.updateUser, u.mailbox)

	go u.react()

	return nil
}

// Stop service
func (u *Unraid) Stop() {
	mlog.Info("stopped service Unraid ...")
}

func (u *Unraid) react() {
	for mbox := range u.mailbox {
		u.dispatch(mbox.Topic, mbox.Content)
	}
}

func (u *Unraid) refresh(msg *pubsub.Message) {
	go func() {
		_dockers, err := u.get("/Docker")
		if err != nil {
			mlog.Warning("Unable to get dockers: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (dockers): %s", err)}
			u.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
			return
		}

		_vms, err := u.get("/VMs")
		if err != nil {
			mlog.Warning("Unable to get vms: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (vms): %s", err)}
			u.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
			return
		}

		_users, err := u.get("/state/users.ini")
		if err != nil {
			mlog.Warning("Unable to get users.ini: %s", err)
			outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: fmt.Sprintf("Unable to get unRAID state (users): %s", err)}
			u.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			return
		}

		users := u.manager.GetUsers(_users)
		apps := u.manager.GetApps(_dockers, _vms)

		outbound := &dto.Packet{Topic: "model/REFRESHED", Payload: map[string]interface{}{"users": users, "apps": apps}}
		u.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
	}()
}

func (u *Unraid) updateUser(msg *pubsub.Message) {
	args := msg.Payload.(map[string]interface{})

	data := map[string]string{
		"userName":    args["name"].(string),
		"userDesc":    args["perms"].(string),
		"cmdUserEdit": "Apply",
	}

	_, err := u.post("/update.htm", data)
	if err != nil {
		mlog.Warning("Unable to post updateUser: %s", err)
		outbound := &dto.Packet{Topic: "model/ACCESS_ERROR", Payload: "Unable to update User"}
		u.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	outbound := &dto.Packet{Topic: "model/USER_UPDATED", Payload: map[string]interface{}{"status": "ok"}}
	u.bus.Pub(&pubsub.Message{Id: msg.Id, Payload: outbound}, "socket:broadcast")
}

func (u *Unraid) get(resource string) (string, error) {
	ep, _ := url.Parse(u.data["backend"])
	ep.Path = path.Join(ep.Path, resource)

	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

func (u *Unraid) post(resource string, args map[string]string) (string, error) {
	ep, _ := url.Parse(u.data["backend"])
	ep.Path = path.Join(ep.Path, resource)

	data := url.Values{}
	for k, v := range args {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", ep.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(resp.Status), nil
}
