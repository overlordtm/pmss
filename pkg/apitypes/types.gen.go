// Package apitypes provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package apitypes

import (
	"github.com/google/uuid"
	"github.com/overlordtm/pmss/pkg/datastore"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// File defines model for File.
type File struct {
	Ctime    *int64  `json:"ctime,omitempty"`
	FileMode uint32  `json:"fileMode"`
	Md5      *string `json:"md5,omitempty"`
	Mtime    int64   `json:"mtime"`
	Path     string  `json:"path"`
	Sha1     *string `json:"sha1,omitempty"`
	Sha256   *string `json:"sha256,omitempty"`
	Size     int64   `json:"size"`
}

// HashQueryResponse defines model for HashQueryResponse.
type HashQueryResponse struct {
	File   *KnownFile           `json:"file,omitempty"`
	Status datastore.FileStatus `json:"status"`
}

// KnownFile defines model for KnownFile.
type KnownFile struct {
	Md5    *string `json:"md5,omitempty"`
	Sha1   *string `json:"sha1,omitempty"`
	Sha256 *string `json:"sha256,omitempty"`
	Size   *int64  `json:"size,omitempty"`
}

// NewReportRequest defines model for NewReportRequest.
type NewReportRequest struct {
	Files       []File     `json:"files"`
	Hostname    string     `json:"hostname"`
	MachineId   string     `json:"machineId"`
	ReportRunId *uuid.UUID `json:"reportRunId,omitempty"`
}

// NewReportResponse defines model for NewReportResponse.
type NewReportResponse struct {
	Id uuid.UUID `json:"id"`
}

// SubmitReportJSONRequestBody defines body for SubmitReport for application/json ContentType.
type SubmitReportJSONRequestBody = NewReportRequest
