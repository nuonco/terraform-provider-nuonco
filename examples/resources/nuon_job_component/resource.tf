resource "nuon_app" "my_app" {
  name = "my_app"
}

resource "nuon_job_component" "job" {
  name   = "my_job"
  app_id = nuon_app.my_app.id

  image_url = "bitnami/kubectl"
  tag       = "latest"
  cmd       = ["kubectl"]
  args      = ["get", "pods", "-A"]

  env_var {
    name  = "NUON_APP_ID"
    value = "{{.nuon.app.id}}"
  }
}
