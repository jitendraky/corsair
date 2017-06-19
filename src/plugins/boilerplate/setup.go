package boilerplate

import (
	"net/http"
  "fmt"

	"corsair/corsair"
	"corsair/network/http/httpserver"
)

// ## Boilerplate
// Boilerplate provides an example plugin skeleton
// that can eventually be generated for rapid custom
// plugin development.
type httpHandler struct {
    Next httpserver.Handler
    *Config
}

type Config struct {
	Path           string
	Endpoint       string
	//Template       *template.Template
	SiteRoot       string
}

func init() {
	corsair.RegisterPlugin("boilerplate", corsair.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

func Setup(c *corsair.Controller) (err error) {
  cfg := httpserver.GetConfig(c)
  var conf *Config

  conf, err = ParseBoilerplateConfig(c, cfg)
  if err != nil {
    return err
  }

  cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
    return &httpHandler{
      Next: next,
      Config: conf,
    }
  })
  return
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	return h.Next.ServeHTTP(w, r)
}

func ParseBoilerplateConfig(c *corsair.Controller, cnf *httpserver.SiteConfig) (*Config, error) {
	config := &Config{
		Endpoint:       `/search`,
		SiteRoot:       cnf.Root,
		//Template:       nil,
	}
  for c.NextBlock() {
		switch c.Val() {
		case "endpoint":
      fmt.Println("endpoint configure")
		}
  }
  return config, nil
}
