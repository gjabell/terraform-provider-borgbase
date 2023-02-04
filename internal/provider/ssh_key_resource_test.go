package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSshKeyResource(t *testing.T) {
	name := "terraform_test"
	publicKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost"

	id := "borgbase_ssh_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSshKeyResourceConfig(name, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(
						id,
						"added_at",
						func(v string) error {
							_, err := time.Parse(time.RFC3339, v)
							return err
						},
					),
					resource.TestCheckResourceAttr(id, "bits", "256"),
					resource.TestCheckResourceAttr(
						id,
						"hash_md5",
						"55:62:b1:68:e7:d9:2f:66:ff:29:b6:fb:41:b6:39:a9",
					),
					resource.TestCheckResourceAttr(
						id,
						"hash_sha256",
						"pZlnOMnSYab3A2b1GDfSXBHR1wKEp8RflbcGXsC6la8",
					),
					resource.TestMatchResourceAttr(id, "id", regexp.MustCompile(`\d+`)),
					resource.TestCheckResourceAttr(id, "last_used_at", ""),
					resource.TestCheckResourceAttr(id, "name", name),
					resource.TestCheckResourceAttr(id, "public_key", publicKey),
					resource.TestCheckResourceAttr(id, "type", "ssh-ed25519"),
				),
			},
			// ImportState testing
			{
				ResourceName:      id,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     name,
			},
			// Update and Read testing
			{
				Config: testAccSshKeyResourceConfig(name+"_new", publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(
						id,
						"added_at",
						func(v string) error {
							_, err := time.Parse(time.RFC3339, v)
							return err
						},
					),
					resource.TestCheckResourceAttr(id, "bits", "256"),
					resource.TestCheckResourceAttr(
						id,
						"hash_md5",
						"55:62:b1:68:e7:d9:2f:66:ff:29:b6:fb:41:b6:39:a9",
					),
					resource.TestCheckResourceAttr(
						id,
						"hash_sha256",
						"pZlnOMnSYab3A2b1GDfSXBHR1wKEp8RflbcGXsC6la8",
					),
					resource.TestMatchResourceAttr(id, "id", regexp.MustCompile(`\d+`)),
					resource.TestCheckResourceAttr(id, "last_used_at", ""),
					resource.TestCheckResourceAttr(id, "name", name+"_new"),
					resource.TestCheckResourceAttr(id, "public_key", publicKey),
					resource.TestCheckResourceAttr(id, "type", "ssh-ed25519"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSshKeyResourceConfig(name, publicKey string) string {
	return fmt.Sprintf(`
resource "borgbase_ssh_key" "test" {
	name = %q
	public_key = %q
}`, name, publicKey)
}
