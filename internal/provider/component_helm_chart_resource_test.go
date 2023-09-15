package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccComponentHelmChartResource(app AppResourceModel, component HelmChartComponentResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_helm_chart_component" "my_component" {
    app_id = nuon_app.my_app.id
    name = %s
    chart_name = %s

    public_repo = {
        repo = %s
        branch = %s
        directory = %s
    }
}
`,
		app.Name,
		component.Name,
		component.ChartName,
		component.PublicRepo.Repo,
		component.PublicRepo.Branch,
		component.PublicRepo.Directory,
	)
}

func TestComponentHelmChartResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	component := HelmChartComponentResourceModel{
		Name:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		ChartName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		PublicRepo: &PublicRepo{
			Repo:      types.StringValue("my-github-org/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
			Branch:    types.StringValue("foobar"),
			Directory: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		},
	}

	updatedComponent := HelmChartComponentResourceModel{
		Name:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		ChartName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		PublicRepo: &PublicRepo{
			Repo:      types.StringValue("my-github-org/" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
			Branch:    types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
			Directory: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentHelmChartResource(app, component),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "name", component.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "chart_name", component.ChartName.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.repo", component.PublicRepo.Repo.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.branch", component.PublicRepo.Branch.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.directory", component.PublicRepo.Directory.ValueString()),
				),
			},
			// Import State
			{
				ResourceName:      "nuon_helm_chart_component.my_component",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccComponentHelmChartResource(app, updatedComponent),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "name", updatedComponent.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "chart_name", updatedComponent.ChartName.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.repo", updatedComponent.PublicRepo.Repo.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.branch", updatedComponent.PublicRepo.Branch.ValueString()),
					resource.TestCheckResourceAttr("nuon_helm_chart_component.my_component", "public_repo.directory", updatedComponent.PublicRepo.Directory.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
