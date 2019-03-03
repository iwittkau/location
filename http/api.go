package http

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/iwittkau/location"

	"github.com/gin-gonic/gin"
)

const (
	defaultAddress       = ":8080"
	defaultStaticPath    = "web/static"
	defaultTemplatesPath = "web/templates"
)

var (
	DefaultOpts = Options{
		Address:       defaultAddress,
		StaticPath:    defaultStaticPath,
		TemplatesPath: defaultTemplatesPath,
		Debug:         false,
	}
)

type Options struct {
	Address       string
	StaticPath    string
	TemplatesPath string
	Debug         bool
	SecretHash    string
}

type API struct {
	opts  Options
	r     *gin.Engine
	store location.Storage
}

func New(opts Options, store location.Storage) (*API, error) {
	if _, err := bcrypt.Cost([]byte(opts.SecretHash)); err != nil {
		return nil, err
	}
	return &API{
		opts: opts,
	}, nil
}

func (a *API) Open() error {
	a.setupRouter()
	return a.r.Run(a.opts.Address)
}

func (a *API) setupRouter() {
	if !a.opts.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	a.r = gin.Default()
	a.r.LoadHTMLGlob(a.opts.TemplatesPath + "/*.html")

	a.r.Static("/static", a.opts.StaticPath)
	a.r.GET("/", a.handleGetIndex)
	a.r.GET("/checkin", a.handleGetLocation)

}

func (a *API) handleGetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (a *API) handleGetLocation(c *gin.Context) {
	c.HTML(http.StatusOK, "location.html", nil)
}
