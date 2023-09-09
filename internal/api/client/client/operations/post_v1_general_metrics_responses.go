// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// PostV1GeneralMetricsReader is a Reader for the PostV1GeneralMetrics structure.
type PostV1GeneralMetricsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostV1GeneralMetricsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostV1GeneralMetricsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[POST /v1/general/metrics] PostV1GeneralMetrics", response, response.Code())
	}
}

// NewPostV1GeneralMetricsOK creates a PostV1GeneralMetricsOK with default headers values
func NewPostV1GeneralMetricsOK() *PostV1GeneralMetricsOK {
	return &PostV1GeneralMetricsOK{}
}

/*
PostV1GeneralMetricsOK describes a response with status code 200, with default header values.

OK
*/
type PostV1GeneralMetricsOK struct {
	Payload string
}

// IsSuccess returns true when this post v1 general metrics o k response has a 2xx status code
func (o *PostV1GeneralMetricsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post v1 general metrics o k response has a 3xx status code
func (o *PostV1GeneralMetricsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post v1 general metrics o k response has a 4xx status code
func (o *PostV1GeneralMetricsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this post v1 general metrics o k response has a 5xx status code
func (o *PostV1GeneralMetricsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this post v1 general metrics o k response a status code equal to that given
func (o *PostV1GeneralMetricsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the post v1 general metrics o k response
func (o *PostV1GeneralMetricsOK) Code() int {
	return 200
}

func (o *PostV1GeneralMetricsOK) Error() string {
	return fmt.Sprintf("[POST /v1/general/metrics][%d] postV1GeneralMetricsOK  %+v", 200, o.Payload)
}

func (o *PostV1GeneralMetricsOK) String() string {
	return fmt.Sprintf("[POST /v1/general/metrics][%d] postV1GeneralMetricsOK  %+v", 200, o.Payload)
}

func (o *PostV1GeneralMetricsOK) GetPayload() string {
	return o.Payload
}

func (o *PostV1GeneralMetricsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
