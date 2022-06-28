package types

import (
	"time"

	"github.com/z2665/git-lfs-cos-agent/pkg/config"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Action struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}
type OperationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Combined request struct which can accept anything
type Request struct {
	Event               string  `json:"event"`
	Operation           string  `json:"operation"`
	Concurrent          bool    `json:"concurrent"`
	ConcurrentTransfers int     `json:"concurrenttransfers"`
	Oid                 string  `json:"oid"`
	Size                int64   `json:"size"`
	Path                string  `json:"path"`
	Remote              string  `json:"remote"`
	Action              *Action `json:"action"`
}

type InitResponse struct {
	Error *OperationError `json:"error,omitempty"`
}
type TransferResponse struct {
	Event string          `json:"event"`
	Oid   string          `json:"oid"`
	Path  string          `json:"path,omitempty"` // always blank for upload
	Error *OperationError `json:"error,omitempty"`
}
type ProgressResponse struct {
	Event          string `json:"event"`
	Oid            string `json:"oid"`
	BytesSoFar     int64  `json:"bytesSoFar"`
	BytesSinceLast int    `json:"bytesSinceLast"`
}

type Client interface {
	Init(config *config.Config, remotepath string) error
	Download(oid string, size int64, a *Action)
	Upload(oid string, size int64, a *Action, fromPath string)
	Destory()
}
