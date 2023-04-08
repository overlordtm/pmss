package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/pmss"
)

type options struct {
	listenAddr string
}

type Option func(*options)

func WithListenAddr(addr string) Option {
	return func(o *options) {
		o.listenAddr = addr
	}
}

type Server struct {
	options options
	pmss    *pmss.Pmss
	httpSrv *http.Server
}

func New(ctx context.Context, pmss *pmss.Pmss, opts ...Option) *Server {

	o := options{
		listenAddr: ":8080",
	}

	for _, opt := range opts {
		opt(&o)
	}

	srv := &Server{
		options: o,
		pmss:    pmss,
	}

	srv.httpSrv = &http.Server{
		Addr:    srv.options.listenAddr,
		Handler: srv.setupRouter(),
	}

	return srv
}

func (s *Server) setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/api/v1/hash/:hash", s.handleHashInfo)
	router.POST("/api/v1/report", s.handleBulkReport)
	return router
}

func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpSrv.Close()
}

func (s *Server) handleHashInfo(c *gin.Context) {

	hash := c.Param("hash")
	if hash == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "hash is required"})
		return
	}
	malformed, result, err := s.pmss.FindByHash(hash)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		fmt.Printf("%v\n", err)
		return
	}
	status := http.StatusOK
	if malformed {
		status = http.StatusBadRequest
	}
	c.IndentedJSON(status, result)
}

type BulkReportRequest struct {
	Reports        []datastore.ScannedFile `json:"reports"`
	MachineHostame string                  `json:"machine_hostname"`
	MachineApiKey  string                  `json:"machine_api_key"`
}

func (s *Server) handleBulkReport(c *gin.Context) {

	machine := &datastore.Machine{}
	//Read json report
	var req BulkReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//Should bind json send response automatically
		return
	}

	//Validate submitters MachineApiKey with combination of MachineHostame
	canReport, err := s.pmss.FindMachineByHostname(req.MachineHostame, req.MachineApiKey, machine)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if !canReport {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Save report to database
	var run datastore.ReportRun
	if err := s.pmss.DoMachineReport(req.Reports, machine, &run); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusCreated, run)

}
