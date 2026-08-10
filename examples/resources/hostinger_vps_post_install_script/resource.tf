resource "hostinger_vps_post_install_script" "example" {
  name = "Install monitoring agent"

  content = <<-EOT
    #!/bin/bash
    apt-get update
    apt-get install -y curl
  EOT
}
