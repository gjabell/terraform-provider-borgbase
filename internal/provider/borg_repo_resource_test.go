package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBorgRepoResource_minimal(t *testing.T) {
	name := "terraform_test"
	region := "eu"

	id := "borgbase_borg_repo.test_minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBorgRepoResourceConfig_minimal(name, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(id, "alert_days", "0"),
					resource.TestCheckResourceAttr(id, "append_only", "false"),
					resource.TestCheckResourceAttr(id, "append_only_keys.#", "0"),
					resource.TestCheckResourceAttr(id, "borg_version", "LATEST"),
					resource.TestCheckResourceAttr(id, "compaction.enabled", "false"),
					resource.TestCheckResourceAttr(id, "compaction.hour", "14"),
					resource.TestCheckResourceAttr(id, "compaction.hour_timezone", "UTC"),
					resource.TestCheckResourceAttr(id, "compaction.interval", "6"),
					resource.TestCheckResourceAttr(
						id,
						"compaction.interval_unit",
						"weeks",
					),
					resource.TestCheckResourceAttrWith(
						id,
						"created_at",
						func(v string) error {
							_, err := time.Parse(time.RFC3339, v)
							return err
						},
					),
					resource.TestCheckResourceAttr(id, "current_usage", "0"),
					resource.TestCheckResourceAttr(id, "encryption", "none"),
					resource.TestCheckResourceAttr(id, "format", "borg1"),
					resource.TestCheckResourceAttr(id, "full_access_keys.#", "0"),
					resource.TestCheckResourceAttrSet(id, "id"),
					resource.TestCheckResourceAttr(id, "last_modified", ""),
					resource.TestCheckResourceAttr(id, "name", name),
					resource.TestCheckResourceAttr(id, "quota", "0"),
					resource.TestCheckResourceAttr(id, "quota_enabled", "false"),
					resource.TestCheckResourceAttr(id, "region", region),
					resource.TestMatchResourceAttr(
						id,
						"repo_path",
						regexp.MustCompile(`ssh://.*@.*\.repo\.borgbase\.com/\./repo`),
					),
					resource.TestCheckResourceAttr(id, "rsync_keys.#", "0"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_ecdsa"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_ed25519"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_rsa"),
					resource.TestCheckResourceAttrSet(id, "server.hostname"),
					resource.TestCheckResourceAttrSet(id, "server.id"),
					resource.TestCheckResourceAttrSet(id, "server.location"),
					resource.TestCheckResourceAttrSet(id, "server.public"),
					resource.TestCheckResourceAttrSet(id, "server.region"),
					resource.TestCheckResourceAttr(id, "sftp_enabled", "false"),
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
				Config: testAccBorgRepoResourceConfig_minimal(name+"_new", region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(id, "name", name+"_new"),
					resource.TestCheckResourceAttr(id, "region", region),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBorgRepoResource_full(t *testing.T) {
	name := "terraform_test"
	region := "eu"

	id := "borgbase_borg_repo.test_full"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBorgRepoResourceConfig_full(name, region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(id, "alert_days", "2"),
					resource.TestCheckResourceAttr(id, "append_only", "true"),
					resource.TestCheckResourceAttr(id, "append_only_keys.#", "1"),
					resource.TestCheckResourceAttr(id, "borg_version", "LATEST"),
					resource.TestCheckResourceAttr(id, "compaction.enabled", "true"),
					resource.TestCheckResourceAttr(id, "compaction.hour", "14"),
					resource.TestCheckResourceAttr(
						id,
						"compaction.hour_timezone",
						"Europe/Berlin",
					),
					resource.TestCheckResourceAttr(id, "compaction.interval", "6"),
					resource.TestCheckResourceAttr(
						id,
						"compaction.interval_unit",
						"weeks",
					),
					resource.TestCheckResourceAttrWith(
						id,
						"created_at",
						func(v string) error {
							_, err := time.Parse(time.RFC3339, v)
							return err
						},
					),
					resource.TestCheckResourceAttr(id, "current_usage", "0"),
					resource.TestCheckResourceAttr(id, "encryption", "none"),
					resource.TestCheckResourceAttr(id, "format", "borg1"),
					resource.TestCheckResourceAttr(id, "full_access_keys.#", "1"),
					resource.TestCheckResourceAttrSet(id, "id"),
					resource.TestCheckResourceAttr(id, "last_modified", ""),
					resource.TestCheckResourceAttr(id, "name", name),
					resource.TestCheckResourceAttr(id, "quota", "10000"),
					resource.TestCheckResourceAttr(id, "quota_enabled", "true"),
					resource.TestCheckResourceAttr(id, "region", region),
					resource.TestMatchResourceAttr(
						id,
						"repo_path",
						regexp.MustCompile(`ssh://.*@.*\.repo\.borgbase\.com/\./repo`),
					),
					resource.TestCheckResourceAttr(id, "rsync_keys.#", "1"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_ecdsa"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_ed25519"),
					resource.TestCheckResourceAttrSet(id, "server.fingerprint_rsa"),
					resource.TestCheckResourceAttrSet(id, "server.hostname"),
					resource.TestCheckResourceAttrSet(id, "server.id"),
					resource.TestCheckResourceAttrSet(id, "server.location"),
					resource.TestCheckResourceAttrSet(id, "server.public"),
					resource.TestCheckResourceAttrSet(id, "server.region"),
					resource.TestCheckResourceAttr(id, "sftp_enabled", "true"),
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
				Config: testAccBorgRepoResourceConfig_full(name+"_new", region),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(id, "name", name+"_new"),
					resource.TestCheckResourceAttr(id, "region", region),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccBorgRepoResourceConfig_minimal(name, region string) string {
	return fmt.Sprintf(`
resource "borgbase_borg_repo" "test_minimal" {
	name = %q
	region = %q
}`, name, region)
}

func testAccBorgRepoResourceConfig_full(name, region string) string {
	return fmt.Sprintf(`
resource "borgbase_ssh_key" "terraform_append_only" {
  name       = "append_only"
	public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost"
}

resource "borgbase_ssh_key" "terraform_full_access" {
  name       = "full_access"
	public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost"
}

resource "borgbase_ssh_key" "terraform_rsync" {
  name       = "rsync"
	public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBAt/X37WDQ3cNPEVHQBsW3lH7XPeea5rUoeXuhoTkzR terraform@localhost"
}

resource "borgbase_borg_repo" "test_full" {
  alert_days       = 2
  append_only      = true
  append_only_keys = [borgbase_ssh_key.terraform_append_only.id]
  borg_version     = "LATEST"
  compaction = {
    enabled       = true
    hour          = 14
    hour_timezone = "Europe/Berlin"
    interval      = 6
    interval_unit = "weeks"
  }
	full_access_keys = [borgbase_ssh_key.terraform_full_access.id]
	name             = %q
  quota            = 10000
  quota_enabled    = true
	region           = %q
  rsync_keys       = [borgbase_ssh_key.terraform_rsync.id]
  sftp_enabled     = true
}`, name, region)
}
