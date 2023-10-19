package overseer

import (
	"github.com/spf13/afero"
)

type Config struct {
	Workdir string
	FS      afero.Fs
}
