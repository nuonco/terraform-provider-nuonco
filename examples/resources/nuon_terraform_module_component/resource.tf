resource "nuon_app" "my_app" {
  name = "my_app"
}

data "nuon_connected_repo" "my_repo" {
  name = "my-github-org/my-repo"
}

resource "nuon_terraform_module_component" "my_component" {
  name   = "my_component"
  app_id = nuon_app.my_app.id

  connected_repo = {
    directory = "example/terraform"
    repo      = nuon_connected_repo.my_repo.name
    branch    = nuon_connected_repo.my_repo.default_branch
  }

  # manually set a variable
  var {
    name  = "some-var-name"
    value = "some-var-value"
  }

  # reference another component
  var {
    name  = "reference-to-other-component"
    value = "{{.nuon.components.some_other_component.outputs.s3_bucket_name}}"
  }

  # reference the install
  var {
    name  = "reference-to-install-attribute"
    value = "{{.nuon.install.public_domain}}"
  }
}
