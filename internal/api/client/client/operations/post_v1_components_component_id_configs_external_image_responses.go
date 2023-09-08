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

// PostV1ComponentsComponentIDConfigsExternalImageReader is a Reader for the PostV1ComponentsComponentIDConfigsExternalImage structure.
type PostV1ComponentsComponentIDConfigsExternalImageReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostV1ComponentsComponentIDConfigsExternalImageReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostV1ComponentsComponentIDConfigsExternalImageCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[POST /v1/components/{component_id}/configs/external-image] PostV1ComponentsComponentIDConfigsExternalImage", response, response.Code())
	}
}

// NewPostV1ComponentsComponentIDConfigsExternalImageCreated creates a PostV1ComponentsComponentIDConfigsExternalImageCreated with default headers values
func NewPostV1ComponentsComponentIDConfigsExternalImageCreated() *PostV1ComponentsComponentIDConfigsExternalImageCreated {
	return &PostV1ComponentsComponentIDConfigsExternalImageCreated{}
}

/*
PostV1ComponentsComponentIDConfigsExternalImageCreated describes a response with status code 201, with default header values.

Created
*/
type PostV1ComponentsComponentIDConfigsExternalImageCreated struct {
	Payload *models.AppExternalImageComponentConfig
}

// IsSuccess returns true when this post v1 components component Id configs external image created response has a 2xx status code
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post v1 components component Id configs external image created response has a 3xx status code
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post v1 components component Id configs external image created response has a 4xx status code
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this post v1 components component Id configs external image created response has a 5xx status code
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this post v1 components component Id configs external image created response a status code equal to that given
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) IsCode(code int) bool {
	return code == 201
}

// Code gets the status code for the post v1 components component Id configs external image created response
func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) Code() int {
	return 201
}

func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) Error() string {
	return fmt.Sprintf("[POST /v1/components/{component_id}/configs/external-image][%d] postV1ComponentsComponentIdConfigsExternalImageCreated  %+v", 201, o.Payload)
}

func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) String() string {
	return fmt.Sprintf("[POST /v1/components/{component_id}/configs/external-image][%d] postV1ComponentsComponentIdConfigsExternalImageCreated  %+v", 201, o.Payload)
}

func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) GetPayload() *models.AppExternalImageComponentConfig {
	return o.Payload
}

func (o *PostV1ComponentsComponentIDConfigsExternalImageCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AppExternalImageComponentConfig)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
