package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/pmss"
	"github.com/stretchr/testify/assert"
)

const dbUrl = "mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true"

func TestHashRetrival(t *testing.T) {
	testCases := []struct {
		hash string
		code int
	}{
		{
			"3c738552873525fda24139e1214c95bbdaf9dcca",
			http.StatusOK,
		},
		{
			"test",
			http.StatusBadRequest,
		},
	}

	db := datastore.MustOpen(dbUrl)
	datastore.MustAutoMigrate(db)

	err := datastore.KnownFiles().Create(datastore.KnownFile{
		MD5:    utils.StringPtr("d3b07384d113edec49eaa6238ad5ff00"),
		SHA1:   utils.StringPtr("3c738552873525fda24139e1214c95bbdaf9dcca"),
		SHA256: utils.StringPtr("d3b07384d113edec49eaa6238ad5ff00d3b07384d113edec49eaa6238ad5ff00"),
		Size:   utils.Int64Ptr(100),
		Status: datastore.FileStatusMalicious,
	})(db)

	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
		return
	}

	pmss, err := pmss.New(pmss.WithDbUrl(dbUrl))

	listenAddr := ":0"
	if err != nil {
		t.Fatalf("failed to initialize PMSS: %v", err)
		return
	}
	srv := New(context.Background(), pmss, WithListenAddr(listenAddr))
	router := srv.httpSrv.Handler

	for _, testcase := range testCases {
		t.Run(testcase.hash, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/hash/%s", testcase.hash), nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, testcase.code, w.Code)
			w.Flush()
		})
	}
}

func TestReporting(t *testing.T) {
	t.SkipNow()
	pmss, err := pmss.New(pmss.WithDbUrl("mysql://pmss:pmss@tcp(mariadb:3306)/pmss?parseTime=true"))
	listenAddr := ":0"
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
