resource "nuon_app" "my_app" {
  name = "My App"
}

resource "nuon_container_image_component" "my_component" {
  name   = "My Component"
  app_id = nuon_app.my_app.id

  public = {
    image_url = "kennethreitz/httpbin"
    tag       = "latest"
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
