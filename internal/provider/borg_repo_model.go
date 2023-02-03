package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BorgRepoModel struct {
	AlertDays      types.Int64   `tfsdk:"alert_days"`
	AppendOnly     types.Bool    `tfsdk:"append_only"`
	AppendOnlyKeys types.List    `tfsdk:"append_only_keys"`
	BorgVersion    types.String  `tfsdk:"borg_version"`
	Compaction     types.Object  `tfsdk:"compaction"`
	CreatedAt      types.String  `tfsdk:"created_at"`
	CurrentUsage   types.Float64 `tfsdk:"current_usage"`
	Encryption     types.String  `tfsdk:"encryption"`
	Format         types.String  `tfsdk:"format"`
	FullAccessKeys types.List    `tfsdk:"full_access_keys"`
	Id             types.String  `tfsdk:"id"`
	LastModified   types.String  `tfsdk:"last_modified"`
	Name           types.String  `tfsdk:"name"`
	Quota          types.Int64   `tfsdk:"quota"`
	QuotaEnabled   types.Bool    `tfsdk:"quota_enabled"`
	Region         types.String  `tfsdk:"region"`
	RepoPath       types.String  `tfsdk:"repo_path"`
	RsyncKeys      types.List    `tfsdk:"rsync_keys"`
	Server         types.Object  `tfsdk:"server"`
	SftpEnabled    types.Bool    `tfsdk:"sftp_enabled"`
}

func (m *BorgRepoModel) update(
	ctx context.Context,
	repo BorgRepoPayload,
) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.AlertDays = types.Int64Value(int64(repo.AlertDays))
	m.AppendOnly = types.BoolValue(repo.AppendOnly)
	m.AppendOnlyKeys, diagnostics = types.ListValueFrom(
		ctx,
		types.StringType,
		repo.AppendOnlyKeys,
	)
	if diagnostics.HasError() {
		return diagnostics
	}

	m.BorgVersion = types.StringValue(string(repo.BorgVersion))
	m.Compaction, diagnostics = types.ObjectValueFrom(
		ctx,
		compactionAttributes,
		CompactionModel{
			Enabled:      types.BoolValue(repo.CompactionEnabled),
			Hour:         types.Int64Value(int64(repo.CompactionHour)),
			HourTimezone: types.StringValue(repo.CompactionHourTimezone),
			Interval:     types.Int64Value(int64(repo.CompactionInterval)),
			IntervalUnit: types.StringValue(repo.CompactionIntervalUnit),
		},
	)
	if diagnostics.HasError() {
		return diagnostics
	}

	m.CreatedAt = types.StringValue(repo.CreatedAt)
	m.CurrentUsage = types.Float64Value(repo.CurrentUsage)
	m.Encryption = types.StringValue(repo.Encryption)
	m.Format = types.StringValue(repo.Format)
	m.FullAccessKeys, diagnostics = types.ListValueFrom(
		ctx,
		types.StringType,
		repo.FullAccessKeys,
	)
	if diagnostics.HasError() {
		return diagnostics
	}

	m.Id = types.StringValue(repo.Id)
	m.LastModified = types.StringValue(repo.LastModified)
	m.Name = types.StringValue(repo.Name)
	m.Quota = types.Int64Value(int64(repo.Quota))
	m.QuotaEnabled = types.BoolValue(repo.QuotaEnabled)
	m.Region = types.StringValue(repo.Region)
	m.RepoPath = types.StringValue(repo.RepoPath)
	m.RsyncKeys, diagnostics = types.ListValueFrom(
		ctx,
		types.StringType,
		repo.RsyncKeys,
	)
	if diagnostics.HasError() {
		return diagnostics
	}

	m.Server, diagnostics = types.ObjectValueFrom(
		ctx,
		serverAttributes,
		ServerModel{
			FingerprintEcdsa:   types.StringValue(repo.Server.FingerprintEcdsa),
			FingerprintEd25519: types.StringValue(repo.Server.FingerprintEd25519),
			FingerprintRsa:     types.StringValue(repo.Server.FingerprintRsa),
			Hostname:           types.StringValue(repo.Server.Hostname),
			Id:                 types.StringValue(repo.Server.Id),
			Location:           types.StringValue(repo.Server.Location),
			Public:             types.BoolValue(repo.Server.Public),
			Region:             types.StringValue(repo.Server.Region),
		},
	)
	if diagnostics.HasError() {
		return diagnostics
	}

	m.SftpEnabled = types.BoolValue(repo.SftpEnabled)
	return diagnostics
}

