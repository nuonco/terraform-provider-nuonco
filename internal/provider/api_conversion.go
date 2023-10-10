package provider

const (
	// orgs, apps & installs:
	statusQueued         string = "queued"
	statusProvisioning   string = "provisioning"
	statusActive         string = "active"
	statusDeleteQueued   string = "delete_queued"
	statusDeprovisioning string = "deprovisioning"

	// local statuses, based on api responses
	statusNotFound               string = "not-found"
	statusTemporarilyUnavailable string = "temporarily-unavailable"

	// builds:
	// queued
	statusPlanning string = "planning"
	statusBuilding string = "building"
	// active

	// deploys:
	// queued
	// planning
	statusDeploying string = "deploying"
	// active
)
