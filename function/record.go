package function

import (
	"fmt"
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
)

type functionRecord struct {
	CID      string                     `json:"cid"`
	URL      string                     `json:"url"`
	Manifest blockless.FunctionManifest `json:"manifest"`
	Archive  string                     `json:"archive"`
	Files    string                     `json:"files"`

	UpdatedAt     time.Time `json:"updated_at"`
	LastRetrieved time.Time `json:"last_retrieved"`
}

func (h *Handler) getFunction(cid string) (*functionRecord, error) {

	// Retrieve function.
	var fn functionRecord
	err := h.store.GetRecord(cid, &fn)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve function record: %w", err)
	}

	// Update the "last retrieved" timestamp.
	fn.LastRetrieved = time.Now().UTC()
	err = h.store.SetRecord(cid, fn)
	if err != nil {
		h.log.Warn().Err(err).Str("cid", cid).Msg("could not update function record timestamp")
	}

	return &fn, nil
}

func (h *Handler) saveFunction(fn functionRecord) error {

	// Clean paths - make them relative to the current working directory.
	fn.Archive = h.cleanPath(fn.Archive)
	fn.Files = h.cleanPath(fn.Files)

	fn.UpdatedAt = time.Now().UTC()
	return h.store.SetRecord(fn.CID, fn)
}
