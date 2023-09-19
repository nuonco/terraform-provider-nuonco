resource "nuon_app" "my_app" {
  name = "My App"
}

resource "nuon_install" "customer_one" {
  app_id = nuon_app.main.id

  name         = "Customer One"
  region       = "us-east-1"
  iam_role_arn = var.customer_one_install_role
}

resource "nuon_install" "customer_two" {
  app_id = nuon_app.main.id

  name         = "Customer Two"
  region       = "us-west-2"
  iam_role_arn = var.customer_two_install_role
}

resource "nuon_install" "customer_three" {
  app_id = nuon_app.main.id

  name         = "Customer Three"
  region       = "us-west-2"
  iam_role_arn = var.customer_three_install_role
}
