resource "hostinger_vps_public_key" "example" {
  name = "example-key"
  key  = file("~/.ssh/id_ed25519.pub")
}
