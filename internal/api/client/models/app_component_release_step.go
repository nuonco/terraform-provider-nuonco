// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AppComponentReleaseStep app component release step
//
// swagger:model app.ComponentReleaseStep
type AppComponentReleaseStep struct {

	// component release
	ComponentRelease *AppComponentRelease `json:"componentRelease,omitempty"`

	// parent release ID
	ComponentReleaseID string `json:"component_release_id,omitempty"`

	// created at
	CreatedAt string `json:"created_at,omitempty"`

	// created by id
	CreatedByID string `json:"created_by_id,omitempty"`

	// fields to control the delay of the individual step, as this is set based on the parent strategy
	Delay string `json:"delay,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// install deploys
	InstallDeploys []*AppInstallDeploy `json:"install_deploys"`

	// When a step is created, a set of installs are targeted. However, by the time the release step goes out, the
	// install might have been setup in any order of ways.
	RequestedInstallIds []string `json:"requested_install_ids"`

	// status
	Status string `json:"status,omitempty"`

	// status description
	StatusDescription string `json:"status_description,omitempty"`

	// updated at
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Validate validates this app component release step
func (m *AppComponentReleaseStep) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComponentRelease(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInstallDeploys(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AppComponentReleaseStep) validateComponentRelease(formats strfmt.Registry) error {
	if swag.IsZero(m.ComponentRelease) { // not required
		return nil
	}

	if m.ComponentRelease != nil {
		if err := m.ComponentRelease.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("componentRelease")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("componentRelease")
			}
			return err
		}
	}

	return nil
}

func (m *AppComponentReleaseStep) validateInstallDeploys(formats strfmt.Registry) error {
	if swag.IsZero(m.InstallDeploys) { // not required
		return nil
	}

	for i := 0; i < len(m.InstallDeploys); i++ {
		if swag.IsZero(m.InstallDeploys[i]) { // not required
			continue
		}

		if m.InstallDeploys[i] != nil {
			if err := m.InstallDeploys[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("install_deploys" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("install_deploys" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this app component release step based on the context it is used
func (m *AppComponentReleaseStep) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateComponentRelease(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstallDeploys(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AppComponentReleaseStep) contextValidateComponentRelease(ctx context.Context, formats strfmt.Registry) error {

	if m.ComponentRelease != nil {

		if swag.IsZero(m.ComponentRelease) { // not required
			return nil
		}

		if err := m.ComponentRelease.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("componentRelease")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("componentRelease")
			}
			return err
		}
	}

	return nil
}

func (m *AppComponentReleaseStep) contextValidateInstallDeploys(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.InstallDeploys); i++ {

		if m.InstallDeploys[i] != nil {

			if swag.IsZero(m.InstallDeploys[i]) { // not required
				return nil
			}

			if err := m.InstallDeploys[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("install_deploys" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("install_deploys" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *AppComponentReleaseStep) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AppComponentReleaseStep) UnmarshalBinary(b []byte) error {
	var res AppComponentReleaseStep
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
