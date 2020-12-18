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

type Service struct {
	Name       string `json:"name" mapstructure:"name"`
	Proto      string `json:"proto" mapstructure:"proto"`
	Port       int    `json:"port" mapstructure:"port"`
	HiddenPort int    `json:"hidden" mapstructure:"hidden"`
}

type Vulnbox struct {
	Host     string   `json:"host" mapstructure:"host"`
	Services []string `json:"services" mapstructure:"services"`
	GoxyPort int      `json:"goxy_port" mapstructure:"goxy_port"`
}

type Config struct {
	Auth        AuthData  `json:"auth" mapstructure:"auth"`
	Mongol      Mongol    `json:"mongol" mapstructure:"mongol"`
	Services    []Service `json:"services" mapstructure:"services"`
	Vulnboxes   []Vulnbox `json:"vulnboxes" mapstructure:"vulnboxes"`
	Farm        Farm      `json:"farm" mapstructure:"farm"`
	StartSploit string    `json:"start_sploit" mapstructure:"start_sploit"`
	KeyFile     string    `json:"key_file" mapstructure:"key_file"`
}

func (a AuthData) GetHeader() (string, string) {
	data := []byte(fmt.Sprintf("%s:%s", a.Username, a.Password))
	return "Authorization", base64.StdEncoding.EncodeToString(data)
}

func (f Farm) GetUrl() string {
	return fmt.Sprintf("http://%s", f.Addr)
}
