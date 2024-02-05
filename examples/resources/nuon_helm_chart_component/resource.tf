resource "nuon_app" "my_app" {
  name = "my_app"
}

data "nuon_connected_repo" "my_repo" {
  name = "my-github-org/my-repo"
}

resource "nuon_helm_chart_component" "my_component" {
  name   = "my_component"
  app_id = nuon_app.my_app.id

  chart_name = "my-helm-chart"

  connected_repo = {
    directory = "example/chart"
    repo      = nuon_connected_repo.my_repo.name
    branch    = nuon_connected_repo.my_repo.default_branch
  }

  # manually set a variable
  value {
    name  = "some-var-name"
    value = "some-var-value"
  }

  # reference another component
  value {
    name  = "reference-to-other-component"
    value = "{{.nuon.components.some_other_component.outputs.s3_bucket_name}}"
  }

  # reference the install
  value {
    name  = "reference-to-install-attribute"
    value = "{{.nuon.install.public_domain}}"
  }
}
