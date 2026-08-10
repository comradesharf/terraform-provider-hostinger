resource "terraform_data" "example" {
  input = "recreate-vps-virtual-machine"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.hostinger_vps_virtual_machine_recreate.example]
    }
  }
}

action "hostinger_vps_virtual_machine_recreate" "example" {
  config {
    virtual_machine_id     = 12345
    template_id            = 67890
    password               = "root-password"
    panel_password         = "panel-password"
    post_install_script_id = 24680
    wait_for_action        = true
  }
}
