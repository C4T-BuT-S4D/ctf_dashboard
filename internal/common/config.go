package common

import (
	"encoding/base64"
	"fmt"
)

type AuthData struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type Mongol struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

type Farm struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

type Neo struct {
	Addr    string `json:"addr" mapstructure:"addr"`
	Version string `json:"version" mapstructure:"version"`
}

type Service struct {
	Name       string `json:"name" mapstructure:"name"`
	Proto      string `json:"proto" mapstructure:"proto"`
	Port       int    `json:"port" mapstructure:"port"`
	HiddenPort int    `json:"hidden" mapstructure:"hidden"`
}

type Vulnbox struct {
	User     string    `json:"user" mapstructure:"user"`
	Host     string    `json:"host" mapstructure:"host"`
	Services []Service `json:"services" mapstructure:"services"`
	GoxyPort int       `json:"goxy_port" mapstructure:"goxy_port"`
}

type Game struct {
	Board string `json:"board" mapstructure:"board"`
	End   string `json:"end" mapstructure:"end"`
}

type NeoRunner struct {
	Path    string `json:"path" mapstructure:"path"`
	Version string `json:"version" mapstructure:"version"`
}

type Config struct {
	Auth          AuthData  `json:"auth" mapstructure:"auth"`
	Game          Game      `json:"game" mapstructure:"game"`
	Mongol        Mongol    `json:"mongol" mapstructure:"mongol"`
	Vulnboxes     []Vulnbox `json:"vulnboxes" mapstructure:"vulnboxes"`
	Farm          Farm      `json:"farm" mapstructure:"farm"`
	Neo           Neo       `json:"neo" mapstructure:"neo"`
	ResourcesPath string    `json:"resources_path" mapstructure:"resources_path"`
	KeyFile       string    `json:"key_file" mapstructure:"key_file"`
}

func (a AuthData) GetHeader() (string, string) {
	data := []byte(fmt.Sprintf("%s:%s", a.Username, a.Password))
	return "Authorization", base64.StdEncoding.EncodeToString(data)
}

func (f Farm) GetUrl() string {
	return fmt.Sprintf("http://%s", f.Addr)
}
