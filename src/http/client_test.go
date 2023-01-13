package http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	// setup test server
	// create a test server to simulate the function repository
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("testdata/testfile")
		defer file.Close()

		// Send the file to the client
		io.Copy(w, file)
	}))
	defer ts.Close()

	// setup context
	config := models.Config{
		Node: models.ConfigNode{
			WorkspaceRoot: "/tmp/test_workspace",
		},
	}
	ctx := context.WithValue(context.Background(), "config", config)

	// setup test function manifest
	functionManifest := models.FunctionManifest{
		Function: models.Function{
			ID: "test_function",
		},
		Deployment: models.Deployment{
			Uri:      ts.URL + "/testfile",
			Checksum: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}

	// test Download
	filepath, err := Download(ctx, functionManifest)
	assert.Nil(t, err)
	assert.Equal(t, "/tmp/test_workspace/test_function/testfile", filepath)

	// check if file exists
	_, err = os.Stat(filepath)
	assert.Nil(t, err)

	// check if file content is correct
	content, _ := ioutil.ReadFile(filepath)
	assert.Equal(t, "hello", string(content))

	// cleanup
	os.RemoveAll("/tmp/test_workspace")
}
