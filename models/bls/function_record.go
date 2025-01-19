package bls

import (
	"time"
)

type FunctionRecord struct {
	CID      string           `json:"cid"`
	URL      string           `json:"url"`
	Manifest FunctionManifest `json:"manifest"`
	Archive  string           `json:"archive"`
	Files    string           `json:"files"`

	UpdatedAt     time.Time `json:"updated_at"`
	LastRetrieved time.Time `json:"last_retrieved"`
}
