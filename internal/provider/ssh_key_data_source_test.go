package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSshKeyDataSource(t *testing.T) {
	name := "terraform_test"

	id := "borgbase_ssh_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSshKeyDataSourceConfig(name),
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
					resource.TestCheckResourceAttr(
						id,
						"public_key",
						"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost",
					),
					resource.TestCheckResourceAttr(id, "type", "ssh-ed25519"),
				),
			},
		},
	})
}

func testAccSshKeyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "borgbase_ssh_key" "test" {
	name = %[1]q
	public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost"
}

data "borgbase_ssh_key" "test" {
	name = borgbase_ssh_key.test.name
}`, name)
}
