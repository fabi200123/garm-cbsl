// Code generated by go-swagger; DO NOT EDIT.

package organizations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	apiserver_params "github.com/cloudbase/garm/apiserver/params"
	garm_params "github.com/cloudbase/garm/params"
)

// CreateOrgPoolReader is a Reader for the CreateOrgPool structure.
type CreateOrgPoolReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateOrgPoolReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateOrgPoolOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCreateOrgPoolDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateOrgPoolOK creates a CreateOrgPoolOK with default headers values
func NewCreateOrgPoolOK() *CreateOrgPoolOK {
	return &CreateOrgPoolOK{}
}

/*
CreateOrgPoolOK describes a response with status code 200, with default header values.

Pool
*/
type CreateOrgPoolOK struct {
	Payload garm_params.Pool
}

// IsSuccess returns true when this create org pool o k response has a 2xx status code
func (o *CreateOrgPoolOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create org pool o k response has a 3xx status code
func (o *CreateOrgPoolOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create org pool o k response has a 4xx status code
func (o *CreateOrgPoolOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this create org pool o k response has a 5xx status code
func (o *CreateOrgPoolOK) IsServerError() bool {
	return false
}

// IsCode returns true when this create org pool o k response a status code equal to that given
func (o *CreateOrgPoolOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the create org pool o k response
func (o *CreateOrgPoolOK) Code() int {
	return 200
}

func (o *CreateOrgPoolOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /organizations/{orgID}/pools][%d] createOrgPoolOK %s", 200, payload)
}

func (o *CreateOrgPoolOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /organizations/{orgID}/pools][%d] createOrgPoolOK %s", 200, payload)
}

func (o *CreateOrgPoolOK) GetPayload() garm_params.Pool {
	return o.Payload
}

func (o *CreateOrgPoolOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateOrgPoolDefault creates a CreateOrgPoolDefault with default headers values
func NewCreateOrgPoolDefault(code int) *CreateOrgPoolDefault {
	return &CreateOrgPoolDefault{
		_statusCode: code,
	}
}

/*
CreateOrgPoolDefault describes a response with status code -1, with default header values.

APIErrorResponse
*/
type CreateOrgPoolDefault struct {
	_statusCode int

	Payload apiserver_params.APIErrorResponse
}

// IsSuccess returns true when this create org pool default response has a 2xx status code
func (o *CreateOrgPoolDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this create org pool default response has a 3xx status code
func (o *CreateOrgPoolDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this create org pool default response has a 4xx status code
func (o *CreateOrgPoolDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this create org pool default response has a 5xx status code
func (o *CreateOrgPoolDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this create org pool default response a status code equal to that given
func (o *CreateOrgPoolDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the create org pool default response
func (o *CreateOrgPoolDefault) Code() int {
	return o._statusCode
}

func (o *CreateOrgPoolDefault) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /organizations/{orgID}/pools][%d] CreateOrgPool default %s", o._statusCode, payload)
}

func (o *CreateOrgPoolDefault) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[POST /organizations/{orgID}/pools][%d] CreateOrgPool default %s", o._statusCode, payload)
}

func (o *CreateOrgPoolDefault) GetPayload() apiserver_params.APIErrorResponse {
	return o.Payload
}

func (o *CreateOrgPoolDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
