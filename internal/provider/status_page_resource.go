// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	kumastatuspage "github.com/breml/go-uptime-kuma-client/statuspage"
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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
				Computed:            true,
				Default:             stringdefault.StaticString("/icon.svg"),
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
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
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

	slug := data.Slug.ValueString()
	title := data.Title.ValueString()

	tflog.Info(ctx, "Creating status page", map[string]interface{}{
		"slug":  slug,
		"title": title,
	})

	// 1. Create Status Page (only takes slug and title)
	if err := r.client.Kuma.AddStatusPage(ctx, title, slug); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status page: %s", err))
		return
	}

	// 2. Prepare full status page object for update
	sp := &kumastatuspage.StatusPage{
		Slug:              slug,
		Title:             title,
		Description:       data.Description.ValueString(),
		Theme:             data.Theme.ValueString(),
		Published:         data.Published.ValueBool(),
		ShowTags:          data.ShowTags.ValueBool(),
		FooterText:        data.FooterText.ValueString(),
		CustomCSS:         data.CustomCSS.ValueString(),
		GoogleAnalyticsID: data.GoogleAnalyticsID.ValueString(),
		Icon:              data.Icon.ValueString(),
		ShowPoweredBy:     data.ShowPoweredBy.ValueBool(),
	}

	// Domain Names
	sp.DomainNameList = []string{}
	if len(data.DomainNameList) > 0 {
		sp.DomainNameList = make([]string, len(data.DomainNameList))
		for i, v := range data.DomainNameList {
			sp.DomainNameList[i] = v.ValueString()
		}
	}

	// Public Groups
	sp.PublicGroupList = []kumastatuspage.PublicGroup{}
	if len(data.PublicGroupList) > 0 {
		sp.PublicGroupList = make([]kumastatuspage.PublicGroup, len(data.PublicGroupList))
		for i, g := range data.PublicGroupList {
			pg := kumastatuspage.PublicGroup{
				Name:        g.Name.ValueString(),
				Weight:      int(g.Weight.ValueInt64()),
				MonitorList: []kumastatuspage.PublicMonitor{},
			}

			if len(g.MonitorList) > 0 {
				pg.MonitorList = make([]kumastatuspage.PublicMonitor, len(g.MonitorList))
				for j, mid := range g.MonitorList {
					pg.MonitorList[j] = kumastatuspage.PublicMonitor{
						ID: mid.ValueInt64(),
					}
				}
			}
			sp.PublicGroupList[i] = pg
		}
	}

	// 3. Update (Save) the status page
	publicGroups, err := r.client.Kuma.SaveStatusPage(ctx, sp)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to save status page details: %s", err))
		// Should we rollback?
		return
	}

	// 4. Read back to get Status Page ID?
	// SaveStatusPage returns PublicGroups but not the page ID.
	fetchedSP, err := r.client.Kuma.GetStatusPage(ctx, slug)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read created status page: %s", err))
		return
	}

	// 5. Update state
	data.ID = types.Int64Value(fetchedSP.ID)

	// Map back groups to get their IDs
	// Use returned publicGroups if available (it has IDs) or fetchedSP
	if len(fetchedSP.PublicGroupList) > 0 && len(data.PublicGroupList) > 0 {
		for i, apiGroup := range fetchedSP.PublicGroupList {
			if i < len(data.PublicGroupList) {
				data.PublicGroupList[i].ID = types.Int64Value(apiGroup.ID)
			}
		}
	} else if len(publicGroups) > 0 && len(data.PublicGroupList) > 0 {
		// Use publicGroups returned from SaveStatusPage if fetchedSP fails or as backup
		for i, apiGroup := range publicGroups {
			if i < len(data.PublicGroupList) {
				data.PublicGroupList[i].ID = types.Int64Value(apiGroup.ID)
			}
		}
	}

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

	slug := data.Slug.ValueString()

	// Read status page from API
	sp, err := r.client.Kuma.GetStatusPage(ctx, slug)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status page '%s': %s", slug, err))
		return
	}

	if sp == nil { // Probably handled by err, but just in case
		resp.State.RemoveResource(ctx)
		return
	}

	// Map API response to model
	data.ID = types.Int64Value(sp.ID)
	data.Title = types.StringValue(sp.Title)

	if sp.Description != "" {
		data.Description = types.StringValue(sp.Description)
	} else {
		data.Description = types.StringNull()
	}

	if sp.Theme != "" {
		data.Theme = types.StringValue(sp.Theme)
	} else {
		data.Theme = types.StringNull()
	}

	data.Published = types.BoolValue(sp.Published)
	data.ShowTags = types.BoolValue(sp.ShowTags)

	if sp.FooterText != "" {
		data.FooterText = types.StringValue(sp.FooterText)
	} else {
		data.FooterText = types.StringNull()
	}

	if sp.CustomCSS != "" {
		data.CustomCSS = types.StringValue(sp.CustomCSS)
	} else {
		data.CustomCSS = types.StringNull()
	}

	if sp.GoogleAnalyticsID != "" {
		data.GoogleAnalyticsID = types.StringValue(sp.GoogleAnalyticsID)
	} else {
		data.GoogleAnalyticsID = types.StringNull()
	}

	if sp.Icon != "" {
		data.Icon = types.StringValue(sp.Icon)
	} else {
		data.Icon = types.StringNull()
	}

	data.ShowPoweredBy = types.BoolValue(sp.ShowPoweredBy)

	// Domain Names
	if len(sp.DomainNameList) > 0 {
		outDomains := make([]types.String, len(sp.DomainNameList))
		for i, v := range sp.DomainNameList {
			outDomains[i] = types.StringValue(v)
		}
		data.DomainNameList = outDomains
	}

	// Public Groups
	// GetStatusPage says "PublicGroupList must be maintained separately" in comment (Step 259)
	// But `GetStatusPage` return struct has `PublicGroupList`.
	// Although checking library code (Step 259, Line 28 comment): "Note: The server does not return PublicGroupList in this endpoint."
	// Wait, if it doesn't return PublicGroupList, we lose that state on Read!
	// This is a known issue in Uptime Kuma v1 API?
	// But `client.state.statusPages` cache might have it?
	// `GetStatusPage` calls `syncEmit("getStatusPage")`.

	// If the API doesn't return PublicGroups on Get, how do we Read them?
	// Maybe `GetStatusPages` (plural) returns everything?
	// `GetMonitor` returns monitor list via state.
	// `GetStatusPages` uses `c.state.statusPages`.

	// Let's check `GetStatusPages` again.
	// Line 11: returns map.
	// The state is updated via socket events.
	// If we use `GetStatusPages`, we rely on cache.
	// But `GetStatusPage(slug)` calls API directly.

	// If `GetStatusPage(slug)` returns incomplete data, we have a problem.
	// However, `go-uptime-kuma-client` `GetStatusPage` implementation calls `emit("getStatusPage")`.
	// Does `getStatusPage` event return groups?
	// The comment says no.

	// If so, we might need to rely on `GetStatusPages` (plural) from state if available?
	// But `GetStatusPages` needs state populate.
	// The client connects and performs full sync. So `c.state.statusPages` should be populated.
	// So maybe we should iterate `GetStatusPages` to find our slug?

	// Try to find matching page in cache which might have more details
	// If `getStatusPage` API is limited.
	allPages, err := r.client.Kuma.GetStatusPages(ctx)
	if err == nil {
		for _, page := range allPages {
			if page.Slug == slug {
				sp.PublicGroupList = page.PublicGroupList
				break
			}
		}
	}

	if len(sp.PublicGroupList) > 0 {
		outGroups := make([]PublicGroupModel, len(sp.PublicGroupList))
		for i, g := range sp.PublicGroupList {
			pgModel := PublicGroupModel{
				ID:     types.Int64Value(g.ID),
				Name:   types.StringValue(g.Name),
				Weight: types.Int64Value(int64(g.Weight)),
			}

			if len(g.MonitorList) > 0 {
				mList := make([]types.Int64, len(g.MonitorList))
				for j, m := range g.MonitorList {
					mList[j] = types.Int64Value(m.ID)
				}
				pgModel.MonitorList = mList
			}
			outGroups[i] = pgModel
		}
		data.PublicGroupList = outGroups
	} else {
		// If no groups found in API/Cache, preserve existing state.
		// We cannot distinguish between "groups deleted" and "API didn't return groups".
		// We assume Terraform manages the state.
		tflog.Warn(ctx, fmt.Sprintf("No public groups found for status page '%s' in API/Cache; preserving state", slug))
	}

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

	slug := data.Slug.ValueString()

	// Prepare update
	sp := &kumastatuspage.StatusPage{
		Slug:              slug,
		Title:             data.Title.ValueString(),
		Description:       data.Description.ValueString(),
		Theme:             data.Theme.ValueString(),
		Published:         data.Published.ValueBool(),
		ShowTags:          data.ShowTags.ValueBool(),
		FooterText:        data.FooterText.ValueString(),
		CustomCSS:         data.CustomCSS.ValueString(),
		GoogleAnalyticsID: data.GoogleAnalyticsID.ValueString(),
		Icon:              data.Icon.ValueString(),
		ShowPoweredBy:     data.ShowPoweredBy.ValueBool(),
	}

	// Domain Names
	sp.DomainNameList = []string{}
	if len(data.DomainNameList) > 0 {
		sp.DomainNameList = make([]string, len(data.DomainNameList))
		for i, v := range data.DomainNameList {
			sp.DomainNameList[i] = v.ValueString()
		}
	}

	// Public Groups
	sp.PublicGroupList = []kumastatuspage.PublicGroup{}
	if len(data.PublicGroupList) > 0 {
		sp.PublicGroupList = make([]kumastatuspage.PublicGroup, len(data.PublicGroupList))
		for i, g := range data.PublicGroupList {
			pg := kumastatuspage.PublicGroup{
				Name:        g.Name.ValueString(),
				Weight:      int(g.Weight.ValueInt64()),
				MonitorList: []kumastatuspage.PublicMonitor{},
			}
			if !g.ID.IsNull() {
				pg.ID = g.ID.ValueInt64()
			}

			if len(g.MonitorList) > 0 {
				pg.MonitorList = make([]kumastatuspage.PublicMonitor, len(g.MonitorList))
				for j, mid := range g.MonitorList {
					pg.MonitorList[j] = kumastatuspage.PublicMonitor{
						ID: mid.ValueInt64(),
					}
				}
			}
			sp.PublicGroupList[i] = pg
		}
	}

	// Update (Save) the status page
	publicGroups, err := r.client.Kuma.SaveStatusPage(ctx, sp)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status page: %s", err))
		return
	}

	// Update IDs from response
	if len(publicGroups) > 0 && len(data.PublicGroupList) > 0 {
		for i, apiGroup := range publicGroups {
			if i < len(data.PublicGroupList) {
				data.PublicGroupList[i].ID = types.Int64Value(apiGroup.ID)
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

	slug := data.Slug.ValueString()

	// Delete the status page
	if err := r.client.Kuma.DeleteStatusPage(ctx, slug); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status page '%s': %s", slug, err))
		return
	}
}

func (r *StatusPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Slug is the primary identifier for status pages
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}
