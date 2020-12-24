package web

import (
	"ctf_dashboard/internal/deploy"
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

func (s Server) addSSHKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UploadIdRSARequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := deploy.UploadIdRsa(s.cfg.Vulnboxes[req.Vulnbox], req.Key, s.cfg.KeyFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

