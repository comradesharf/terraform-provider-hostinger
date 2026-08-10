resource "terraform_data" "example" {
  input = "start-vps-virtual-machine"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.hostinger_vps_virtual_machine_start.example]
    }
  }
}

action "hostinger_vps_virtual_machine_start" "example" {
  config {
    virtual_machine_id = 12345
    wait_for_action    = true
  }
}
