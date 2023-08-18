package pbft

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func getDigest(req Request) string {
	payload, _ := json.Marshal(req)
	hash := sha256.Sum256(payload)

	return hex.EncodeToString(hash[:])
}

func digestOK(req Request, digest string) bool {
	d := getDigest(req)
	return d == digest
}
