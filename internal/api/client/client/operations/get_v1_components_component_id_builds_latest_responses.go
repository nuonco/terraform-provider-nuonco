// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/nuonco/terraform-provider-nuon/internal/api/client/models"
)

// GetV1ComponentsComponentIDBuildsLatestReader is a Reader for the GetV1ComponentsComponentIDBuildsLatest structure.
type GetV1ComponentsComponentIDBuildsLatestReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1ComponentsComponentIDBuildsLatestReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1ComponentsComponentIDBuildsLatestOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/components/{component_id}/builds/latest] GetV1ComponentsComponentIDBuildsLatest", response, response.Code())
	}
}

// NewGetV1ComponentsComponentIDBuildsLatestOK creates a GetV1ComponentsComponentIDBuildsLatestOK with default headers values
func NewGetV1ComponentsComponentIDBuildsLatestOK() *GetV1ComponentsComponentIDBuildsLatestOK {
	return &GetV1ComponentsComponentIDBuildsLatestOK{}
}

/*
GetV1ComponentsComponentIDBuildsLatestOK describes a response with status code 200, with default header values.

OK
*/
type GetV1ComponentsComponentIDBuildsLatestOK struct {
	Payload *models.AppComponentBuild
}

// IsSuccess returns true when this get v1 components component Id builds latest o k response has a 2xx status code
func (o *GetV1ComponentsComponentIDBuildsLatestOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 components component Id builds latest o k response has a 3xx status code
func (o *GetV1ComponentsComponentIDBuildsLatestOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 components component Id builds latest o k response has a 4xx status code
func (o *GetV1ComponentsComponentIDBuildsLatestOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 components component Id builds latest o k response has a 5xx status code
func (o *GetV1ComponentsComponentIDBuildsLatestOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 components component Id builds latest o k response a status code equal to that given
func (o *GetV1ComponentsComponentIDBuildsLatestOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get v1 components component Id builds latest o k response
func (o *GetV1ComponentsComponentIDBuildsLatestOK) Code() int {
	return 200
}

func (o *GetV1ComponentsComponentIDBuildsLatestOK) Error() string {
	return fmt.Sprintf("[GET /v1/components/{component_id}/builds/latest][%d] getV1ComponentsComponentIdBuildsLatestOK  %+v", 200, o.Payload)
}

func (o *GetV1ComponentsComponentIDBuildsLatestOK) String() string {
	return fmt.Sprintf("[GET /v1/components/{component_id}/builds/latest][%d] getV1ComponentsComponentIdBuildsLatestOK  %+v", 200, o.Payload)
}

func (o *GetV1ComponentsComponentIDBuildsLatestOK) GetPayload() *models.AppComponentBuild {
	return o.Payload
}

func (o *GetV1ComponentsComponentIDBuildsLatestOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AppComponentBuild)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
