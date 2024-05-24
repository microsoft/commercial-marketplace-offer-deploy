resource "null_resource" "provider_registration" {
  provisioner "local-exec" {
    # Sleep for 30 seconds after registering. Sometimes the resource provider won't be
    # registered in the current region quick enough
    command     = "az provider register --n ${var.name} --subscription ${var.subscription_id}; sleep 30"
    interpreter = var.platform == "windows" ? ["PowerShell", "-Command"] : null
  }
}
