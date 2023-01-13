package keygen

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeys(t *testing.T) {
	// setup
	outputFolder := "/tmp/test_keygen"
	os.RemoveAll(outputFolder)

	// test GenerateKeys
	err := GenerateKeys(outputFolder)
	assert.Nil(t, err)

	// check that the files exist
	_, err = os.Stat(outputFolder + "/pub.bin")
	assert.Nil(t, err)
	_, err = os.Stat(outputFolder + "/priv.bin")
	assert.Nil(t, err)
	_, err = os.Stat(outputFolder + "/identity")
	assert.Nil(t, err)

	// check that the files are not empty
	pubBytes, err := ioutil.ReadFile(outputFolder + "/pub.bin")
	assert.Nil(t, err)
	assert.NotEmpty(t, pubBytes)
	privBytes, err := ioutil.ReadFile(outputFolder + "/priv.bin")
	assert.Nil(t, err)
	assert.NotEmpty(t, privBytes)
	peerIdBytes, err := ioutil.ReadFile(outputFolder + "/identity")
	assert.Nil(t, err)
	assert.NotEmpty(t, peerIdBytes)
}
