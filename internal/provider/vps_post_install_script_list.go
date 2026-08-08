// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ list.ListResource              = &VPSPostInstallScriptList{}
	_ list.ListResourceWithConfigure = &VPSPostInstallScriptList{}
)

func NewVPSPostInstallScriptList() list.ListResource {
	return &VPSPostInstallScriptList{}
}

// VPSPostInstallScriptList defines the list resource implementation.
type VPSPostInstallScriptList struct {
	client *client.ClientWithResponses
}

type VPSPostInstallScriptListModel struct{}

func (l *VPSPostInstallScriptList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected List Resource Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	l.client = c
}

func (l *VPSPostInstallScriptList) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_post_install_script"
}

func (l *VPSPostInstallScriptList) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource provides a list of VPS post-install scripts.",
		Attributes:          map[string]schema.Attribute{},
	}
}

func (l *VPSPostInstallScriptList) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data VPSPostInstallScriptListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		params := &client.VPSGetPostInstallScriptsV1Params{}
		page := 1

		for {
			params.Page = &page
			ctx = tflog.SetField(ctx, "page", params.Page)

			response, err := l.client.VPSGetPostInstallScriptsV1WithResponse(ctx, params)
			if err != nil {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Post-install Scripts",
					fmt.Sprintf("Got error: %s", err),
				)
				push(result)
				return
			}
			if response.StatusCode() != http.StatusOK {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Post-install Scripts",
					fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
				)
				push(result)
				return
			}
			if response.JSON200 == nil {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Post-install Scripts",
					"Response body is nil",
				)
				push(result)
				return
			}
			if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
				break
			}

			for _, item := range *response.JSON200.Data {
				var d VPSPostInstallScriptResourceModel
				d.Merge(item)

				result := req.NewListResult(ctx)
				result.DisplayName = d.Name.ValueString()

				identity := VPSPostInstallScriptIdentity{ID: d.ID}
				result.Diagnostics.Append(result.Identity.Set(ctx, &identity)...)

				if req.IncludeResource {
					d.Timeouts = resourcetimeouts.Value{
						Object: types.ObjectNull(map[string]attr.Type{
							"create": types.StringType,
							"read":   types.StringType,
							"update": types.StringType,
							"delete": types.StringType,
						}),
					}

					result.Diagnostics.Append(result.Resource.Set(ctx, &d)...)
				}

				if !push(result) {
					return
				}
			}

			meta := response.JSON200.Meta
			if meta == nil || meta.CurrentPage == nil || meta.PerPage == nil || meta.Total == nil {
				break
			}
			fetched := (*meta.CurrentPage) * (*meta.PerPage)
			if fetched >= *meta.Total {
				break
			}
			page++
		}
	}
}
