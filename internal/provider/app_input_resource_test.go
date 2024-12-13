package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccAppInputResource(app AppResourceModel, input AppInputResourceModel) string {
	fmt.Println(app.Name)
	fmt.Println(input.Groups[0].Name)

	return fmt.Sprintf(providerConfig+`
resource "nuon_app" "my_app" {
    name = %s
}

resource "nuon_app_input" "my_app" {
    app_id = nuon_app.my_app.id

    group {
        name = %s
        display_name = %s
        description = %s
    }

    input {
        name = %s
        display_name = %s
        description = %s
        default = %s
        required = %s
        sensitive = %s
        group = %s
    }
}
`,
		app.Name,

		input.Groups[0].Name,
		input.Groups[0].DisplayName,
		input.Groups[0].Description,

		input.Inputs[0].Name,
		input.Inputs[0].DisplayName,
		input.Inputs[0].Description,
		input.Inputs[0].Default,
		input.Inputs[0].Required,
		input.Inputs[0].Sensitive,
		input.Inputs[0].Group,
	)
}

func TestAppInputResource(t *testing.T) {
	app := AppResourceModel{
		Name: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
	}

	groupName := types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	input := AppInputResourceModel{
		Inputs: []AppInput{
			AppInput{
				Name:        types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				DisplayName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Description: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Default:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Required:    types.BoolValue(false),
				Sensitive:   types.BoolValue(false),
				Group:       groupName,
			},
		},
		Groups: []AppInputGroup{
			AppInputGroup{
				Name:        groupName,
				DisplayName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Description: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
			},
		},
	}

	updatedGroupName := types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	updatedInput := AppInputResourceModel{
		Inputs: []AppInput{
			AppInput{
				Name:        types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				DisplayName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Description: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Default:     types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Required:    types.BoolValue(false),
				Sensitive:   types.BoolValue(false),
				Group:       updatedGroupName,
			},
		},
		Groups: []AppInputGroup{
			AppInputGroup{
				Name:        updatedGroupName,
				DisplayName: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
				Description: types.StringValue(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)),
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppInputResource(app, input),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_app_input.my_app.input[0]", "input", input.Inputs[0].Name.ValueString()),
				),
			},
			// ImportState
			{
				ResourceName:      "nuon_app_input.my_app",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccAppInputResource(app, updatedInput),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuon_app_input.my_app.input[0]", "name", updatedInput.Inputs[0].Name.ValueString()),
				),
			},
			// Delete testing will happen automatically.
		},
	})
}
