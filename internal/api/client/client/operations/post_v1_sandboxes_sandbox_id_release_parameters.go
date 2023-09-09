// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/nuonco/terraform-provider-nuon/internal/api/client/models"
)

// NewPostV1SandboxesSandboxIDReleaseParams creates a new PostV1SandboxesSandboxIDReleaseParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostV1SandboxesSandboxIDReleaseParams() *PostV1SandboxesSandboxIDReleaseParams {
	return &PostV1SandboxesSandboxIDReleaseParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostV1SandboxesSandboxIDReleaseParamsWithTimeout creates a new PostV1SandboxesSandboxIDReleaseParams object
// with the ability to set a timeout on a request.
func NewPostV1SandboxesSandboxIDReleaseParamsWithTimeout(timeout time.Duration) *PostV1SandboxesSandboxIDReleaseParams {
	return &PostV1SandboxesSandboxIDReleaseParams{
		timeout: timeout,
	}
}

// NewPostV1SandboxesSandboxIDReleaseParamsWithContext creates a new PostV1SandboxesSandboxIDReleaseParams object
// with the ability to set a context for a request.
func NewPostV1SandboxesSandboxIDReleaseParamsWithContext(ctx context.Context) *PostV1SandboxesSandboxIDReleaseParams {
	return &PostV1SandboxesSandboxIDReleaseParams{
		Context: ctx,
	}
}

// NewPostV1SandboxesSandboxIDReleaseParamsWithHTTPClient creates a new PostV1SandboxesSandboxIDReleaseParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostV1SandboxesSandboxIDReleaseParamsWithHTTPClient(client *http.Client) *PostV1SandboxesSandboxIDReleaseParams {
	return &PostV1SandboxesSandboxIDReleaseParams{
		HTTPClient: client,
	}
}

/*
PostV1SandboxesSandboxIDReleaseParams contains all the parameters to send to the API endpoint

	for the post v1 sandboxes sandbox ID release operation.

	Typically these are written to a http.Request.
*/
type PostV1SandboxesSandboxIDReleaseParams struct {

	/* Req.

	   Input
	*/
	Req *models.ServiceCreateSandboxReleaseRequest

	/* SandboxID.

	   sandbox ID
	*/
	SandboxID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post v1 sandboxes sandbox ID release params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1SandboxesSandboxIDReleaseParams) WithDefaults() *PostV1SandboxesSandboxIDReleaseParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post v1 sandboxes sandbox ID release params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1SandboxesSandboxIDReleaseParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) WithTimeout(timeout time.Duration) *PostV1SandboxesSandboxIDReleaseParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) WithContext(ctx context.Context) *PostV1SandboxesSandboxIDReleaseParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) WithHTTPClient(client *http.Client) *PostV1SandboxesSandboxIDReleaseParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithReq adds the req to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) WithReq(req *models.ServiceCreateSandboxReleaseRequest) *PostV1SandboxesSandboxIDReleaseParams {
	o.SetReq(req)
	return o
}

// SetReq adds the req to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) SetReq(req *models.ServiceCreateSandboxReleaseRequest) {
	o.Req = req
}

// WithSandboxID adds the sandboxID to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) WithSandboxID(sandboxID string) *PostV1SandboxesSandboxIDReleaseParams {
	o.SetSandboxID(sandboxID)
	return o
}

// SetSandboxID adds the sandboxId to the post v1 sandboxes sandbox ID release params
func (o *PostV1SandboxesSandboxIDReleaseParams) SetSandboxID(sandboxID string) {
	o.SandboxID = sandboxID
}

// WriteToRequest writes these params to a swagger request
func (o *PostV1SandboxesSandboxIDReleaseParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Req != nil {
		if err := r.SetBodyParam(o.Req); err != nil {
			return err
		}
	}

	// path param sandbox_id
	if err := r.SetPathParam("sandbox_id", o.SandboxID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
