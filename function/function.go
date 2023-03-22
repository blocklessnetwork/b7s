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

// NewHandler creates a new function handler.
func NewHandler(log zerolog.Logger, store Store, workdir string) *Handler {

	// Create an HTTP client.
	cli := http.Client{
		Timeout: defaultTimeout,
	}

	// Create a download client.
	downloader := grab.NewClient()
	downloader.UserAgent = defaultUserAgent

	h := Handler{
		log:        log.With().Str("component", "function_store").Logger(),
		store:      store,
		http:       &cli,
		downloader: downloader,
		workdir:    workdir,
	}

	return &h
}
