resource "hostinger_vps_virtual_machine_nameservers" "example" {
  virtual_machine_id = 12345
  ns1                = "1.1.1.1"
  ns2                = "8.8.8.8"
  wait_for_action    = true
}
