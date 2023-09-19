resource "nuon_app" "my_app" {
  name = "My App"
}

data "nuon_connected_repo" "my_repo" {
  name = "my-github-org/my-repo"
}

resource "nuon_docker_build_component" "my_component" {
  name   = "My Component"
  app_id = nuon_app.my_app.id

  connected_repo = {
    directory = "example/docker"
    repo      = var.my_repo
    branch    = "main"
  }

  # manually set a variable
  env_var {
    name  = "some-var-name"
    value = "some-var-value"
  }

  # reference another component
  env_var {
    name  = "reference-to-other-component"
    value = "{{.nuon.components.some_other_component.outputs.s3_bucket_name}}"
  }

  # reference the install
  env_var {
    name  = "reference-to-install-attribute"
    value = "{{.nuon.install.public_domain}}"
  }
}
