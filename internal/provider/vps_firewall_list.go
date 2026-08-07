package provider

import (
	"context"
	"fmt"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/list/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ list.ListResource              = &VPSFirewallList{}
	_ list.ListResourceWithConfigure = &VPSFirewallList{}
)

func NewVPSFirewallList() list.ListResource {
	return &VPSFirewallList{}
}

// VPSFirewallList defines the list resource implementation.
type VPSFirewallList struct {
	client *client.ClientWithResponses
}

type VPSFirewallListModel struct {
	Timeouts          timeouts.Type `tfsdk:"timeouts"`
	ResourceGroupName types.String  `tfsdk:"resource_group_name"`
}

func (l *VPSFirewallList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	l.client = c
}

func (l *VPSFirewallList) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall"
}

func (l *VPSFirewallList) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource provides a list of VPS firewalls.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx),
			"resource_group_name": schema.StringAttribute{
				Description: "Name of the resource group to list things in.",
				Required:    true,
			},
		},
	}
}

func (l *VPSFirewallList) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data VPSFirewallListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		result := req.NewListResult(ctx)
		result.DisplayName = "test"
		result.Diagnostics.Append(result.Identity.Set(ctx, "test")...)

		if !push(result) {
			return
		}
	}
}
