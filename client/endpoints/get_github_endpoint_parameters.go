// Code generated by go-swagger; DO NOT EDIT.

package endpoints

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
)

// NewGetGithubEndpointParams creates a new GetGithubEndpointParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetGithubEndpointParams() *GetGithubEndpointParams {
	return &GetGithubEndpointParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetGithubEndpointParamsWithTimeout creates a new GetGithubEndpointParams object
// with the ability to set a timeout on a request.
func NewGetGithubEndpointParamsWithTimeout(timeout time.Duration) *GetGithubEndpointParams {
	return &GetGithubEndpointParams{
		timeout: timeout,
	}
}

// NewGetGithubEndpointParamsWithContext creates a new GetGithubEndpointParams object
// with the ability to set a context for a request.
func NewGetGithubEndpointParamsWithContext(ctx context.Context) *GetGithubEndpointParams {
	return &GetGithubEndpointParams{
		Context: ctx,
	}
}

// NewGetGithubEndpointParamsWithHTTPClient creates a new GetGithubEndpointParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetGithubEndpointParamsWithHTTPClient(client *http.Client) *GetGithubEndpointParams {
	return &GetGithubEndpointParams{
		HTTPClient: client,
	}
}

/*
GetGithubEndpointParams contains all the parameters to send to the API endpoint

	for the get github endpoint operation.

	Typically these are written to a http.Request.
*/
type GetGithubEndpointParams struct {

	/* Name.

	   The name of the GitHub endpoint.
	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get github endpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGithubEndpointParams) WithDefaults() *GetGithubEndpointParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get github endpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGithubEndpointParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get github endpoint params
func (o *GetGithubEndpointParams) WithTimeout(timeout time.Duration) *GetGithubEndpointParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get github endpoint params
func (o *GetGithubEndpointParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get github endpoint params
func (o *GetGithubEndpointParams) WithContext(ctx context.Context) *GetGithubEndpointParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get github endpoint params
func (o *GetGithubEndpointParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get github endpoint params
func (o *GetGithubEndpointParams) WithHTTPClient(client *http.Client) *GetGithubEndpointParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get github endpoint params
func (o *GetGithubEndpointParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get github endpoint params
func (o *GetGithubEndpointParams) WithName(name string) *GetGithubEndpointParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get github endpoint params
func (o *GetGithubEndpointParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetGithubEndpointParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
