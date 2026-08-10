resource "hostinger_vps_virtual_machine_hostname" "example" {
  virtual_machine_id = 12345
  hostname           = "server.example.com"
  wait_for_action    = true
}
