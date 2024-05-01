package store

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

func encodeKey(prefix uint8, segments ...any) ([]byte, error) {

	key := []byte{prefix}

	for _, segment := range segments {
		switch s := segment.(type) {

		case peer.ID:

			id, err := s.MarshalBinary()
			if err != nil {
				return nil, err
			}

			key = append(key, Separator)
			key = append(key, id...)

		case string:
			key = append(key, Separator)
			key = append(key, []byte(s)...)

		default:
			panic(fmt.Sprintf("unsupported type (%T)", segment))
		}
	}

	return key, nil
}
