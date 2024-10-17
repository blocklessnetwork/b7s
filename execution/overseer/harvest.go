package overseer

import (
	"os"

	"github.com/coder/websocket"
)

// Collect any artefacts created and remove traces of an executed job.
func (o *Overseer) harvest(id string) {

	h, ok := o.jobs[id]
	if !ok {
		return
	}

	if h.outputStream != nil {
		err := h.outputStream.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			o.log.Error().Err(err).Msg("could not close output stream")
		} else {
			o.log.Debug().Str("job", id).Msg("closed output stream")
		}
	}

	if h.errorStream != nil {
		err := h.errorStream.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			o.log.Error().Err(err).Msg("could not close error stream")
		} else {
			o.log.Debug().Str("job", id).Msg("closed error stream")
		}
	}

	err := os.RemoveAll(h.workdir)
	if err != nil {
		o.log.Error().Err(err).Msg("could not remove job workdir")
	}

	if h.source.Limits != nil {
		err = o.cfg.Limiter.DeleteGroup(id)
		if err != nil {
			o.log.Error().Err(err).Str("job", id).Msg("could not delete limit group")
		}
	}
}
