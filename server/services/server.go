package services

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"plugin/dto"
	"plugin/lib"
	"plugin/model"
	"plugin/ntk"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/tredoe/osutil/user/crypt"
	"github.com/tredoe/osutil/user/crypt/md5_crypt"
	"github.com/tredoe/osutil/user/crypt/sha256_crypt"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"

	"golang.org/x/net/websocket"
)

// const (
// 	apiVersion = "/api/v1"
// 	capacity   = 3
// )

// Server type
type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	engine *echo.Echo
	actor  *actor.Actor

	pool   map[uint64]*ntk.Connection
	state  *model.State
	secret string

	proxy *httputil.ReverseProxy
}

// NewServer - constructor
func NewServer(bus *pubsub.PubSub, settings *lib.Settings, state *model.State) *Server {
	server := &Server{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
		pool:     make(map[uint64]*ntk.Connection),
		state:    state,
	}
	return server
}

func redirector(sPort string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req, scheme := c.Request(), c.Scheme()
			host, _, err := net.SplitHostPort(req.Host)
			if err != nil {
				log.Printf("err(%s)", err)
				return next(c)
			}

			if scheme != "https" {
				return c.Redirect(http.StatusMovedPermanently, "https://"+host+sPort+req.RequestURI)
			}

			return next(c)
		}
	}
}

// Start service
func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	cwd, _ := os.Getwd()

	locations := []string{
		"/usr/local/emhttp/plugins/controlr",
		cwd,
		s.settings.WebDir,
	}

	location := lib.SearchFile("index.html", locations)
	if location == "" {
		msg := ""
		for _, loc := range locations {
			msg += fmt.Sprintf("%s, ", loc)
		}
		mlog.Fatalf("Unable to find index.html. Exiting now. (searched in %s)", msg)
	}

	mlog.Info("Serving files from %s", location)

	// create JWT secret
	h := sha256.Sum256([]byte(s.state.Name + s.state.Timezone + s.state.Version + s.state.CsrfToken))
	s.secret = base64.StdEncoding.EncodeToString(h[:])

	// port for https is port for http + 1
	var iPort int
	var err error
	if iPort, err = strconv.Atoi(s.settings.Port[:1]); err != nil {
		iPort = 2378
	}
	sPort := fmt.Sprintf(":%d", iPort+1)

	s.engine = echo.New()

	s.engine.HideBanner = true

	s.engine.Use(mw.Logger())
	s.engine.Use(mw.Recover())
	s.engine.Use(mw.CORS())
	if s.state.Secure {
		s.engine.Use(redirector(sPort))
	}

	s.engine.Static("/", filepath.Join(location, "index.html"))
	s.engine.Static("/favicon.ico", filepath.Join(location, "app", "favicon.ico"))
	s.engine.Static("/img", filepath.Join(location, "app", "img"))
	s.engine.Static("/js", filepath.Join(location, "app", "js"))
	s.engine.Static("/css", filepath.Join(location, "app", "css"))

	s.engine.GET("/version", s.getVersion)
	s.engine.POST("/login", s.login)

	s.engine.GET("/state/plugins/*", s.proxyHandler)
	s.engine.GET("/plugins/*", s.proxyHandler)

	r := s.engine.Group("/skt")
	r.Use(mw.JWTWithConfig(mw.JWTConfig{
		SigningKey:  []byte(s.secret),
		TokenLookup: "query:token",
	}))
	r.GET("/", s.handleWs)

	s.actor.Register("socket:broadcast", s.broadcast)
	go s.actor.React()

	targetURL, _ := url.Parse(s.state.Host)

	// Always listen on http port, but based on above setting, we could be redirecting to https
	go func() {
		err := s.engine.Start(s.settings.Port)
		if err != nil {
			mlog.Fatalf("Unable to start http server: %s", err)
		}
	}()

	mlog.Info("Server started listening http on %s", s.settings.Port)

	if s.state.Secure {
		s.proxy = CreateReverseProxy(targetURL)

		go func() {
			err := s.engine.StartTLS(sPort, filepath.Join(s.settings.CertDir, s.state.Cert), filepath.Join(s.settings.CertDir, s.state.Cert))
			if err != nil {
				mlog.Fatalf("Unable to start https server: %s", err)
			}
		}()

		mlog.Info("Server started listening https on %s", sPort)
	} else {
		s.proxy = httputil.NewSingleHostReverseProxy(targetURL)
	}
}

// Stop service
func (s *Server) Stop() {
	mlog.Info("stopped service Server ...")
}

func (s *Server) getVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"version": s.settings.Version})
}

func (s *Server) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "root" {
		mlog.Warning("Must log in as root")
		return c.JSON(http.StatusUnauthorized, map[string]string{"token": "invalid"})
	}

	if !s.settings.Dev {
		// get the /etc/shadow entry for root
		shadowLine := ""
		lib.Shell("getent shadow root", func(line string) {
			shadowLine = line
		})

		re := regexp.MustCompile(`root:(\$(.*?)\$(.*?)\$.*?):`)

		saltString := ""
		actualHash := ""
		encType := ""
		for _, match := range re.FindAllStringSubmatch(shadowLine, -1) {
			actualHash = match[1]
			encType = match[2]
			saltString = match[3]
		}

		var crypto crypt.Crypter
		saltPrefix := ""
		switch encType {
		case "1":
			crypto = crypt.New(crypt.MD5)
			saltPrefix = md5_crypt.MagicPrefix
		case "5":
			crypto = crypt.New(crypt.SHA256)
			saltPrefix = sha256_crypt.MagicPrefix
		case "6":
			crypto = crypt.New(crypt.SHA512)
			saltPrefix = sha512_crypt.MagicPrefix
		default:
			mlog.Warning("Unknown encryption type: (%s)", encType)
			return c.JSON(http.StatusUnauthorized, map[string]string{"token": "invalid"})
		}

		saltString = fmt.Sprintf("%s%s", saltPrefix, saltString)

		shadowHash, err := crypto.Generate([]byte(password), []byte(saltString))
		if err != nil {
			mlog.Warning("Unable to create hash: %s", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"token": "invalid"})
		}

		if shadowHash != actualHash {
			mlog.Warning("shadowHash != actualHash")
			return c.JSON(http.StatusUnauthorized, map[string]string{"token": "invalid"})
		}
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = now.Unix()
	claims["name"] = username
	claims["admin"] = true
	claims["exp"] = now.Add(time.Minute * 60).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"token": fmt.Sprintf("%s", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}

// WEBSOCKET handler
func (s *Server) handleWs(c echo.Context) (err error) {
	websocket.Handler(func(ws *websocket.Conn) {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := uint64(claims["id"].(float64))

		conn := ntk.NewConnection(id, ws, s.onMessage, s.onClose)
		s.pool[id] = conn
		if err := conn.Read(); err != nil {
			mlog.Warning("error reading from connection: %s", err)
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func (s *Server) onMessage(packet *dto.Packet) {
	s.bus.Pub(&pubsub.Message{Id: packet.ID, Payload: packet.Payload}, packet.Topic)
}

func (s *Server) onClose(c *ntk.Connection, _ error) {
	delete(s.pool, c.ID)
}

func (s *Server) broadcast(msg *pubsub.Message) {
	packet := msg.Payload.(*dto.Packet)
	if _, ok := s.pool[msg.Id]; ok {
		conn := s.pool[msg.Id]
		if err := conn.Write(packet); err != nil {
			mlog.Warning("error writing to connection: %s", err)
		}
	}
}

// PROXY for images
func (s *Server) proxyHandler(c echo.Context) (err error) {
	s.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// CreateReverseProxy -
func CreateReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}
