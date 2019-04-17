package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMap(t *testing.T) {
	rr := httptest.NewRecorder()
	GetMap(rr, httptest.NewRequest("GET", "/random/map", nil))
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)

	data, err := ioutil.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	require.Equal(t, mapSize, len(data), "map size doesn't match")

	rows := bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
	require.Equal(t, mapHeight, len(rows), "invalid number of rows")

	for _, row := range rows {
		require.Equal(t, mapWidth, len(row), "invalid number characters in row")
		for _, c := range row {
			require.Contains(t, []byte("ox"), c, "unexpected map character")
		}
	}
}
