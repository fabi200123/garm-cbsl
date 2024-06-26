// Code generated by go-swagger; DO NOT EDIT.

package organizations

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

	garm_params "github.com/cloudbase/garm/params"
)

// NewUpdateOrgParams creates a new UpdateOrgParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateOrgParams() *UpdateOrgParams {
	return &UpdateOrgParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateOrgParamsWithTimeout creates a new UpdateOrgParams object
// with the ability to set a timeout on a request.
func NewUpdateOrgParamsWithTimeout(timeout time.Duration) *UpdateOrgParams {
	return &UpdateOrgParams{
		timeout: timeout,
	}
}

// NewUpdateOrgParamsWithContext creates a new UpdateOrgParams object
// with the ability to set a context for a request.
func NewUpdateOrgParamsWithContext(ctx context.Context) *UpdateOrgParams {
	return &UpdateOrgParams{
		Context: ctx,
	}
}

// NewUpdateOrgParamsWithHTTPClient creates a new UpdateOrgParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateOrgParamsWithHTTPClient(client *http.Client) *UpdateOrgParams {
	return &UpdateOrgParams{
		HTTPClient: client,
	}
}

/*
UpdateOrgParams contains all the parameters to send to the API endpoint

	for the update org operation.

	Typically these are written to a http.Request.
*/
type UpdateOrgParams struct {

	/* Body.

	   Parameters used when updating the organization.
	*/
	Body garm_params.UpdateEntityParams

	/* OrgID.

	   ID of the organization to update.
	*/
	OrgID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update org params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateOrgParams) WithDefaults() *UpdateOrgParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update org params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateOrgParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update org params
func (o *UpdateOrgParams) WithTimeout(timeout time.Duration) *UpdateOrgParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update org params
func (o *UpdateOrgParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update org params
func (o *UpdateOrgParams) WithContext(ctx context.Context) *UpdateOrgParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update org params
func (o *UpdateOrgParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update org params
func (o *UpdateOrgParams) WithHTTPClient(client *http.Client) *UpdateOrgParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update org params
func (o *UpdateOrgParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update org params
func (o *UpdateOrgParams) WithBody(body garm_params.UpdateEntityParams) *UpdateOrgParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update org params
func (o *UpdateOrgParams) SetBody(body garm_params.UpdateEntityParams) {
	o.Body = body
}

// WithOrgID adds the orgID to the update org params
func (o *UpdateOrgParams) WithOrgID(orgID string) *UpdateOrgParams {
	o.SetOrgID(orgID)
	return o
}

// SetOrgID adds the orgId to the update org params
func (o *UpdateOrgParams) SetOrgID(orgID string) {
	o.OrgID = orgID
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateOrgParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param orgID
	if err := r.SetPathParam("orgID", o.OrgID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
