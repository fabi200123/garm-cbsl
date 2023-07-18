// Code generated by go-swagger; DO NOT EDIT.

package metrics_token

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	apiserver_params "github.com/cloudbase/garm/apiserver/params"
	garm_params "github.com/cloudbase/garm/params"
)

// MetricsTokenReader is a Reader for the MetricsToken structure.
type MetricsTokenReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *MetricsTokenReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewMetricsTokenOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewMetricsTokenUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /metrics-token] MetricsToken", response, response.Code())
	}
}

// NewMetricsTokenOK creates a MetricsTokenOK with default headers values
func NewMetricsTokenOK() *MetricsTokenOK {
	return &MetricsTokenOK{}
}

/*
MetricsTokenOK describes a response with status code 200, with default header values.

JWTResponse
*/
type MetricsTokenOK struct {
	Payload garm_params.JWTResponse
}

// IsSuccess returns true when this metrics token o k response has a 2xx status code
func (o *MetricsTokenOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this metrics token o k response has a 3xx status code
func (o *MetricsTokenOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this metrics token o k response has a 4xx status code
func (o *MetricsTokenOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this metrics token o k response has a 5xx status code
func (o *MetricsTokenOK) IsServerError() bool {
	return false
}

// IsCode returns true when this metrics token o k response a status code equal to that given
func (o *MetricsTokenOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the metrics token o k response
func (o *MetricsTokenOK) Code() int {
	return 200
}

func (o *MetricsTokenOK) Error() string {
	return fmt.Sprintf("[GET /metrics-token][%d] metricsTokenOK  %+v", 200, o.Payload)
}

func (o *MetricsTokenOK) String() string {
	return fmt.Sprintf("[GET /metrics-token][%d] metricsTokenOK  %+v", 200, o.Payload)
}

func (o *MetricsTokenOK) GetPayload() garm_params.JWTResponse {
	return o.Payload
}

func (o *MetricsTokenOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewMetricsTokenUnauthorized creates a MetricsTokenUnauthorized with default headers values
func NewMetricsTokenUnauthorized() *MetricsTokenUnauthorized {
	return &MetricsTokenUnauthorized{}
}

/*
MetricsTokenUnauthorized describes a response with status code 401, with default header values.

APIErrorResponse
*/
type MetricsTokenUnauthorized struct {
	Payload apiserver_params.APIErrorResponse
}

// IsSuccess returns true when this metrics token unauthorized response has a 2xx status code
func (o *MetricsTokenUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this metrics token unauthorized response has a 3xx status code
func (o *MetricsTokenUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this metrics token unauthorized response has a 4xx status code
func (o *MetricsTokenUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this metrics token unauthorized response has a 5xx status code
func (o *MetricsTokenUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this metrics token unauthorized response a status code equal to that given
func (o *MetricsTokenUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the metrics token unauthorized response
func (o *MetricsTokenUnauthorized) Code() int {
	return 401
}

func (o *MetricsTokenUnauthorized) Error() string {
	return fmt.Sprintf("[GET /metrics-token][%d] metricsTokenUnauthorized  %+v", 401, o.Payload)
}

func (o *MetricsTokenUnauthorized) String() string {
	return fmt.Sprintf("[GET /metrics-token][%d] metricsTokenUnauthorized  %+v", 401, o.Payload)
}

func (o *MetricsTokenUnauthorized) GetPayload() apiserver_params.APIErrorResponse {
	return o.Payload
}

func (o *MetricsTokenUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}