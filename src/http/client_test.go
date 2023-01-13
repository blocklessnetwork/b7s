package http

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func TestGetJson(t *testing.T) {
	// setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonData := testStruct{
			Field1: "value1",
			Field2: 2,
		}
		json.NewEncoder(w).Encode(jsonData)
	}))
	defer ts.Close()

	// test GetJson
	var target testStruct
	err := GetJson(ts.URL, &target)
	assert.Nil(t, err)
	assert.Equal(t, "value1", target.Field1)
	assert.Equal(t, 2, target.Field2)
}
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
