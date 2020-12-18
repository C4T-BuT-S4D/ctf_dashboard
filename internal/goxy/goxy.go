package goxy

import (
	"ctf_dashboard/internal/common"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func New(host common.Vulnbox, auth common.AuthData) *Goxy {
	return &Goxy{
		host: host,
		auth: auth,
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}
}

type Goxy struct {
	host   common.Vulnbox
	auth   common.AuthData
	client http.Client
}

func (g *Goxy) GetProxies() (*ProxyDescription, error) {
	url := g.getApiUrl("proxies")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	hName, hVal := g.auth.GetHeader()
	req.Header.Add(hName, hVal)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("Error closing response body: %v", err)
		}
	}()
	dec := json.NewDecoder(resp.Body)
	result := new(ProxyDescription)
	if err := dec.Decode(result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}

func (g *Goxy) getApiUrl(suffix string) string {
	return fmt.Sprintf("http://%s:%d/api/%s/", g.host.Host, g.host.GoxyPort, suffix)
}
