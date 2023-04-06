package fstore

import (
	"net/http"

	"github.com/cavaliergopher/grab/v3"
	"github.com/rs/zerolog"
)

// FStore - function store - deals with all of the function-related actions - saving/reading them from backing storage,
// downloading them, unpacking them etc.
type FStore struct {
	log        zerolog.Logger
	store      Store
	http       *http.Client
	downloader *grab.Client

	workdir string
}

// New creates a new function store.
func New(log zerolog.Logger, store Store, workdir string) *FStore {

	// Create an HTTP client.
	cli := http.Client{
		Timeout: defaultTimeout,
	}

	// Create a download client.
	downloader := grab.NewClient()
	downloader.UserAgent = defaultUserAgent

	h := FStore{
		log:        log.With().Str("component", "function_store").Logger(),
		store:      store,
		http:       &cli,
		downloader: downloader,
		workdir:    workdir,
	}

	return &h
}

// InstalledFunctions will return the CIDs of all functions found in local storage.
func (h *FStore) InstalledFunctions() []string {
	return h.store.Keys()
}
