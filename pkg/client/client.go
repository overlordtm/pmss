package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/overlordtm/pmss/internal/apiclient"
)

type Client interface {
	// ScanFeatures(f FileFeatures) (detector.Result, error)
	SubmitFiles(files apiclient.NewReportRequest)
}

type HTTPClient struct {
	client      apiclient.ClientWithResponsesInterface
	reqHandlers []apiclient.RequestEditorFn
}

func New(apiUrl string) (*HTTPClient, error) {

	c := &HTTPClient{}

	opts := []apiclient.ClientOption{
		apiclient.WithRequestEditorFn(c.authorizeReq()),
	}

	if client, err := apiclient.NewClientWithResponses(apiUrl, opts...); err != nil {
		return nil, fmt.Errorf("failed to initialize api client: %v", err)
	} else {
		c.client = client
	}

	return c, nil
}

func (c *HTTPClient) authorizeReq() apiclient.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer dummy")
		return nil
	}
}

func (c *HTTPClient) SubmitFiles(ctx context.Context, files apiclient.NewReportRequest) (*apiclient.SubmitReportResponse, error) {
	return c.client.SubmitReportWithResponse(ctx, files, c.reqHandlers...)
}
