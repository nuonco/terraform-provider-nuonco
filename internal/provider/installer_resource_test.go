package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccAppInstallerResource(app AppResourceModel, appInstaller AppInstallerResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_app_installer" "my_app_installer" {
    app_ids = [nuon_app.my_app.id]
    name = %s
    community_url = %s
    description = %s
    documentation_url = %s
    github_url = %s
    homepage_url = %s
    logo_url = %s
    slug = %s
}
`,
		app.Name,

		appInstaller.Name,
		appInstaller.CommunityURL,
		appInstaller.Description,
		appInstaller.DocumentationURL,
		appInstaller.GithubURL,
		appInstaller.HomepageURL,
		appInstaller.LogoURL,
		appInstaller.Slug,
	)
}

func TestAppInstallerResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	appInstaller := AppInstallerResourceModel{
		AppID:            app.Id,
		Name:             types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		CommunityURL:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Description:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		DocumentationURL: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		GithubURL:        types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		HomepageURL:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		LogoURL:          types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Slug:             types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}

	updatedAppInstaller := AppInstallerResourceModel{
		AppID:            app.Id,
		Name:             types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		CommunityURL:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Description:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		DocumentationURL: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		GithubURL:        types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		HomepageURL:      types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		LogoURL:          types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		Slug:             types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppInstallerResource(app, appInstaller),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_install.my_install", "name", appInstaller.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "community_url", appInstaller.CommunityURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "description", appInstaller.Description.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "documentation_url", appInstaller.DocumentationURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "github_url", appInstaller.GithubURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "homepage_url", appInstaller.HomepageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "logo_url", appInstaller.LogoURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "slug", appInstaller.Slug.ValueString()),
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
				Config: testAccAppInstallerResource(app, updatedAppInstaller),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_install.my_install", "name", updatedAppInstaller.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "community_url", updatedAppInstaller.CommunityURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "description", updatedAppInstaller.Description.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "documentation_url", updatedAppInstaller.DocumentationURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "github_url", updatedAppInstaller.GithubURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "homepage_url", updatedAppInstaller.HomepageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "logo_url", updatedAppInstaller.LogoURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_install.my_install", "slug", updatedAppInstaller.Slug.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
