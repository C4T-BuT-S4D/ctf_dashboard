package web

type UploadIdRSARequest struct {
	Vulnbox int    `json:"vulnbox"`
	Key     string `json:"key"`
}
