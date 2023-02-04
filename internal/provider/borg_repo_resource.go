package provider

import (
	"context"
	"fmt"

	"github.com/gjabell/terraform-provider-borgbase/gql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BorgRepoResource{}
var _ resource.ResourceWithImportState = &BorgRepoResource{}

func NewBorgRepoResource() resource.Resource {
	return &BorgRepoResource{}
}

type BorgRepoResource struct {
	client *gql.Client
}

func (r *BorgRepoResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_borg_repo"
}

func (r *BorgRepoResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "BorgBase borg repository.",
		Attributes: map[string]schema.Attribute{
			"alert_days": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Number of days after which an alert should be triggered if no new backups are made.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"append_only": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Whether the repository should allow old data to be deleted.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"append_only_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "IDs of SSH keys which are only allowed to append data to the repository.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"borg_version": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Borg version to use for the repository (defaults to latest stable version).",
			},
			"compaction": schema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Settings for repository compaction.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Whether to enable repository compaction.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"hour": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Hour of the day when the repository should be compacted.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int64{
							int64validator.Between(0, 23),
						},
					},
					"hour_timezone": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Timezone of repository compaction hour.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"interval": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Repository compaction interval value (1-24).",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int64{
							int64validator.Between(1, 24),
						},
					},
					"interval_unit": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Repository compaction interval unit (days, weeks, or months).",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("days", "weeks", "months"),
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the repository was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"current_usage": schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: "Current usage of the repository in megabytes.",
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"encryption": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the repository is encrypted.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"format": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Format of the repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"full_access_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "IDs of SSH keys which have full access to the repository.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal BorgBase repository identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the repository was last modified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User-defined repository identifier.",
			},
			"quota": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Max allowed size of the repository in megabytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"quota_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Whether the repository quota should be enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Region where the repository is hosted (eu or us).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("eu", "us"),
				},
			},
			"repo_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH path where the repository can be accessed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rsync_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "IDs of SSH keys which can access the repository via rsync.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.ObjectAttribute{
				AttributeTypes:      serverAttributes,
				Computed:            true,
				MarkdownDescription: "Information about the server where the repository is hosted.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"sftp_enabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Whether SFTP access to the repository should be enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *BorgRepoResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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
	r.client = client
}

func setArguments(
	ctx context.Context,
	args map[string]gql.Argument,
	data BorgRepoModel,
) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	for name, attr := range map[string]basetypes.Int64Value{
		"alertDays": data.AlertDays,
		"quota":     data.Quota,
	} {
		if !attr.IsNull() || !attr.IsUnknown() {
			args[name] = gql.Optional(attr.ValueInt64())
		}
	}

	for name, attr := range map[string]basetypes.BoolValue{
		"appendOnly":   data.AppendOnly,
		"sftpEnabled":  data.SftpEnabled,
		"quotaEnabled": data.QuotaEnabled,
	} {
		if !attr.IsNull() || !attr.IsUnknown() {
			args[name] = gql.Optional(attr.ValueBool())
		}
	}

	if !data.BorgVersion.IsNull() && !data.BorgVersion.IsUnknown() {
		args["borgVersion"] = gql.Optional(data.BorgVersion.ValueString())
	}

	if !data.Compaction.IsNull() && !data.Compaction.IsUnknown() {
		var compaction CompactionModel
		diagnostics = data.Compaction.As(
			ctx,
			&compaction,
			basetypes.ObjectAsOptions{},
		)
		if diagnostics.HasError() {
			return diagnostics
		}
		args["compactionEnabled"] = gql.Optional(compaction.Enabled.ValueBool())
		args["compactionHour"] = gql.Optional(compaction.Hour.ValueInt64())
		args["compactionHourTimezone"] = gql.Optional(
			compaction.HourTimezone.ValueString(),
		)
		args["compactionInterval"] = gql.Optional(compaction.Interval.ValueInt64())
		args["compactionIntervalUnit"] = gql.Optional(
			compaction.IntervalUnit.ValueString(),
		)
	}

	for name, attr := range map[string]basetypes.ListValue{
		"appendOnlyKeys": data.AppendOnlyKeys,
		"fullAccessKeys": data.FullAccessKeys,
		"rsyncKeys":      data.RsyncKeys,
	} {
		if attr.IsNull() || attr.IsUnknown() {
			continue
		}

		var keys []string
		diagnostics = attr.ElementsAs(ctx, &keys, false)
		if diagnostics.HasError() {
			return diagnostics
		}
		args[name] = gql.Optional(keys)
	}

	return diagnostics
}

func (r *BorgRepoResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data BorgRepoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := gql.Arguments{
		"name":   gql.Required(data.Name.ValueString()),
		"region": gql.Required(data.Region.ValueString()),
	}
	if diagnostics := setArguments(ctx, args, data); diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	var payload BorgRepoAddPayload
	if err := r.client.Mutation("repoAdd", &payload, args); err != nil {
		resp.Diagnostics.AddError("Failed to create borg repo", err.Error())
		return
	}
	resp.Diagnostics.Append(data.update(ctx, payload.RepoAdded)...)

	tflog.Trace(ctx, "created repo", map[string]interface{}{
		"id":   data.Id,
		"name": data.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BorgRepoResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data BorgRepoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload BorgReposPayload
	args := gql.Arguments{"name": gql.Optional(data.Name.ValueString())}
	if err := r.client.Query("repoList", &payload, args); err != nil {
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

func (r *BorgRepoResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data BorgRepoModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := gql.Arguments{
		"id":   gql.Required(data.Id.ValueString()),
		"name": gql.Optional(data.Name.ValueString()),
	}
	if diagnostics := setArguments(ctx, args, data); diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	var payload BorgRepoEditPayload
	if err := r.client.Mutation("repoEdit", &payload, args); err != nil {
		resp.Diagnostics.AddError("Failed to update borg repo", err.Error())
		return
	}
	resp.Diagnostics.Append(data.update(ctx, payload.RepoEdited)...)

	tflog.Trace(ctx, "updated repo", map[string]interface{}{
		"id":   data.Id,
		"name": data.Name,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BorgRepoResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data BorgRepoModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := gql.Arguments{"id": gql.Required(data.Id.ValueString())}
	if err := r.client.Mutation("repoDelete", &BorgRepoDeletePayload{}, args); err != nil {
		resp.Diagnostics.AddError("Failed to delete borg repo", err.Error())
	}

	tflog.Trace(ctx, "deleted repo", map[string]interface{}{
		"id":   data.Id,
		"name": data.Name,
	})
}

func (r *BorgRepoResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
