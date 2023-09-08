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

// PatchV1InstallsInstallIDReader is a Reader for the PatchV1InstallsInstallID structure.
type PatchV1InstallsInstallIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchV1InstallsInstallIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchV1InstallsInstallIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[PATCH /v1/installs/{install_id}] PatchV1InstallsInstallID", response, response.Code())
	}
}

// NewPatchV1InstallsInstallIDOK creates a PatchV1InstallsInstallIDOK with default headers values
func NewPatchV1InstallsInstallIDOK() *PatchV1InstallsInstallIDOK {
	return &PatchV1InstallsInstallIDOK{}
}

/*
PatchV1InstallsInstallIDOK describes a response with status code 200, with default header values.

OK
*/
type PatchV1InstallsInstallIDOK struct {
	Payload *models.AppInstall
}

// IsSuccess returns true when this patch v1 installs install Id o k response has a 2xx status code
func (o *PatchV1InstallsInstallIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this patch v1 installs install Id o k response has a 3xx status code
func (o *PatchV1InstallsInstallIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this patch v1 installs install Id o k response has a 4xx status code
func (o *PatchV1InstallsInstallIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this patch v1 installs install Id o k response has a 5xx status code
func (o *PatchV1InstallsInstallIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this patch v1 installs install Id o k response a status code equal to that given
func (o *PatchV1InstallsInstallIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the patch v1 installs install Id o k response
func (o *PatchV1InstallsInstallIDOK) Code() int {
	return 200
}

func (o *PatchV1InstallsInstallIDOK) Error() string {
	return fmt.Sprintf("[PATCH /v1/installs/{install_id}][%d] patchV1InstallsInstallIdOK  %+v", 200, o.Payload)
}

func (o *PatchV1InstallsInstallIDOK) String() string {
	return fmt.Sprintf("[PATCH /v1/installs/{install_id}][%d] patchV1InstallsInstallIdOK  %+v", 200, o.Payload)
}

func (o *PatchV1InstallsInstallIDOK) GetPayload() *models.AppInstall {
	return o.Payload
}

func (o *PatchV1InstallsInstallIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AppInstall)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
