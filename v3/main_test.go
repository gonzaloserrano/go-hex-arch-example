package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gonzaloserrano/go-hex-arch-example/v3/app"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	require := require.New(t)

	var sum int
	repo := &CounterRepositoryMock{
		FindByIDFunc: func(ID string) app.Counter {
			return app.Counter{ID: ID, Value: sum}
		},
		UpsertFunc: func(c app.Counter) {
			sum = c.Value
		},
	}

	getHandler := newGetHandler(repo)
	addHandler := newAddHandler(repo)

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
			var c app.Counter
			err := json.NewDecoder(w.Body).Decode(&c)
			require.NoError(err)
			require.Equal(counterID, c.ID)
			require.Equal(tc.respValue, c.Value)
		}
	}
}
