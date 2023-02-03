data "borgbase_ssh_key" "append_only" {
  name = "append_only"
}

data "borgbase_ssh_key" "full_access" {
  name = "full_access"
}

data "borgbase_ssh_key" "rsync" {
  name = "rsync"
}

resource "borgbase_borg_repo" "repo_minimal" {
  name   = "repo_minimal"
  region = "eu"
}

resource "borgbase_borg_repo" "repo_full" {
  alert_days       = 2
  append_only      = true
  append_only_keys = [data.borgbase_ssh_key.append_only]
  borg_version     = "LATEST"
  compaction = {
    enabled          = true
    hour             = 14
    hour_timezone    = "Europe/Berlin"
    interval         = 6
    interval_unit    = "weeks"
    full_access_keys = [data.borgbase_ssh_key.full_access]
  }
  name          = "repo_full"
  quota         = 10000
  quota_enabled = true
  region        = "eu"
  rsync_keys    = [data.borgbase_ssh_key.rsync]
  sftp_enabled  = true
}