type CompactionModel struct {
	Enabled      types.Bool   `tfsdk:"enabled"`
	Hour         types.Int64  `tfsdk:"hour"`
	HourTimezone types.String `tfsdk:"hour_timezone"`
	Interval     types.Int64  `tfsdk:"interval"`
	IntervalUnit types.String `tfsdk:"interval_unit"`
}

var compactionAttributes = map[string]attr.Type{
	"enabled":       types.BoolType,
	"hour":          types.Int64Type,
	"hour_timezone": types.StringType,
	"interval":      types.Int64Type,
	"interval_unit": types.StringType,
}

type ServerModel struct {
	FingerprintEcdsa   types.String `tfsdk:"fingerprint_ecdsa"`
	FingerprintEd25519 types.String `tfsdk:"fingerprint_ed25519"`
	FingerprintRsa     types.String `tfsdk:"fingerprint_rsa"`
	Hostname           types.String `tfsdk:"hostname"`
	Id                 types.String `tfsdk:"id"`
	Location           types.String `tfsdk:"location"`
	Public             types.Bool   `tfsdk:"public"`
	Region             types.String `tfsdk:"region"`
}

var serverAttributes = map[string]attr.Type{
	"fingerprint_ecdsa":   types.StringType,
	"fingerprint_ed25519": types.StringType,
	"fingerprint_rsa":     types.StringType,
	"hostname":            types.StringType,
	"id":                  types.StringType,
	"location":            types.StringType,
	"public":              types.BoolType,
	"region":              types.StringType,
}

type BorgRepoPayload struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Server struct {
		Id                 string `json:"id"`
		Hostname           string `json:"hostname"`
		Region             string `json:"region"`
		Public             bool   `json:"public"`
		Location           string `json:"location"`
		FingerprintRsa     string `json:"fingerprintRsa"`
		FingerprintEcdsa   string `json:"fingerprintEcdsa"`
		FingerprintEd25519 string `json:"fingerprintEd25519"`
	} `json:"server"`
	Quota                  int      `json:"quota"`
	QuotaEnabled           bool     `json:"quotaEnabled"`
	AlertDays              int      `json:"alertDays"`
	Region                 string   `json:"region"`
	Format                 string   `json:"format"`
	BorgVersion            string   `json:"borgVersion"`
	ResticVersion          string   `json:"resticVersion"`
	Htpasswd               string   `json:"htpasswd"`
	AppendOnly             bool     `json:"appendOnly"`
	AppendOnlyKeys         []string `json:"appendOnlyKeys"`
	FullAccessKeys         []string `json:"fullAccessKeys"`
	RsyncKeys              []string `json:"rsyncKeys"`
	SftpEnabled            bool     `json:"sftpEnabled"`
	Encryption             string   `json:"encryption"`
	CreatedAt              string   `json:"createdAt"`
	LastModified           string   `json:"lastModified"`
	CompactionEnabled      bool     `json:"compactionEnabled"`
	CompactionInterval     int      `json:"compactionInterval"`
	CompactionIntervalUnit string   `json:"compactionIntervalUnit"`
	CompactionHour         int      `json:"compactionHour"`
	CompactionHourTimezone string   `json:"compactionHourTimezone"`
	RepoPath               string   `json:"repoPath"`
	CurrentUsage           float64  `json:"currentUsage"`
}

type BorgReposPayload []BorgRepoPayload

type BorgRepoAddPayload struct {
	RepoAdded BorgRepoPayload `json:"repoAdded"`
}

type BorgRepoEditPayload struct {
	RepoEdited BorgRepoPayload `json:"repoEdited"`
}

type BorgRepoDeletePayload struct {
	Ok bool `json:"ok"`
}
