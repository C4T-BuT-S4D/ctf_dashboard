package goxy

type ServiceDescription struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Listen string `json:"listen"`
	Target string `json:"target"`
}

type FilterDescription struct {
	ID      int    `json:"id"`
	ProxyID int    `json:"proxy_id"`
	Rule    string `json:"rule"`
	Verdict string `json:"verdict"`
	Enabled bool   `json:"enabled"`
}

type ProxyDescription struct {
	ID                 int                 `json:"id"`
	Service            ServiceDescription  `json:"service"`
	Listening          bool                `json:"listening"`
	FilterDescriptions []FilterDescription `json:"filter_descriptions"`
}
