package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/response"
)

func TestAPI_HealthResult(t *testing.T) {
	t.Run("nominal case", func(t *testing.T) {
		t.Parallel()

		api := setupAPI(t)

		rec, ctx, err := setupRecorder(healthEndpoint, nil)
		require.NoError(t, err)

		err = api.Health(ctx)
		require.NoError(t, err)

		var res response.Health
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

		require.Equal(t, http.StatusOK, res.Code)
		require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	})
}
