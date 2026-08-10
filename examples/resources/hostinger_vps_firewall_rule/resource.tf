resource "hostinger_vps_firewall" "example" {
  name = "My Firewall Group"
}

resource "hostinger_vps_firewall_rule" "example" {
  firewall_id   = hostinger_vps_firewall.example.id
  protocol      = "TCP"
  port          = "80"
  source        = "custom"
  source_detail = "any"
}
