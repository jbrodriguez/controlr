package services

import (
	"net/http"
	"path/filepath"

	"plugin/dto"
	"plugin/lib"
	"plugin/model"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const (
	apiVersion = "/api/v1"
	capacity   = 3
)

// API type
type API struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	engine *echo.Echo

	state *model.State
}

// NewAPI - constructor
func NewAPI(bus *pubsub.PubSub, settings *lib.Settings, state *model.State) *API {
	server := &API{
		bus:      bus,
		settings: settings,
		state:    state,
	}
	return server
}

const basic = "Basic"

// Start service
func (a *API) Start() {
	mlog.Info("Starting service Api ...")

	a.engine = echo.New()

	a.engine.HideBanner = true

	a.engine.Use(mw.Logger())
	a.engine.Use(mw.Recover())

	r := a.engine.Group(apiVersion)
	r.Use(mw.BasicAuthWithConfig(mw.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && auth[:l] == basic {
				if auth[l+1:] == "null" {
					return true
				}
			}

			return false
		},
		Validator: func(usr, pwd string, c echo.Context) (bool, error) {
			return true, nil
		},
	}))

	r.GET("/log/:logType", a.getLog)
	r.GET("/debug", a.debugGet)
	r.POST("/debug", a.debugPost)
	r.GET("/info", a.getInfo)
	r.GET("/mac", a.getMac)
	r.GET("/prefs", a.getPrefs)

	go func() {
		var err error
		if a.state.UseSelfCerts {
			err = a.engine.StartTLS(a.settings.APIPort, filepath.Join(a.settings.CertDir, "controlr_cert.pem"), filepath.Join(a.settings.CertDir, "controlr_key.pem"))
		} else {
			err = a.engine.StartTLS(a.settings.APIPort, filepath.Join(a.settings.CertDir, a.state.Cert), filepath.Join(a.settings.CertDir, a.state.Cert))
		}
		if err != nil {
			mlog.Fatalf("Unable to start https api: %s", err)
		}
	}()

	mlog.Info("Api started listening https on %s", a.settings.APIPort)
}

// Stop service
func (a *API) Stop() {
	mlog.Info("stopped service Api ...")
}

func (a *API) getLog(c echo.Context) (err error) {
	logType := c.Param("logType")
	mlog.Info("log (%s) requested", logType)

	msg := &pubsub.Message{Payload: logType, Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_LOG")

	reply := <-msg.Reply
	resp := reply.([]string)

	return c.JSON(http.StatusOK, &resp)
}

func (a *API) debugGet(c echo.Context) (err error) {
	mlog.Info("received debug/get")
	return c.String(http.StatusOK, "Ok")
}

func (a *API) debugPost(c echo.Context) (err error) {
	req, _ := c.FormParams()
	mlog.Info("req=%+v", req)

	return c.String(http.StatusOK, "Ok")
}

func (a *API) getInfo(c echo.Context) (err error) {
	mlog.Info("received /info")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_INFO")

	reply := <-msg.Reply
	resp := reply.(dto.Info)

	mlog.Info("info(%+v)", resp)

	return c.JSON(http.StatusOK, &resp)
}

func (a *API) getMac(c echo.Context) (err error) {
	mlog.Info("received /mac")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_MAC")

	reply := <-msg.Reply
	resp := reply.(string)

	return c.JSON(http.StatusOK, &resp)
}

func (a *API) getPrefs(c echo.Context) (err error) {
	mlog.Info("received /prefs")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_PREFS")

	reply := <-msg.Reply
	resp := reply.(dto.Prefs)

	return c.JSON(http.StatusOK, &resp)
}
