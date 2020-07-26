package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	go run()

	time.Sleep(10 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/run")
	require.NoError(t, err)

	data, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, welcomeMessage, string(data))
}
