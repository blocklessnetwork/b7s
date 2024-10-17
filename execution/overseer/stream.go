package overseer

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
)

const (
	wsTimeout = 100 * time.Millisecond
)

func wsConnect(address string) (*websocket.Conn, error) {

	opts := websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	}

	conn, _, err := websocket.Dial(context.Background(), address, &opts)
	if err != nil {
		return nil, fmt.Errorf("could not dial: %w", err)
	}

	return conn, nil
}

type wsWriter struct {
	conn *websocket.Conn
	log  zerolog.Logger
}

// Implement the write interface but ignore the error.
func (w *wsWriter) Write(p []byte) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), wsTimeout)
	defer cancel()

	err := w.conn.Write(ctx, websocket.MessageBinary, p)
	if err != nil {
		w.log.Trace().Err(err).Msg("could not write to stream")
	}

	return len(p), nil
}
