package store

import (
	"fmt"
)

func encodeKey(prefix uint8, segments ...any) []byte {

	key := []byte{prefix}

	for _, segment := range segments {
		switch s := segment.(type) {

		case []byte: // e.g. peer.ID
			key = append(key, Separator)
			key = append(key, s...)

		case string:
			key = append(key, Separator)
			key = append(key, []byte(s)...)

		default:
			panic(fmt.Sprintf("unsupported type (%T)", segment))
		}
	}

	return key
}
