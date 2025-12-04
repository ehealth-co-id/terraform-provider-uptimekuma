// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ehealth-co-id/terraform-provider-uptimekuma/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StatusPageResource{}
var _ resource.ResourceWithImportState = &StatusPageResource{}

func NewStatusPageResource() resource.Resource {
	return &StatusPageResource{}
}

// StatusPageResource defines the resource implementation.
type StatusPageResource struct {
	client *client.Client
}

// PublicGroupModel describes a group of monitors on a status page.
type PublicGroupModel struct {
	ID          types.Int64   `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Weight      types.Int64   `tfsdk:"weight"`
	MonitorList []types.Int64 `tfsdk:"monitor_list"`
}

// StatusPageResourceModel describes the resource data model.
type StatusPageResourceModel struct {
	ID                types.Int64        `tfsdk:"id"`
	Slug              types.String       `tfsdk:"slug"`
	Title             types.String       `tfsdk:"title"`
	Description       types.String       `tfsdk:"description"`
	Theme             types.String       `tfsdk:"theme"`
	Published         types.Bool         `tfsdk:"published"`
	ShowTags          types.Bool         `tfsdk:"show_tags"`
	DomainNameList    []types.String     `tfsdk:"domain_name_list"`
	FooterText        types.String       `tfsdk:"footer_text"`
	CustomCSS         types.String       `tfsdk:"custom_css"`
	GoogleAnalyticsID types.String       `tfsdk:"google_analytics_id"`
	Icon              types.String       `tfsdk:"icon"`
	ShowPoweredBy     types.Bool         `tfsdk:"show_powered_by"`
	PublicGroupList   []PublicGroupModel `tfsdk:"public_group_list"`
}

func (r *StatusPageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uptime Kuma Status Page resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Status page identifier",
				PlanModifiers:       []planmodifier.Int64{
					// UseStateForUnknown tells Terraform to keep the value from the prior state if it's not explicitly set in the configuration.
					// This is useful for computed attributes that don't change.
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Status page URL slug",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Status page title",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Status page description",
				Optional:            true,
			},
			"theme": schema.StringAttribute{
				MarkdownDescription: "Status page theme",
				Optional:            true,
			},
			"published": schema.BoolAttribute{
				MarkdownDescription: "Whether the status page is published",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"show_tags": schema.BoolAttribute{
				MarkdownDescription: "Whether to show tags on the status page",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"domain_name_list": schema.ListAttribute{
				MarkdownDescription: "List of custom domain names for the status page",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"footer_text": schema.StringAttribute{
				MarkdownDescription: "Custom footer text",
				Optional:            true,
			},
			"custom_css": schema.StringAttribute{
				MarkdownDescription: "Custom CSS for the status page",
				Optional:            true,
			},
			"google_analytics_id": schema.StringAttribute{
				MarkdownDescription: "Google Analytics ID",
				Optional:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Status page icon",
				Optional:            true,
			},
			"show_powered_by": schema.BoolAttribute{
				MarkdownDescription: "Whether to show 'Powered by Uptime Kuma' text",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"public_group_list": schema.ListNestedAttribute{
				MarkdownDescription: "List of monitor groups displayed on the status page",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Group identifier",
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Group name",
							Required:            true,
						},
						"weight": schema.Int64Attribute{
							MarkdownDescription: "Group order weight",
							Optional:            true,
						},
						"monitor_list": schema.ListAttribute{
							MarkdownDescription: "List of monitor IDs in the group",
							Optional:            true,
							ElementType:         types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (r *StatusPageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *StatusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StatusPageResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the status page first with just slug and title
	createRequest := &client.AddStatusPageRequest{
		Slug:  data.Slug.ValueString(),
		Title: data.Title.ValueString(),
	}

	tflog.Info(ctx, "Creating status page", map[string]interface{}{
		"slug":  createRequest.Slug,
		"title": createRequest.Title,
	})

	createResp, err := r.client.CreateStatusPage(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status page: %s", err))
		return
	}

	tflog.Info(ctx, "Status page created", map[string]interface{}{
		"message": createResp.Msg,
	})

	// Now get the status page to find its ID and other details
	createdPage, err := r.client.GetStatusPage(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve created status page: %s", err))
		return
	}

	// Now update the status page with all other attributes
	updateRequest := &client.SaveStatusPageRequest{
		Title:     data.Title.ValueString(),
		Published: data.Published.ValueBool(),
		ShowTags:  data.ShowTags.ValueBool(),
	}

	// Set optional fields
	if !data.Description.IsNull() {
		updateRequest.Description = data.Description.ValueString()
	}

	if !data.Theme.IsNull() {
		updateRequest.Theme = data.Theme.ValueString()
	}

	// Convert domain name list
	if len(data.DomainNameList) > 0 {
		domainNames := make([]string, 0, len(data.DomainNameList))
		for _, domain := range data.DomainNameList {
			domainNames = append(domainNames, domain.ValueString())
		}
		updateRequest.DomainNameList = domainNames
	}

	if !data.FooterText.IsNull() {
		updateRequest.FooterText = data.FooterText.ValueString()
	}

	if !data.CustomCSS.IsNull() {
		updateRequest.CustomCSS = data.CustomCSS.ValueString()
	}

	if !data.GoogleAnalyticsID.IsNull() {
		updateRequest.GoogleAnalyticsID = data.GoogleAnalyticsID.ValueString()
	}

	if !data.Icon.IsNull() {
		updateRequest.Icon = data.Icon.ValueString()
	}

	if !data.ShowPoweredBy.IsNull() {
		updateRequest.ShowPoweredBy = data.ShowPoweredBy.ValueBool()
	}

	// Add public groups
	if len(data.PublicGroupList) > 0 {
		groups := make([]client.PublicGroup, 0, len(data.PublicGroupList))
		for _, group := range data.PublicGroupList {
			newGroup := client.PublicGroup{
				Name:   group.Name.ValueString(),
				Weight: int(group.Weight.ValueInt64()),
			}

			// Convert monitor list
			monitors := make([]int, 0, len(group.MonitorList))
			for _, monitorID := range group.MonitorList {
				monitors = append(monitors, int(monitorID.ValueInt64()))
			}
			newGroup.MonitorList = monitors

			groups = append(groups, newGroup)
		}
		updateRequest.PublicGroupList = groups
	}

	// Update the status page with all attributes
	_, err = r.client.UpdateStatusPage(ctx, data.Slug.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status page attributes: %s", err))
		return
	}

	// Update local state
	data.ID = types.Int64Value(int64(createdPage.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StatusPageResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read status page from API
	statusPage, err := r.client.GetStatusPage(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status page '%s': %s", data.Slug.ValueString(), err))
		return
	}

	// Update model with API data
	data.ID = types.Int64Value(int64(statusPage.ID))
	data.Title = types.StringValue(statusPage.Title)
	data.Description = types.StringValue(statusPage.Description)
	data.Theme = types.StringValue(statusPage.Theme)
	data.Published = types.BoolValue(statusPage.Published)
	data.ShowTags = types.BoolValue(statusPage.ShowTags)

	// Convert domain names
	domainNames := make([]types.String, 0, len(statusPage.DomainNameList))
	for _, domain := range statusPage.DomainNameList {
		domainNames = append(domainNames, types.StringValue(domain))
	}
	data.DomainNameList = domainNames

	data.FooterText = types.StringValue(statusPage.FooterText)
	data.CustomCSS = types.StringValue(statusPage.CustomCSS)
	data.GoogleAnalyticsID = types.StringValue(statusPage.GoogleAnalyticsID)
	data.Icon = types.StringValue(statusPage.Icon)
	data.ShowPoweredBy = types.BoolValue(statusPage.ShowPoweredBy)

	// Convert public groups
	groups := make([]PublicGroupModel, 0, len(statusPage.PublicGroupList))
	for _, apiGroup := range statusPage.PublicGroupList {
		group := PublicGroupModel{
			ID:     types.Int64Value(int64(apiGroup.ID)),
			Name:   types.StringValue(apiGroup.Name),
			Weight: types.Int64Value(int64(apiGroup.Weight)),
		}

		// Convert monitor list
		monitors := make([]types.Int64, 0, len(apiGroup.MonitorList))
		for _, monitorID := range apiGroup.MonitorList {
			monitors = append(monitors, types.Int64Value(int64(monitorID)))
		}
		group.MonitorList = monitors

		groups = append(groups, group)
	}
	data.PublicGroupList = groups

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StatusPageResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare update request
	updateRequest := &client.SaveStatusPageRequest{
		Title:     data.Title.ValueString(),
		Published: data.Published.ValueBool(),
		ShowTags:  data.ShowTags.ValueBool(),
	}

	// Set optional fields
	if !data.Description.IsNull() {
		updateRequest.Description = data.Description.ValueString()
	}

	if !data.Theme.IsNull() {
		updateRequest.Theme = data.Theme.ValueString()
	}

	// Convert domain name list
	if len(data.DomainNameList) > 0 {
		domainNames := make([]string, 0, len(data.DomainNameList))
		for _, domain := range data.DomainNameList {
			domainNames = append(domainNames, domain.ValueString())
		}
		updateRequest.DomainNameList = domainNames
	}

	if !data.FooterText.IsNull() {
		updateRequest.FooterText = data.FooterText.ValueString()
	}

	if !data.CustomCSS.IsNull() {
		updateRequest.CustomCSS = data.CustomCSS.ValueString()
	}

	if !data.GoogleAnalyticsID.IsNull() {
		updateRequest.GoogleAnalyticsID = data.GoogleAnalyticsID.ValueString()
	}

	if !data.Icon.IsNull() {
		updateRequest.Icon = data.Icon.ValueString()
	}

	if !data.ShowPoweredBy.IsNull() {
		updateRequest.ShowPoweredBy = data.ShowPoweredBy.ValueBool()
	}

	// Add public groups
	if len(data.PublicGroupList) > 0 {
		groups := make([]client.PublicGroup, 0, len(data.PublicGroupList))
		for _, group := range data.PublicGroupList {
			newGroup := client.PublicGroup{
				Name:   group.Name.ValueString(),
				Weight: int(group.Weight.ValueInt64()),
			}

			// Convert monitor list
			monitors := make([]int, 0, len(group.MonitorList))
			for _, monitorID := range group.MonitorList {
				monitors = append(monitors, int(monitorID.ValueInt64()))
			}
			newGroup.MonitorList = monitors

			groups = append(groups, newGroup)
		}
		updateRequest.PublicGroupList = groups
	}

	// Update the status page
	tflog.Info(ctx, "Updating status page", map[string]interface{}{
		"slug": data.Slug.ValueString(),
	})

	_, err := r.client.UpdateStatusPage(ctx, data.Slug.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status page: %s", err))
		return
	}

	// Refresh the data from the API
	updatedPage, err := r.client.GetStatusPage(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated status page: %s", err))
		return
	}

	// Update the resource's groups with their IDs from the API
	if len(updatedPage.PublicGroupList) > 0 && len(data.PublicGroupList) > 0 {
		for i, apiGroup := range updatedPage.PublicGroupList {
			if i < len(data.PublicGroupList) {
				data.PublicGroupList[i].ID = types.Int64Value(int64(apiGroup.ID))
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StatusPageResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the status page
	tflog.Info(ctx, "Deleting status page", map[string]interface{}{
		"slug": data.Slug.ValueString(),
	})

	_, err := r.client.DeleteStatusPage(ctx, data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status page '%s': %s", data.Slug.ValueString(), err))
		return
	}
}

func (r *StatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Slug is the primary identifier for status pages
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}
