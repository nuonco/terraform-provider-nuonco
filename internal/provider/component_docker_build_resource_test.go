package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccComponentDockerBuildResource(app AppResourceModel, component DockerBuildComponentResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_docker_build_component" "my_component" {
    app_id = nuon_app.my_app.id
    name = %s
    dockerfile = %s

    public_repo = {
	repo = %s
	directory = %s
	branch = %s
    }
}
`,
		app.Name,
		component.Name,
		component.Dockerfile,
		component.PublicRepo.Repo,
		component.PublicRepo.Directory,
		component.PublicRepo.Branch,
	)
}

func TestComponentDockerBuildResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	component := DockerBuildComponentResourceModel{
		Name:       types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		SyncOnly:   types.BoolValue(true),
		Dockerfile: types.StringValue("Dockerfile"),
		PublicRepo: &PublicRepo{
			Repo:      types.StringValue("https://github.com/postmanlabs/httpbin.git"),
			Directory: types.StringValue("."),
			Branch:    types.StringValue("master"),
		},
		ConnectedRepo: nil,
		EnvVar:        []EnvVar{},
	}

	updatedComponent := DockerBuildComponentResourceModel{
		Name:       types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		SyncOnly:   types.BoolValue(true),
		Dockerfile: types.StringValue("Dockerfile"),
		PublicRepo: &PublicRepo{
			Repo:      types.StringValue("https://github.com/postmanlabs/httpbin.git"),
			Directory: types.StringValue("."),
			Branch:    types.StringValue("master"),
		},
		ConnectedRepo: nil,
		EnvVar:        []EnvVar{},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentDockerBuildResource(app, component),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "name", component.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "sync_only", component.SyncOnly.String()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.repo", component.PublicRepo.Repo.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.directory", component.PublicRepo.Directory.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.branch", component.PublicRepo.Branch.ValueString()),
				),
			},
			// Import State
			{
				ResourceName:      "nuon_docker_build_component.my_component",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccComponentDockerBuildResource(app, updatedComponent),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "name", updatedComponent.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "sync_only", updatedComponent.SyncOnly.String()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.repo", updatedComponent.PublicRepo.Repo.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.directory", updatedComponent.PublicRepo.Directory.ValueString()),
					resource.TestCheckResourceAttr("nuon_docker_build_component.my_component", "public_repo.branch", updatedComponent.PublicRepo.Branch.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
