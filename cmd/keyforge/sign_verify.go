package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// HandleSignAndVerify performs the signing and verification based on the provided keys and flags
func HandleSignAndVerify(priv crypto.PrivKey, pub crypto.PubKey, flagString string, flagFile string, flagOutput string) {
	// Sign and verify string
	if flagString != "" {
		signature, err := priv.Sign([]byte(flagString))
		if err != nil {
			log.Fatalf("Could not sign string: %s", err)
		}

		isValid, err := pub.Verify([]byte(flagString), signature)
		if err != nil {
			log.Fatalf("Could not verify string: %s", err)
		}

		signatureStr := base64.StdEncoding.EncodeToString(signature)
		if isValid {
			fmt.Printf("String signature verified successfully. Signature: %s\n", signatureStr)
		} else {
			fmt.Println("String signature verification failed.")
		}
	}

	// Sign and verify file
	if flagFile != "" {
		fileData, err := ioutil.ReadFile(flagFile)
		if err != nil {
			log.Fatalf("Could not read file: %s", err)
		}

		fileSignature, err := priv.Sign(fileData)
		if err != nil {
			log.Fatalf("Could not sign file: %s", err)
		}

		isFileValid, err := pub.Verify(fileData, fileSignature)
		if err != nil {
			log.Fatalf("Could not verify file: %s", err)
		}

		signatureFileStr := base64.StdEncoding.EncodeToString(fileSignature)
		if isFileValid {
			fmt.Printf("File signature verified successfully. Signature: %s\n", signatureFileStr)
		} else {
			fmt.Println("File signature verification failed.")
		}

		if flagOutput != "" {
			signatureFilePath := filepath.Join(flagOutput, "file_signature.enc.txt")
			err := ioutil.WriteFile(signatureFilePath, []byte(signatureFileStr), 0644)
			if err != nil {
				log.Fatalf("Could not write signature to file: %s", err)
			}
		}
	}
}

// VerifyGivenSignature verifies a message with a given signature and public key
func VerifyGivenSignature(encodedPubKey string, message string, encodedSignature string) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(encodedPubKey)
	if err != nil {
		log.Fatalf("Could not decode public key: %s", err)
	}

	pubKey, err := crypto.UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		log.Fatalf("Could not unmarshal public key: %s", err)
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		log.Fatalf("Could not decode signature: %s", err)
	}

	isValid, err := pubKey.Verify([]byte(message), signatureBytes)
	if err != nil {
		log.Fatalf("Could not verify string: %s", err)
	}

	if isValid {
		fmt.Println("Signature verified successfully.")
	} else {
		fmt.Println("Signature verification failed.")
	}
}

// VerifyGivenSignatureWithPeerID verifies a message with a given signature and PeerID
func VerifyGivenSignatureWithPeerID(peerIDStr string, message string, encodedSignature string) {
	log.Println("Verifying signature with peerid")
	peerID, err := peer.Decode(peerIDStr)
	if err != nil {
		log.Fatalf("Could not decode PeerID: %s", err)
	}

	pubKey, err := peerID.ExtractPublicKey()
	if err != nil || pubKey == nil {
		log.Fatalf("Could not extract public key from PeerID: %s", err)
	}

	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Could not marshal public key to bytes: %s", err)
	}

	// This will give you a base64 string representation of the public key
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKeyBytes)
	log.Printf("Extracted public key from PeerID (Base64): %s", pubKeyBase64)

	VerifyGivenSignature(pubKeyBase64, message, encodedSignature)
}
