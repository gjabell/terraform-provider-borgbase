package provider

import (
	"context"
	"fmt"

	"github.com/gjabell/terraform-provider-borgbase/gql"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SshKeyResource{}
var _ resource.ResourceWithImportState = &SshKeyResource{}

func NewSshKeyResource() resource.Resource {
	return &SshKeyResource{}
}

type SshKeyResource struct {
	client *gql.Client
}

func (r *SshKeyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SshKeyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Public SSH key for accessing repositories",
		Attributes: map[string]schema.Attribute{
			"added_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the key was added to BorgBase",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bits": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of bits in the key",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"hash_md5": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "MD5 hash of the SSH key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hash_sha256": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA256 hash of the SSH key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal BorgBase key identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_used_at": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Date when the key was last used to access " +
					"BorgBase",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined key identifier",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public SSH key",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of the SSH key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SshKeyResource) Configure(
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
			fmt.Sprintf("Expected *gql.Client, got: %T. "+
				"Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *SshKeyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data SshKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload SshAddPayload
	args := gql.Arguments{
		"name":    gql.Optional(data.Name.ValueString()),
		"keyData": gql.Optional(data.PublicKey.ValueString()),
	}
	if err := r.client.Mutation("sshAdd", &payload, args); err != nil {
		resp.Diagnostics.AddError("Failed to create SSH key", err.Error())
		return
	}
	data.update(payload.KeyAdded)

	tflog.Trace(ctx, "created SSH key", map[string]interface{}{
		"id":         data.Id,
		"name":       data.Name,
		"public_key": data.PublicKey,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshKeyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data SshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload SshKeysPayload
	args := gql.Arguments{"name": gql.Optional(data.Name.ValueString())}
	if err := r.client.Query("sshList", &payload, args); err != nil {
		resp.Diagnostics.AddError("Failed to read SSH key", err.Error())
		return
	}

	var key *SshKeyPayload
	for _, item := range payload {
		if item.Name == data.Name.ValueString() {
			key = &item
			break
		}
	}
	if key == nil {
		resp.Diagnostics.AddError("Unknown SSH key", data.Name.String())
		return
	}
	data.update(*key)

	tflog.Trace(ctx, "read SSH key", map[string]interface{}{
		"id":         data.Id,
		"name":       data.Name,
		"public_key": data.PublicKey,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshKeyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Update is not supported for SSH keys.
}

func (r *SshKeyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data SshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := gql.Arguments{"id": gql.Required(data.Id.ValueString())}
	if err := r.client.Mutation("sshDelete", &SshDeletePayload{}, args); err != nil {
		resp.Diagnostics.AddError("Failed to delete SSH key", err.Error())
	}

	tflog.Trace(ctx, "deleted SSH key", map[string]interface{}{
		"id":         data.Id,
		"name":       data.Name,
		"public_key": data.PublicKey,
	})
}

func (r *SshKeyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
