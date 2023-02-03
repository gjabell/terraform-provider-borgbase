package provider

import (
	"context"
	"fmt"

	"github.com/gjabell/terraform-provider-borgbase/gql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &SshKeyDataSource{}

func NewSshKeyDataSource() datasource.DataSource {
	return &SshKeyDataSource{}
}

type SshKeyDataSource struct {
	client *gql.Client
}

func (d *SshKeyDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (d *SshKeyDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Public SSH key for accessing repositories",
		Attributes: map[string]schema.Attribute{
			"added_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date when the key was added to BorgBase",
			},
			"bits": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of bits in the key",
			},
			"hash_md5": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "MD5 hash of the SSH key",
			},
			"hash_sha256": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA256 hash of the SSH key",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal BorgBase key identifier",
			},
			"last_used_at": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Date when the key was last used to access " +
					"BorgBase",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User-defined key identifier",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Public SSH key",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of the SSH key",
			},
		},
	}
}

func (d *SshKeyDataSource) Configure(
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
			fmt.Sprintf("Expected *gql.Client, got: %T. "+
				"Please report this issue to the provider developers.",
				req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *SshKeyDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data SshKeyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var payload SshKeysPayload
	args := gql.Arguments{"name": gql.Optional(data.Name.ValueString())}
	if err := d.client.Query("sshList", &payload, args); err != nil {
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
