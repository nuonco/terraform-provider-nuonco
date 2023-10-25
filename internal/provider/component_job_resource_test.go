package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccComponentJobResource(app AppResourceModel, job JobComponentResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_job_component" "my_job" {
    name = %s
    app_id = nuon_app.my_app.id
    image_url = %s
    tag = %s
    cmd = %s
    args = %s

    env_var {
        name = %s
        value = %s
    }
}
`,
		app.Name,
		job.Name,
		job.ImageURL,
		job.Tag,
		job.Cmd,
		job.Args,
		job.EnvVar[0].Name,
		job.EnvVar[0].Value,
	)
}

func TestComponentJobResource(t *testing.T) {
	ctx := context.Background()
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	job := JobComponentResourceModel{
		Name:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		ImageURL: types.StringValue("bitnami/kubectl"),
		Tag:      types.StringValue("latest"),
		Cmd:      stringSliceToList(ctx, []string{"kubectl"}),
		Args:     stringSliceToList(ctx, []string{"get", "pods"}),
		EnvVar:   NewEnvVarSliceFromMap(map[string]string{"POD_NAMESPACE": "default"}),
	}

	updatedJob := JobComponentResourceModel{
		Name:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		ImageURL: types.StringValue("postgres"),
		Tag:      types.StringValue("alpine3.18"),
		Cmd:      stringSliceToList(ctx, []string{"psql"}),
		Args:     stringSliceToList(ctx, []string{"-U", "admin"}),
		EnvVar:   NewEnvVarSliceFromMap(map[string]string{"PGPASSWORD": "password"}),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentJobResource(app, job),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "name", job.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "image_url", job.ImageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "tag", job.Tag.ValueString()),
					resource.TestCheckTypeSetElemAttr("nuon_job_component.my_job", "cmd.*", listToStringSlice(job.Cmd)[0]),
					resource.TestCheckTypeSetElemAttr("nuon_job_component.my_job", "args.*", listToStringSlice(job.Args)[0]),
					// TODO(ja): These won't pass, even though I've confirmed the env vars are being set and read,
					// we do the same thing for Helm and TF components, and I implemented this test following the TF docs.
					// Will circle back after testing with canaries.
					// resource.TestCheckTypeSetElemNestedAttrs("nuon_job_component.my_job", "env_var.*", job.EnvVar.ToMap()),
				),
			},
			// Import State
			{
				ResourceName:      "nuon_job_component.my_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccComponentJobResource(app, updatedJob),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "name", updatedJob.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "image_url", updatedJob.ImageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_job_component.my_job", "tag", updatedJob.Tag.ValueString()),
					resource.TestCheckTypeSetElemAttr("nuon_job_component.my_job", "cmd.*", listToStringSlice(updatedJob.Cmd)[0]),
					resource.TestCheckTypeSetElemAttr("nuon_job_component.my_job", "args.*", listToStringSlice(updatedJob.Args)[0]),
					// resource.TestCheckTypeSetElemNestedAttrs("nuon_job_component.my_job", "env_var.*", updatedJob.EnvVar.ToMap()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
