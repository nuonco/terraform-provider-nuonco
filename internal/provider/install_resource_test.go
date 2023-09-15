package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccInstallResource(app AppResourceModel, install InstallResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_install" "my_install" {
    app_id = nuon_app.my_app.id
    name = %s
    region = %s
    iam_role_arn = %s
}
`,
		app.Name,
		install.Name,
		install.Region,
		install.IAMRoleARN,
	)
}

func TestInstallResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	install := InstallResourceModel{
		AppID:      app.Id,
		Name:       types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Region:     types.StringValue("us-west-2"),
		IAMRoleARN: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}

	updatedInstall := InstallResourceModel{
		AppID:      app.Id,
		Name:       types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Region:     types.StringValue("us-west-2"),
		IAMRoleARN: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccInstallResource(app, install),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_install.my_install", "name", install.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "region", install.Region.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "iam_role_arn", install.IAMRoleARN.ValueString()),
				),
			},
			// Import State
			{
				ResourceName:      "nuon_install.my_install",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccInstallResource(app, updatedInstall),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_install.my_install", "name", updatedInstall.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "region", updatedInstall.Region.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "iam_role_arn", updatedInstall.IAMRoleARN.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
