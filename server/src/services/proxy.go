package services

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	// "jbrodriguez/controlr/plugin/server/src/dto"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"net/http"
)

const (
	proxyVersion  = "/api/v1"
	proxyCapacity = 3
)

// Proxy type
type Proxy struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	engine *echo.Echo

	data   map[string]string
	secret string
}

// NewProxy - constructor
func NewProxy(bus *pubsub.PubSub, settings *lib.Settings, data map[string]string) *Proxy {
	server := &Proxy{
		bus:      bus,
		settings: settings,
		data:     data,
	}
	// server.init()
	return server
}

const basic = "Basic"

// Start service
func (p *Proxy) Start() {
	mlog.Info("Starting service Proxy ...")

	// h := sha256.Sum256([]byte(s.data["gateway"] + s.data["name"] + s.data["timezone"] + s.data["version"]))
	// targetURL, _ := url.Parse(s.data["backend"])

	p.engine = echo.New()

	p.engine.Use(mw.Logger())
	p.engine.Use(mw.Recover())

	// p.engine.Static("/", filepath.Join(location, "index.html"))

	r := p.engine.Group(proxyVersion)
	r.Use(mw.BasicAuthWithConfig(mw.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {
			auth := c.Request().Header().Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && auth[:l] == basic {
				if auth[l+1:] == "null" {
					// mlog.Info("auth: %s", auth)
					return true
				}
			}

			return false
		},
		Validator: func(usr, pwd string) bool {
			// mlog.Info("auth:usr:%s", usr)
			return true
		},
	}))
	r.Get("/log/:logType", p.getLog)
	r.Get("/debug", p.debugGet)
	r.Post("/debug", p.debugPost)

	port := fmt.Sprintf(":%s", p.settings.ProxyPort)
	go p.engine.Run(standard.New(port))

	mlog.Info("Proxy started listening on %s", port)
}

// Stop service
func (p *Proxy) Stop() {
	mlog.Info("stopped service Proxy ...")
}

func (p *Proxy) getLog(c echo.Context) (err error) {
	logType := c.Param("logType")
	mlog.Info("log (%s) requested", logType)

	msg := &pubsub.Message{Payload: logType, Reply: make(chan interface{}, capacity)}
	p.bus.Pub(msg, "api/GET_LOG")

	reply := <-msg.Reply
	resp := reply.([]string)
	// c.JSON(200, &resp)

	return c.JSON(http.StatusOK, &resp)
}

func (p *Proxy) debugGet(c echo.Context) (err error) {
	mlog.Info("received debug/get")
	return c.String(http.StatusOK, "Ok")
}

func (p *Proxy) debugPost(c echo.Context) (err error) {
	req := c.FormParams()
	mlog.Info("req=%+v", req)

	return c.String(http.StatusOK, "Ok")
}
