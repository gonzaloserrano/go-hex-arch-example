package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	dialect "upper.io/db.v3/postgresql"
)

func TestCounter(t *testing.T) {
	require := require.New(t)

	db, err := dialect.Open(dialect.ConnectionURL{
		Database: "test",
		User:     "test",
		Host:     "localhost:54321",
	})
	require.NoError(err)

	getHandler := newGetHandler(db)
	addHandler := newAddHandler(db)

	counterID := uuid.New().String()
	for _, tc := range []struct {
		handler   http.HandlerFunc
		reqMethod string
		reqBody   string
		respValue int
	}{
		{addHandler, "POST", fmt.Sprintf(`{"id":"%s", "value":%d}`, counterID, 5), 0},
		{getHandler, "GET", fmt.Sprintf(`{"id":"%s"}`, counterID), 5},
		{addHandler, "POST", fmt.Sprintf(`{"id":"%s", "value":%d}`, counterID, 11), 0},
		{getHandler, "GET", fmt.Sprintf(`{"id":"%s"}`, counterID), 16},
	} {
		body := strings.NewReader(tc.reqBody)
		req := httptest.NewRequest(tc.reqMethod, "http://dont.care/", body)
		w := httptest.NewRecorder()
		tc.handler(w, req)

		require.Equal(http.StatusOK, w.Code)

		if tc.reqMethod == "GET" {
			var c counter
			err = json.NewDecoder(w.Body).Decode(&c)
			require.NoError(err)
			require.Equal(counterID, c.ID)
			require.Equal(tc.respValue, c.Value)
		}
	}
}
