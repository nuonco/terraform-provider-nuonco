package provider

const (
	// orgs, apps & installs:
	statusQueued         = "queued"
	statusProvisioning   = "provisioning"
	statusActive         = "active"
	statusDeleteQueued   = "delete_queued"
	statusDeprovisioning = "deprovisioning"

	// local statuses, based on api responses
	statusNotFound               = "not-found"
	statusTemporarilyUnavailable = "temporarily-unavailable"

	// builds:
	// queued
	statusPlanning = "planning"
	statusBuilding = "building"
	// active

	// deploys:
	// queued
	// planning
	statusDeploying = "deploying"
	// active
)
