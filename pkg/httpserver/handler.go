package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/overlordtm/pmss/internal/apiserver"
	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/apitypes"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/hashtype"
	"github.com/overlordtm/pmss/pkg/pmss"
	"gorm.io/gorm"
)

type handler struct {
	*pmss.Pmss
}

// typecheck
var _ apiserver.ServerInterface = &handler{}

func (h *handler) QueryByHash(c *gin.Context, hash string) {
	if hash == "" {
		c.Error(fmt.Errorf("hash is required"))
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := h.Pmss.FindByHash(hash)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, apitypes.HashQueryResponse{
				Status: datastore.FileStatusUnknown,
			})
			return
		case hashtype.ErrUnknown:
			c.Error(fmt.Errorf("error: %s %s %d", err.Error(), hash, len(hash)))
			c.Status(http.StatusBadRequest)
			return
		default:
			c.Error(fmt.Errorf("error: %s %s", err.Error(), hash))
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	res := apitypes.HashQueryResponse{
		Status: result.File.Status,
		File: &apitypes.KnownFile{
			Md5:    result.File.MD5,
			Sha1:   result.File.SHA1,
			Sha256: result.File.SHA256,
			Size:   result.File.Size,
		},
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) SubmitReport(c *gin.Context) {

	//Read json report
	var req apitypes.NewReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//Should bind json send response automatically
		c.Error(fmt.Errorf("ShouldBindJSON error: %s", err.Error()))
		return
	}

	files := make([]datastore.ScannedFile, len(req.Files))
	reportFiles := make([]apitypes.ReportFile, len(req.Files))
	for i, f := range req.Files {
		files[i] = datastore.ScannedFile{
			Path:   f.Path,
			MD5:    f.Md5,
			SHA1:   f.Sha1,
			SHA256: f.Sha256,

			Ctime: f.Ctime,
			Mtime: &f.Mtime,

			Size: f.Size,
			Mode: f.FileMode,
		}

		knownFile, err := h.Pmss.ScanFile(&files[i])

		reportFiles[i] = apitypes.ReportFile{
			Path:   f.Path,
			Status: knownFile.Status,
			Error:  utils.ErrToStrPtr(err),
		}
	}

	run, err := h.Pmss.DoMachineReport(&pmss.ScanReport{
		Hostname:  req.Hostname,
		Files:     files,
		IP:        c.ClientIP(),
		MachineId: req.MachineId,
		ScanRunId: req.ReportRunId,
	})
	if err != nil {
		c.Error(fmt.Errorf("DoMachineReport error: %s", err.Error()))
		return
	}

	c.IndentedJSON(http.StatusCreated, apitypes.NewReportResponse{
		Id:    run.ID,
		Files: reportFiles,
	})
}
