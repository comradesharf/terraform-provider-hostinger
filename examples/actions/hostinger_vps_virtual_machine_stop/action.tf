resource "terraform_data" "example" {
  input = "stop-vps-virtual-machine"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.hostinger_vps_virtual_machine_stop.example]
    }
  }
}

action "hostinger_vps_virtual_machine_stop" "example" {
  config {
    virtual_machine_id = 12345
    wait_for_action    = true
  }
}
