resource "terraform_data" "example" {
  input = "attach-vps-public-keys"

  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.hostinger_vps_public_keys_attach.example]
    }
  }
}

action "hostinger_vps_public_keys_attach" "example" {
  config {
    public_key_ids     = [65224, 65225]
    virtual_machine_id = 12345
    wait_for_action    = true
  }
}
