package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/stretchr/testify/assert"
)

func TestHashRetrival(t *testing.T) {

	testCases := []struct {
		hash string
		code int
		file string
	}{
		{
			"test",
			http.StatusBadRequest,
			"",
		},
		{
			"3c738552873525fda24139e1214c95bbdaf9dcca",
			http.StatusOK,
			"/bin/ps",
		},
	}

	pmss, err := pmss.New("mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true")
	listenAddr := ":6069"
	if err != nil {
		t.Fatalf("failed to initialize PMSS: %v", err)
		return
	}
	srv := New(context.Background(), pmss, WithListenAddr(listenAddr))
	router := srv.httpSrv.Handler

	for _, testcase := range testCases {
		t.Run(testcase.hash, func(t *testing.T) {
			w := httptest.NewRecorder()
			endpoint := fmt.Sprintf("/api/v1/hash/%s", testcase.hash)
			req, _ := http.NewRequest("GET", endpoint, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, testcase.code, w.Code)
			// TODO add file check
		})
	}
}

func TestReporting(t *testing.T) {
	t.Error("Not implemented")
	pmss, err := pmss.New("mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true")
	listenAddr := ":6069"
	if err != nil {
		t.Fatalf("failed to initialize PMSS: %v", err)
		return
	}
	srv := New(context.Background(), pmss, WithListenAddr(listenAddr))
	router := srv.httpSrv.Handler
	w := httptest.NewRecorder()
	endpoint := fmt.Sprintf("/api/v1/report")
	req, _ := http.NewRequest("POST", endpoint, nil)
	router.ServeHTTP(w, req)
	//assert.Equal(t, testcase.code, w.Code)
	// TODO Implement report testing
}
