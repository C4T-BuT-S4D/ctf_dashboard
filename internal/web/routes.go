package web

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"ctf_dashboard/internal/deploy"

	"github.com/gin-gonic/gin"
)

func (s *Server) statusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Server) configHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.cfg)
	}
}

func (s *Server) serveFileList() gin.HandlerFunc {
	type response struct {
		Files []string `json:"files"`
	}

	return func(c *gin.Context) {
		entries, err := os.ReadDir(s.cfg.ResourcesPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp := &response{Files: make([]string, 0, len(entries))}
		for _, entry := range entries {
			if !entry.IsDir() {
				resp.Files = append(resp.Files, entry.Name())
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) serveFile() gin.HandlerFunc {
	type request struct {
		Name string `form:"name" binding:"required"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		path := filepath.Join(s.cfg.ResourcesPath, req.Name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Data(
			http.StatusOK,
			"application/octet-stream",
			s.replaceVariables(data),
		)
	}
}

func (s *Server) addSSHKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UploadIdRSARequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for _, vulnbox := range s.cfg.Vulnboxes {
			if err := deploy.UploadSSHKey(vulnbox, req.Key, s.cfg.KeyFile); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Server) replaceVariables(data []byte) []byte {
	replacements := map[string]string{
		"NEO_VERSION": s.cfg.Neo.Version,
		"NEO_ADDR":    s.cfg.Neo.Addr,
		"FARM_ADDR":   s.cfg.Farm.Addr,
		"PROXY_ADDR":  s.cfg.Proxy.Addr,
		"PASSWORD":    s.cfg.Auth.Password,
	}
	for k, v := range replacements {
		data = bytes.ReplaceAll(data, []byte(fmt.Sprintf("$$%s$$", k)), []byte(v))
	}
	return data
}
