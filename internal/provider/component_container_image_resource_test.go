package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccComponentContainerImageResource(app AppResourceModel, component ContainerImageComponentResourceModel) string {
	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_container_image_component" "my_component" {
    app_id = nuon_app.my_app.id
    name = %s

    public = {
	image_url = %s
	tag = %s
    }

}
`,
		app.Name,
		component.Name,
		component.Public.ImageURL,
		component.Public.Tag,
	)
}

func TestComponentContainerImageResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}
	component := ContainerImageComponentResourceModel{
		Name:   types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		AwsEcr: nil,
		Public: &Public{
			ImageURL: types.StringValue("kennethreitz/httpbin"),
			Tag:      types.StringValue("latest"),
		},
		EnvVar: []EnvVar{},
	}

	updatedComponent := ContainerImageComponentResourceModel{
		Name:   types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
		AwsEcr: component.AwsEcr,
		Public: component.Public,
		EnvVar: component.EnvVar,
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentContainerImageResource(app, component),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "name", component.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "public.image_url", component.Public.ImageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "public.tag", component.Public.Tag.ValueString()),
				),
			},
			// Import State
			{
				ResourceName:      "nuon_container_image_component.my_component",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccComponentContainerImageResource(app, updatedComponent),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "name", updatedComponent.Name.ValueString()),
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "public.image_url", updatedComponent.Public.ImageURL.ValueString()),
					resource.TestCheckResourceAttr("nuon_container_image_component.my_component", "public.tag", updatedComponent.Public.Tag.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
