package pbft

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func getDigest(rec any) string {
	payload, _ := json.Marshal(rec)
	hash := sha256.Sum256(payload)

	return hex.EncodeToString(hash[:])
}
