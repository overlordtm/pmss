package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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
		Handler: srv.router(),
	}

	return srv
}

func (s *Server) router() *gin.Engine {
	router := gin.Default()
	router.GET("/api/v1/hash/:hash", s.findByHash)
	router.POST("/api/v1/report", s.bulkReport)
	return router
}

func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpSrv.Close()
}

func (s *Server) findByHash(c *gin.Context) {

	hash := c.Param("hash")
	if hash == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "hash is required"})
		return
	}
	result, err := s.pmss.FindByHash(hash)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

type FileReport struct {
	Path   string `json:"path"`
	MD5    string `json:"md5,omitempty"`
	SHA1   string `json:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

type BulkReportRequest struct {
	Reports        []FileReport `json:"reports"`
	MachineHostame string       `json:"machine_hostname"`
	MachineApiKey  string       `json:"machine_api_key"`
}

func (s *Server) canMachineReport(machineHostname, machineApiKey string) (bool, error) {
	// SELECT * FROM machines `m` WHERE m.hostname = ? AND m.api_key = ? AND m.allow_submit LIMIT 1
	return false, nil
}

func (s *Server) insertReports(reports []FileReport) error {

	// INSERT INTO runs (time) VALUES (NOW())

	// INSERT INTO reports (path, md5, sha1, sha256) VALUES (?, ?, ?, ?)

	return nil
}

func (s *Server) bulkReport(c *gin.Context) {
	//Read json report
	var req BulkReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//Should bind json send response automatically
		return
	}

	//Validate submitters MachineApiKey with combination of MachineHostame
	canReport, err := s.canMachineReport(req.MachineHostame, req.MachineApiKey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if !canReport {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Save report to database
	if err := s.insertReports(req.Reports); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}
