package web

import (
	"net/http"
	"os"
	"path"

	"ctf_dashboard/internal/common"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewServer(cfg common.Config) *Server {
	ms := &Server{
		cfg:       cfg,
		Router:    gin.New(),
		StaticDir: viper.GetString("web.static_dir"),
	}
	ms.registerMiddleware()
	ms.registerRoutes()
	return ms
}

type Server struct {
	cfg       common.Config
	Router    *gin.Engine
	StaticDir string
}

func (s *Server) registerRoutes() {
	api := s.Router.Group("/api")
	{
		api.GET("/status/", s.statusHandler())
		api.GET("/config/", s.configHandler())
		api.GET("/files/", s.serveFileList())
		api.GET("/file/", s.serveFile())
		api.POST("/add_ssh_key/", s.addSSHKey())
	}

	logrus.Infof("Serving static dir: %s", s.StaticDir)

	s.Router.NoRoute(func(c *gin.Context) {
		realPath := path.Join(s.StaticDir, c.Request.URL.Path)
		if _, err := os.Stat(realPath); os.IsNotExist(err) {
			realPath = path.Join(s.StaticDir, "index.html")
		}
		c.File(realPath)
	})

	logrus.Info("Routes registered successfully")
}

func (s *Server) registerMiddleware() {
	s.Router.Use(gin.Recovery())
	s.Router.Use(loggerMiddleware())

	s.Router.Use(gzip.Gzip(gzip.DefaultCompression))
	s.Router.Use(cors.Default())

	if !gin.IsDebugging() {
		s.Router.Use(gin.BasicAuth(gin.Accounts{
			s.cfg.Auth.Username: s.cfg.Auth.Password,
		}))
	}

	logrus.Info("Middleware registered successfully")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
