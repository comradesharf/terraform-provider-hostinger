// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Uuid               types.String      `tfsdk:"uuid"`
	Email              types.String      `tfsdk:"email"`
	Name               types.String      `tfsdk:"name"`
	Surname            types.String      `tfsdk:"surname"`
	SubscriptionStatus types.String      `tfsdk:"subscription_status"`
	SubscribedAt       timetypes.RFC3339 `tfsdk:"subscribed_at"`
	Source             types.String      `tfsdk:"source"`
	Note               types.String      `tfsdk:"note"`
}

// DataSourceReachContactsModel describes the data source data model.
type DataSourceReachContactsModel struct {
	GroupUuid          types.String             `tfsdk:"group_uuid"`
	SubscriptionStatus types.String             `tfsdk:"subscription_status"`
	Contacts           []ReachContactsItemModel `tfsdk:"contacts"`
}

func (d *DataSourceReachContacts) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reach_contacts"
}

func (d *DataSourceReachContacts) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists email marketing contacts from the Hostinger Reach API.",
		Attributes: map[string]schema.Attribute{
			"group_uuid": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter contacts by group UUID",
			},
			"subscription_status": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter contacts by subscription status",
				Validators: []validator.String{
					stringvalidator.OneOf("subscribed", "unsubscribed"),
				},
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
							CustomType:          timetypes.RFC3339Type{},
						},
						"source": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Source of the contact.",
						},
						"note": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Note of the contact.",
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

	if data.GroupUuid.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Group UUID",
			"The 'group_uuid' attribute cannot be unknown.",
		)
	}

	if data.SubscriptionStatus.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Subscription Status",
			"The 'subscription_status' attribute cannot be unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	params := client.ReachListContactsV1Params{}

	if !data.GroupUuid.IsNull() && data.GroupUuid.ValueString() != "" {
		params.GroupUuid = data.GroupUuid.ValueStringPointer()
		ctx = tflog.SetField(ctx, "group_uuid", &params.GroupUuid)
	}

	if !data.SubscriptionStatus.IsNull() && data.SubscriptionStatus.ValueString() != "" {
		params.SubscriptionStatus = (*client.ReachListContactsV1ParamsSubscriptionStatus)(data.SubscriptionStatus.ValueStringPointer())
		ctx = tflog.SetField(ctx, "subscription_status", &params.SubscriptionStatus)
	}

	page := 0
	for {
		params.Page = &page

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
			var m ReachContactsItemModel
			m.Uuid = types.StringPointerValue(item.Uuid)
			m.Email = types.StringPointerValue(item.Email)
			m.Name = types.StringPointerValue(item.Name)
			m.Surname = types.StringPointerValue(item.Surname)
			m.SubscriptionStatus = types.StringPointerValue((*string)(item.SubscriptionStatus))
			m.Source = types.StringPointerValue((*string)(item.Source))
			m.Note = types.StringPointerValue(item.Note)
			m.SubscribedAt = timetypes.NewRFC3339TimePointerValue(item.SubscribedAt)

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
