package services

import (
	"net/http"
	"path/filepath"

	"jbrodriguez/controlr/plugin/server/src/dto"
	"jbrodriguez/controlr/plugin/server/src/lib"
	"jbrodriguez/controlr/plugin/server/src/model"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const (
	apiVersion = "/api/v1"
	capacity   = 3
)

// Api type
type Api struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	engine *echo.Echo

	state  *model.State
	secret string
}

// NewApi - constructor
func NewApi(bus *pubsub.PubSub, settings *lib.Settings, state *model.State) *Api {
	server := &Api{
		bus:      bus,
		settings: settings,
		state:    state,
	}
	return server
}

const basic = "Basic"

// Start service
func (a *Api) Start() {
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
					// mlog.Info("auth: %s", auth)
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

	if a.state.Secure {
		go func() {
			err := a.engine.StartTLS(a.settings.ApiPort, filepath.Join(a.settings.CertDir, "certificate_bundle.pem"), filepath.Join(a.settings.CertDir, "certificate_bundle.pem"))
			if err != nil {
				mlog.Fatalf("Unable to start https api: %s", err)
			}
		}()

		mlog.Info("Api started listening https on %s", a.settings.ApiPort)
	} else {
		go func() {
			err := a.engine.Start(a.settings.ApiPort)
			if err != nil {
				mlog.Fatalf("Unable to start http api: %s", err)
			}
		}()

		mlog.Info("Api started listening http on %s", a.settings.ApiPort)
	}
}

// Stop service
func (a *Api) Stop() {
	mlog.Info("stopped service Api ...")
}

func (a *Api) getLog(c echo.Context) (err error) {
	logType := c.Param("logType")
	mlog.Info("log (%s) requested", logType)

	msg := &pubsub.Message{Payload: logType, Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_LOG")

	reply := <-msg.Reply
	resp := reply.([]string)

	return c.JSON(http.StatusOK, &resp)
}

func (a *Api) debugGet(c echo.Context) (err error) {
	mlog.Info("received debug/get")
	return c.String(http.StatusOK, "Ok")
}

func (a *Api) debugPost(c echo.Context) (err error) {
	req, _ := c.FormParams()
	mlog.Info("req=%+v", req)

	return c.String(http.StatusOK, "Ok")
}

func (a *Api) getInfo(c echo.Context) (err error) {
	mlog.Info("received /info")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_INFO")

	reply := <-msg.Reply
	resp := reply.(dto.Info)

	mlog.Info("info(%+v)", resp)

	return c.JSON(http.StatusOK, &resp)
}

func (a *Api) getMac(c echo.Context) (err error) {
	mlog.Info("received /mac")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_MAC")

	reply := <-msg.Reply
	resp := reply.(string)
	// c.JSON(200, &resp)

	return c.JSON(http.StatusOK, &resp)
}

func (a *Api) getPrefs(c echo.Context) (err error) {
	mlog.Info("received /prefs")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	a.bus.Pub(msg, "api/GET_PREFS")

	reply := <-msg.Reply
	resp := reply.(dto.Prefs)
	// c.JSON(200, &resp)

	return c.JSON(http.StatusOK, &resp)
}
