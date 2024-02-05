resource "nuon_app" "my_app" {
  name = "my_app"
}

resource "nuon_install" "customer_one" {
  app_id = nuon_app.main.id

  name         = "customer_one"
  region       = "us-east-1"
  iam_role_arn = var.customer_one_install_role
}

resource "nuon_install" "customer_two" {
  app_id = nuon_app.main.id

  name         = "customer_two"
  region       = "us-west-2"
  iam_role_arn = var.customer_two_install_role
}

resource "nuon_install" "customer_three" {
  app_id = nuon_app.main.id

  name         = "customer_three"
  region       = "us-west-2"
  iam_role_arn = var.customer_three_install_role
}
