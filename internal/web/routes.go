package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s Server) statusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s Server) configHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.cfg)
	}
}

func (s Server) serveKeyFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.FileAttachment(s.cfg.KeyFile, "ssh_key")
	}
}

func (s Server) serveStartSploit() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := ioutil.ReadFile(s.cfg.StartSploit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		content := string(data)
		content = strings.ReplaceAll(content, "$$SERVER_URL$$", s.cfg.Farm.GetUrl())
		content = strings.ReplaceAll(content, "$$PASSWORD$$", s.cfg.Auth.Password)
		c.Data(http.StatusOK, "application/octet-stream", []byte(content))
	}
}
