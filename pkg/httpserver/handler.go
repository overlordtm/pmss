package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlordtm/pmss/internal/apiserver"
	"github.com/overlordtm/pmss/pkg/apitypes"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/pmss"
)

type handler struct {
	*pmss.Pmss
}

// typecheck
var _ apiserver.ServerInterface = &handler{}

func (h *handler) QueryByHash(c *gin.Context, hash string) {
	if hash == "" {
		c.Error(fmt.Errorf("hash is required"))
		return
	}

	result, err := h.Pmss.FindByHash(hash)

	if err != nil {
		c.Error(fmt.Errorf("error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) SubmitReport(c *gin.Context) {

	//Read json report
	var req apitypes.NewReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//Should bind json send response automatically
		c.Error(fmt.Errorf("error: %s", err.Error()))
		return
	}

	files := make([]datastore.ScannedFile, len(req.Files))
	for i, f := range req.Files {
		files[i] = datastore.ScannedFile{
			MD5:    f.Md5,
			SHA1:   f.Sha1,
			SHA256: f.Sha256,

			Ctime: f.Ctime,
			Mtime: &f.Mtime,

			Size: f.Size,
			Mode: f.FileMode,
		}
	}

	run, err := h.Pmss.DoMachineReport(&pmss.ScanReport{
		Hostname:  req.Hostname,
		Files:     files,
		MachineId: req.MachineId,
	})
	if err != nil {
		c.Error(fmt.Errorf("error: %s", err.Error()))
		return
	}

	c.IndentedJSON(http.StatusCreated, apitypes.NewReportResponse{
		Id: run.ID,
	})
}
