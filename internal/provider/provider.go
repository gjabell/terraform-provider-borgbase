package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/gjabell/terraform-provider-borgbase/gql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const borgBaseApi = "https://api.borgbase.com/graphql"

const apiTokenEnvVar = "BORGBASE_API_TOKEN"

// Ensure BorgBaseProvider satisfies various provider interfaces.
var _ provider.Provider = &BorgBaseProvider{}

type BorgBaseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type BorgBaseProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *BorgBaseProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "borgbase"
	resp.Version = p.version
}

func (p *BorgBaseProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "BorgBase API token",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *BorgBaseProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data BorgBaseProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv(apiTokenEnvVar)
	if data.ApiToken.ValueString() != "" {
		apiToken = data.ApiToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddError("Missing API token",
			fmt.Sprintf("A BorgBase API token must be provided in the provider "+
				"configuration block in the api_token attribute or in the %s env var.",
				apiTokenEnvVar))
	}

	client := gql.NewClient(borgBaseApi, apiToken)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *BorgBaseProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	return []func() resource.Resource{
		NewBorgRepoResource,
		NewSshKeyResource,
	}
}

func (p *BorgBaseProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBorgRepoDataSource,
		NewSshKeyDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BorgBaseProvider{
			version: version,
		}
	}
}
