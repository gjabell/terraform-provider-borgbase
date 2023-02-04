package provider

import (
	"context"
	"fmt"

	"github.com/gjabell/terraform-provider-borgbase/gql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &BorgRepoDataSource{}

func NewBorgRepoDataSource() datasource.DataSource {
	return &BorgRepoDataSource{}
}

type BorgRepoDataSource struct {
	client *gql.Client
}

func (d *BorgRepoDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_borg_repo"
}

func (d *BorgRepoDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "BorgBase borg repository.",
		Attributes: map[string]schema.Attribute{
			"alert_days": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of days after which an alert should be triggered if no new backups are made.",
			},
			"append_only": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the repository should allow old data to be deleted.",
			},
			"append_only_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "IDs of SSH keys which are only allowed to append data to the repository.",
			},
			"borg_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Borg version to use for the repository (defaults to latest stable version).",
			},
			"compaction": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Settings for repo compaction.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Whether to enable repository compaction.",
					},
					"hour": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Hour of the day when the repository should be compacted.",
					},
					"hour_timezone": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Timezone of repository compaction hour.",
					},
					"interval": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Repository compaction interval value (1-24).",
					},
					"interval_unit": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Repository compaction interval unit (days, weeks, or months).",
					},
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the repository was created.",
			},
			"current_usage": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Current usage of the repository in megabytes.",
			},
			"encryption": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the repository is encrypted.",
			},
			"format": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Format of the repository.",
			},
			"full_access_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "IDs of SSH keys which have full access to the repository.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal BorgBase repository identifier.",
			},
			"last_modified": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the repository was last modified.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User-defined repository identifier.",
			},
			"quota": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Max allowed size of the repository in megabytes.",
			},
			"quota_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the repository quota should be enabled.",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Region where the repository is hosted (eu or us).",
			},
			"repo_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH path where the repository can be accessed.",
			},
			"rsync_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "IDs of SSH keys which can access the repository via rsync.",
			},
			"server": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Information about the server where the repository is hosted.",
				Attributes: map[string]schema.Attribute{
					"fingerprint_ecdsa": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Fingerprint of the server's ECDSA SSH key.",
					},
					"fingerprint_ed25519": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Fingerprint of the server's ED25519 SSH key.",
					},
					"fingerprint_rsa": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Fingerprint of the server's RSA SSH key.",
					},
					"hostname": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Hostname of the server.",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Internal ID of the server.",
					},
					"location": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Location of the server.",
					},
					"public": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Whether the server is public.",
					},
					"region": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Region in which the server is located.",
					},
				},
			},
			"sftp_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether SFTP access to the repository should be enabled.",
			},
		},
	}
}

func (d *BorgRepoDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *gql.Client, got: %T. "+
					"Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *BorgRepoDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data BorgRepoModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload BorgReposPayload
	args := gql.Arguments{"name": gql.Optional(data.Name.ValueString())}
	if err := d.client.Query("repoList", &payload, args); err != nil {
		resp.Diagnostics.AddError("Failed to read borg repo", err.Error())
		return
	}

	var repo *BorgRepoPayload
	for _, item := range payload {
		if item.Name == data.Name.ValueString() {
			repo = &item
			break
		}
	}
	if repo == nil {
		resp.Diagnostics.AddError("Unknown borg repo", data.Name.String())
		return
	}
	resp.Diagnostics = append(resp.Diagnostics, data.update(ctx, *repo)...)

	tflog.Trace(ctx, "read repo", map[string]interface{}{
		"id":   data.Id,
		"name": data.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
