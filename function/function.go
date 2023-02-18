package function

import (
	"net/http"

	"github.com/cavaliergopher/grab/v3"
	"github.com/rs/zerolog"
)

// Handler deals with all of the function-related actions - saving/reading them from backing storage,
// downloading them, unpacking them etc.
type Handler struct {
	log        zerolog.Logger
	store      Store
	http       *http.Client
	downloader *grab.Client

	workdir string
}

// New creates a new function handler.
func New(log zerolog.Logger, store Store, workdir string) *Handler {

	// Create an HTTP client.
	cli := http.Client{
		Timeout: defaultTimeout,
	}

	// Create a download client.
	downloader := grab.NewClient()
	downloader.UserAgent = defaultUserAgent

	h := Handler{
		log:        log,
		store:      store,
		http:       &cli,
		downloader: downloader,
		workdir:    workdir,
	}

	return &h
}