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

// GetV1AppsAppIDComponentsReader is a Reader for the GetV1AppsAppIDComponents structure.
type GetV1AppsAppIDComponentsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1AppsAppIDComponentsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1AppsAppIDComponentsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/apps/{app_id}/components] GetV1AppsAppIDComponents", response, response.Code())
	}
}

// NewGetV1AppsAppIDComponentsOK creates a GetV1AppsAppIDComponentsOK with default headers values
func NewGetV1AppsAppIDComponentsOK() *GetV1AppsAppIDComponentsOK {
	return &GetV1AppsAppIDComponentsOK{}
}

/*
GetV1AppsAppIDComponentsOK describes a response with status code 200, with default header values.

OK
*/
type GetV1AppsAppIDComponentsOK struct {
	Payload []*models.AppComponent
}

// IsSuccess returns true when this get v1 apps app Id components o k response has a 2xx status code
func (o *GetV1AppsAppIDComponentsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 apps app Id components o k response has a 3xx status code
func (o *GetV1AppsAppIDComponentsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 apps app Id components o k response has a 4xx status code
func (o *GetV1AppsAppIDComponentsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 apps app Id components o k response has a 5xx status code
func (o *GetV1AppsAppIDComponentsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 apps app Id components o k response a status code equal to that given
func (o *GetV1AppsAppIDComponentsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get v1 apps app Id components o k response
func (o *GetV1AppsAppIDComponentsOK) Code() int {
	return 200
}

func (o *GetV1AppsAppIDComponentsOK) Error() string {
	return fmt.Sprintf("[GET /v1/apps/{app_id}/components][%d] getV1AppsAppIdComponentsOK  %+v", 200, o.Payload)
}

func (o *GetV1AppsAppIDComponentsOK) String() string {
	return fmt.Sprintf("[GET /v1/apps/{app_id}/components][%d] getV1AppsAppIdComponentsOK  %+v", 200, o.Payload)
}

func (o *GetV1AppsAppIDComponentsOK) GetPayload() []*models.AppComponent {
	return o.Payload
}

func (o *GetV1AppsAppIDComponentsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
