resource "borgbase_ssh_key" "example" {
  name       = "example"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR user@hostname"
}
