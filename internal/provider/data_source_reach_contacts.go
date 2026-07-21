// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceReachContacts{}
	_ datasource.DataSourceWithConfigure = &DataSourceReachContacts{}
)

func NewDataSourceReachContacts() datasource.DataSource {
	return &DataSourceReachContacts{}
}

// DataSourceReachContacts defines the data source implementation.
type DataSourceReachContacts struct {
	client *client.ClientWithResponses
}

// ReachContactsItemModel maps a single contact from the API response.
type ReachContactsItemModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Email              types.String `tfsdk:"email"`
	Name               types.String `tfsdk:"name"`
	Surname            types.String `tfsdk:"surname"`
	SubscriptionStatus types.String `tfsdk:"subscription_status"`
	SubscribedAt       types.String `tfsdk:"subscribed_at"`
}

// DataSourceReachContactsModel describes the data source data model.
type DataSourceReachContactsModel struct {
	Search   types.String             `tfsdk:"search"`
	Contacts []ReachContactsItemModel `tfsdk:"contacts"`
}

func (d *DataSourceReachContacts) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reach_contacts"
}

func (d *DataSourceReachContacts) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists email marketing contacts from the Hostinger Reach API.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter contacts client-side by matching the given string against each contact's email or name (case-insensitive substring match).",
			},
			"contacts": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of email marketing contacts.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier of the contact.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Email address of the contact.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "First name of the contact.",
						},
						"surname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Last name (surname) of the contact.",
						},
						"subscription_status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Subscription status of the contact. Possible values: `subscribed`, `unsubscribed`.",
						},
						"subscribed_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "RFC3339 timestamp of when the contact subscribed.",
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceReachContacts) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = c
}

func (d *DataSourceReachContacts) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceReachContactsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	search := ""
	if !data.Search.IsNull() && data.Search.ValueString() != "" {
		search = data.Search.ValueString()
		ctx = tflog.SetField(ctx, "search", search)
	}

	page := 1
	for {
		p := page
		params := client.ReachListContactsV1Params{
			Page: &p,
		}

		response, err := d.client.ReachListContactsV1WithResponse(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Reach Contacts",
				fmt.Sprintf("Got error: %s", err),
			)
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read Reach Contacts",
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError(
				"Unable to Read Reach Contacts",
				"Response body is nil",
			)
			return
		}

		contacts := response.JSON200.Data
		if contacts == nil || len(*contacts) == 0 {
			break
		}

		for _, item := range *contacts {
			if search != "" {
				email := ""
				if item.Email != nil {
					email = *item.Email
				}
				name := ""
				if item.Name != nil {
					name = *item.Name
				}
				if !strings.Contains(strings.ToLower(email), strings.ToLower(search)) &&
					!strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
					continue
				}
			}

			var m ReachContactsItemModel
			m.Uuid = types.StringPointerValue(item.Uuid)
			m.Email = types.StringPointerValue(item.Email)
			m.Name = types.StringPointerValue(item.Name)
			m.Surname = types.StringPointerValue(item.Surname)
			m.SubscriptionStatus = types.StringPointerValue((*string)(item.SubscriptionStatus))
			if item.SubscribedAt != nil {
				m.SubscribedAt = types.StringValue(item.SubscribedAt.Format("2006-01-02T15:04:05Z07:00"))
			} else {
				m.SubscribedAt = types.StringNull()
			}

			data.Contacts = append(data.Contacts, m)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
