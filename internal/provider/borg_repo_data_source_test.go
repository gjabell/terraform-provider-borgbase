package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBorgRepoDataSource(t *testing.T) {
	name := "terraform_test"

	id := "borgbase_borg_repo.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBorgRepoDataSourceConfig(name),
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
					resource.TestCheckResourceAttr(id, "region", "eu"),
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
		},
	})
}

func testAccBorgRepoDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "borgbase_borg_repo" "test" {
	name = %[1]q
	region = "eu"
}

data "borgbase_borg_repo" "test" {
	name = borgbase_borg_repo.test.name
}`, name)
}
