resource "terraform_data" "example" {
  input = "sync-vps-firewall"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.hostinger_vps_firewall_sync.example]
    }
  }
}

action "hostinger_vps_firewall_sync" "example" {
  config {
    firewall_id        = 65224
    virtual_machine_id = 12345
    wait_for_action    = true
  }
}
