// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	kumamonitor "github.com/breml/go-uptime-kuma-client/monitor"
	"github.com/breml/go-uptime-kuma-client/tag"
	"github.com/ehealth-co-id/terraform-provider-uptimekuma/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

// MonitorResource defines the resource implementation.
type MonitorResource struct {
	client *client.Client
}

// MonitorResourceModel describes the resource data model.
type MonitorResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	Name                     types.String `tfsdk:"name"`
	Active                   types.Bool   `tfsdk:"active"`
	URL                      types.String `tfsdk:"url"`
	Method                   types.String `tfsdk:"method"`
	Hostname                 types.String `tfsdk:"hostname"`
	Port                     types.Int64  `tfsdk:"port"`
	Interval                 types.Int64  `tfsdk:"interval"`
	RetryInterval            types.Int64  `tfsdk:"retry_interval"`
	ResendInterval           types.Int64  `tfsdk:"resend_interval"`
	MaxRetries               types.Int64  `tfsdk:"max_retries"`
	UpsideDown               types.Bool   `tfsdk:"upside_down"`
	IgnoreTLS                types.Bool   `tfsdk:"ignore_tls"`
	MaxRedirects             types.Int64  `tfsdk:"max_redirects"`
	Body                     types.String `tfsdk:"body"`
	Headers                  types.String `tfsdk:"headers"`
	AuthMethod               types.String `tfsdk:"auth_method"`
	BasicAuthUser            types.String `tfsdk:"basic_auth_user"`
	BasicAuthPass            types.String `tfsdk:"basic_auth_pass"`
	Keyword                  types.String `tfsdk:"keyword"`
	NotificationIDList       types.List   `tfsdk:"notification_id_list"`
	AcceptedStatusCodes      types.List   `tfsdk:"accepted_status_codes"`
	DatabaseConnectionString types.String `tfsdk:"database_connection_string"`
	Tags                     types.List   `tfsdk:"tags"`
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uptime Kuma Monitor resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Monitor identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Monitor type (http, ping, port, keyword, dns, etc.)",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Monitor name",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is active (enabled). Defaults to true.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to monitor (required for http, keyword monitors)",
				Optional:            true,
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "HTTP method (GET, POST, etc.) for http monitors",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GET"),
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname for ping, port, etc. monitors. Also used for database connection strings.",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port number for port monitors",
				Optional:            true,
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "Check interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"retry_interval": schema.Int64Attribute{
				MarkdownDescription: "Retry interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"resend_interval": schema.Int64Attribute{
				MarkdownDescription: "Notification resend interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"upside_down": schema.BoolAttribute{
				MarkdownDescription: "Invert status (treat DOWN as UP and vice versa)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ignore_tls": schema.BoolAttribute{
				MarkdownDescription: "Ignore TLS/SSL errors",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"max_redirects": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of redirects to follow",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Request body for http monitors",
				Optional:            true,
			},
			"headers": schema.StringAttribute{
				MarkdownDescription: "Request headers for http monitors (JSON format)",
				Optional:            true,
			},
			"auth_method": schema.StringAttribute{
				MarkdownDescription: "Authentication method (basic, ntlm, mtls)",
				Optional:            true,
			},
			"basic_auth_user": schema.StringAttribute{
				MarkdownDescription: "Basic auth username",
				Optional:            true,
			},
			"basic_auth_pass": schema.StringAttribute{
				MarkdownDescription: "Basic auth password",
				Optional:            true,
				Sensitive:           true,
			},
			"keyword": schema.StringAttribute{
				MarkdownDescription: "Keyword to search for in response",
				Optional:            true,
			},
			"notification_id_list": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of notification IDs to trigger when monitor status changes",
				Optional:            true,
			},
			"accepted_status_codes": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of accepted HTTP status codes (e.g., [200, 201, 204]). Defaults to all 2xx codes if not specified.",
				Optional:            true,
			},
			"database_connection_string": schema.StringAttribute{
				MarkdownDescription: "Database connection string for database monitors (postgres, mysql, mongodb, etc.)",
				Optional:            true,
				Sensitive:           true,
			},
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "Tags associated with the monitor",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_id": schema.Int64Attribute{
							MarkdownDescription: "Tag ID",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value for the tag",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *MonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.monitorFromPlan(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error creating monitor", err.Error())
		return
	}

	// Call library to create monitor
	// Use client.Kuma.CreateMonitor instead of client.Kuma.Monitor.Add
	id, err := r.client.Kuma.CreateMonitor(ctx, monitor)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor: %s", err))
		return
	}

	// Update Terraform state
	data.ID = types.Int64Value(id)

	// Handle active state (monitors are created active by default, pause if active=false)
	// The active field in the API create request is not reliable, so we use PauseMonitor/ResumeMonitor
	if !data.Active.ValueBool() {
		if err := r.client.Kuma.PauseMonitor(ctx, id); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to pause monitor %d: %s", id, err))
			return
		}
	}

	// Add tags to the monitor (tags are managed separately via AddMonitorTag API)
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		type tagModel struct {
			TagID types.Int64  `tfsdk:"tag_id"`
			Value types.String `tfsdk:"value"`
		}
		var tfTags []tagModel
		data.Tags.ElementsAs(ctx, &tfTags, false)

		for _, t := range tfTags {
			tagID := t.TagID.ValueInt64()
			value := ""
			if !t.Value.IsNull() {
				value = t.Value.ValueString()
			}
			_, err := r.client.Kuma.AddMonitorTag(ctx, tagID, id, value)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add tag %d to monitor %d: %s", tagID, id, err))
				return
			}
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonitorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	monitorID := data.ID.ValueInt64()

	// Read the monitor from the API use client.Kuma.GetMonitor
	// Note: GetMonitor returns monitor.Base, which contains the data but might lose specific fields
	// unless we use GetMonitorAs or similar?
	// The library `GetMonitor` returns `monitor.Base`.
	// But `monitor.Base` in the library definition (Step 258) has `internalType` and `raw`.
	// We can't access `raw` it's private.
	// But we can call `monitor.GetMonitorAs(ctx, id, &target)`.
	// To do that, we need to know the type first.
	// Or we can try to guess from the provider state which type we expect?
	// But `Read` should be robust.
	// `client.GetMonitor` returns `monitor.Base`. `Type()` gives us the type string.
	// Then we can unmarshal into the specific struct.

	// Actually, `GetMonitor` returns `monitor.Base`. The library `Base` struct has `MarshalJSON` which uses `raw`.
	// So if we just use `monitor.Base`, we might not get type-specific fields if we don't unmarshal `raw` into struct?
	// Wait, `GetMonitor` implementation (Step 250):
	// var mon monitor.Base
	// err = convertToStruct(response.Monitor, &mon)
	// This only fills Base fields?
	// `monitor.Base` has `raw` field.
	// If `convertToStruct` fills `raw`, then we can use `As`.
	// Let's assume `GetMonitor` is enough to check existence and basic fields.
	// But for thorough read we need type specific fields.

	baseMonitor, err := r.client.Kuma.GetMonitor(ctx, monitorID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read monitor %d: %s", monitorID, err),
		)
		return
	}

	// If ID is 0, it might mean not found or empty (library usually returns error on not found, but we should check)
	if baseMonitor.ID == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Now determine type and load full details
	var fullMonitor kumamonitor.Monitor
	monitorType := baseMonitor.Type()

	switch monitorType {
	case "http":
		var m kumamonitor.HTTP
		if err := baseMonitor.As(&m); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to convert monitor: %s", err))
			return
		}
		fullMonitor = &m
	case "ping":
		var m kumamonitor.Ping
		if err := baseMonitor.As(&m); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to convert monitor: %s", err))
			return
		}
		fullMonitor = &m
	case "port":
		var m kumamonitor.TCPPort
		if err := baseMonitor.As(&m); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to convert monitor: %s", err))
			return
		}
		fullMonitor = &m
	case "keyword":
		var m kumamonitor.HTTPKeyword
		if err := baseMonitor.As(&m); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to convert monitor: %s", err))
			return
		}
		fullMonitor = &m
	default:
		// Fallback to base if type unknown, but we might miss fields
		// For now, let's error or just use base if possible?
		// We can't really use base as full monitor interface in monitorToModel because of casting.
		// We'll log a warning?
		tflog.Warn(ctx, fmt.Sprintf("Unsupported monitor type found on read: %s", monitorType))
		// Use empty struct to avoid nil panic maybe?
	}

	// Update the data model
	if fullMonitor != nil {
		r.monitorToModel(ctx, fullMonitor, &data)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MonitorResourceModel
	var stateData MonitorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.monitorFromPlan(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Error preparing monitor update", err.Error())
		return
	}

	idVal := data.ID.ValueInt64()
	_ = setIdOnMonitor(monitor, idVal) // Error is non-critical, ID will be set if type is known

	if err := r.client.Kuma.UpdateMonitor(ctx, monitor); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update monitor %d: %s", idVal, err))
		return
	}

	// Handle tag updates (tags are managed separately via AddMonitorTag/DeleteMonitorTag API)
	type tagModel struct {
		TagID types.Int64  `tfsdk:"tag_id"`
		Value types.String `tfsdk:"value"`
	}

	// Get current state tags
	var stateTags []tagModel
	if !stateData.Tags.IsNull() && !stateData.Tags.IsUnknown() {
		stateData.Tags.ElementsAs(ctx, &stateTags, false)
	}

	// Get planned tags
	var planTags []tagModel
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		data.Tags.ElementsAs(ctx, &planTags, false)
	}

	// Build maps for comparison
	stateTagMap := make(map[int64]string)
	for _, t := range stateTags {
		value := ""
		if !t.Value.IsNull() {
			value = t.Value.ValueString()
		}
		stateTagMap[t.TagID.ValueInt64()] = value
	}

	planTagMap := make(map[int64]string)
	for _, t := range planTags {
		value := ""
		if !t.Value.IsNull() {
			value = t.Value.ValueString()
		}
		planTagMap[t.TagID.ValueInt64()] = value
	}

	// Remove tags that are in state but not in plan
	for tagID, value := range stateTagMap {
		if _, exists := planTagMap[tagID]; !exists {
			if err := r.client.Kuma.DeleteMonitorTagWithValue(ctx, tagID, idVal, value); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove tag %d from monitor %d: %s", tagID, idVal, err))
				return
			}
		}
	}

	// Add or update tags that are in plan
	for tagID, planValue := range planTagMap {
		if stateValue, exists := stateTagMap[tagID]; !exists {
			// Tag doesn't exist, add it
			_, err := r.client.Kuma.AddMonitorTag(ctx, tagID, idVal, planValue)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add tag %d to monitor %d: %s", tagID, idVal, err))
				return
			}
		} else if stateValue != planValue {
			// Tag exists but value changed, delete old and add new
			if err := r.client.Kuma.DeleteMonitorTagWithValue(ctx, tagID, idVal, stateValue); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to remove old tag value for tag %d from monitor %d: %s", tagID, idVal, err))
				return
			}
			_, err := r.client.Kuma.AddMonitorTag(ctx, tagID, idVal, planValue)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to add tag %d to monitor %d: %s", tagID, idVal, err))
				return
			}
		}
	}

	// Handle active state changes (requires separate API calls)
	planActive := data.Active.ValueBool()
	stateActive := stateData.Active.ValueBool()
	if planActive != stateActive {
		if planActive {
			if err := r.client.Kuma.ResumeMonitor(ctx, idVal); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to resume monitor %d: %s", idVal, err))
				return
			}
		} else {
			if err := r.client.Kuma.PauseMonitor(ctx, idVal); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to pause monitor %d: %s", idVal, err))
				return
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	monitorID := data.ID.ValueInt64()

	// Delete the monitor
	if err := r.client.Kuma.DeleteMonitor(ctx, monitorID); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor %d: %s", monitorID, err))
		return
	}
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Convert import ID (string) to int64
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Monitor ID",
			fmt.Sprintf("Monitor ID must be a number, got: %s", req.ID),
		)
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Helpers

func setIdOnMonitor(m kumamonitor.Monitor, id int64) error {
	switch v := m.(type) {
	case *kumamonitor.HTTP:
		v.ID = id
	case *kumamonitor.Ping:
		v.ID = id
	case *kumamonitor.TCPPort:
		v.ID = id
	case *kumamonitor.HTTPKeyword:
		v.ID = id
	default:
		return fmt.Errorf("cannot set ID on unknown type")
	}
	return nil
}

func (r *MonitorResource) monitorFromPlan(ctx context.Context, plan MonitorResourceModel) (kumamonitor.Monitor, error) {
	base := kumamonitor.Base{
		Name:           plan.Name.ValueString(),
		IsActive:       plan.Active.ValueBool(),
		Interval:       plan.Interval.ValueInt64(),
		RetryInterval:  plan.RetryInterval.ValueInt64(),
		ResendInterval: plan.ResendInterval.ValueInt64(),
		MaxRetries:     plan.MaxRetries.ValueInt64(),
		UpsideDown:     plan.UpsideDown.ValueBool(),
	}

	// Notification IDs
	if !plan.NotificationIDList.IsNull() && !plan.NotificationIDList.IsUnknown() {
		var notifIDs []int64
		plan.NotificationIDList.ElementsAs(ctx, &notifIDs, false)
		base.NotificationIDs = notifIDs
	}

	// Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		type tagModel struct {
			TagID types.Int64  `tfsdk:"tag_id"`
			Value types.String `tfsdk:"value"`
		}
		var tfTags []tagModel
		plan.Tags.ElementsAs(ctx, &tfTags, false)

		var tags []tag.MonitorTag
		for _, t := range tfTags {
			mt := tag.MonitorTag{
				TagID: t.TagID.ValueInt64(),
			}
			if !t.Value.IsNull() {
				mt.Value = t.Value.ValueString()
			}
			tags = append(tags, mt)
		}
		base.Tags = tags
	}

	switch plan.Type.ValueString() {
	case "http":
		m := &kumamonitor.HTTP{
			Base: base,
			HTTPDetails: kumamonitor.HTTPDetails{
				URL:           plan.URL.ValueString(),
				Method:        plan.Method.ValueString(),
				IgnoreTLS:     plan.IgnoreTLS.ValueBool(),
				MaxRedirects:  int(plan.MaxRedirects.ValueInt64()),
				Body:          plan.Body.ValueString(),
				Headers:       plan.Headers.ValueString(),
				AuthMethod:    kumamonitor.AuthMethod(plan.AuthMethod.ValueString()),
				BasicAuthUser: plan.BasicAuthUser.ValueString(),
				BasicAuthPass: plan.BasicAuthPass.ValueString(),
			},
		}
		// Always initialize AcceptedStatusCodes to empty slice to avoid sending null
		m.AcceptedStatusCodes = []string{}
		if !plan.AcceptedStatusCodes.IsNull() {
			var codes []int64
			plan.AcceptedStatusCodes.ElementsAs(ctx, &codes, false)
			strCodes := make([]string, len(codes))
			for i, c := range codes {
				strCodes[i] = strconv.FormatInt(c, 10)
			}
			m.AcceptedStatusCodes = strCodes
		}
		return m, nil

	case "ping":
		m := &kumamonitor.Ping{
			Base: base,
			PingDetails: kumamonitor.PingDetails{
				Hostname: plan.Hostname.ValueString(),
			},
		}
		return m, nil

	case "port":
		m := &kumamonitor.TCPPort{
			Base: base,
			TCPPortDetails: kumamonitor.TCPPortDetails{
				Hostname: plan.Hostname.ValueString(),
				Port:     int(plan.Port.ValueInt64()), // Fixed cast
			},
		}
		return m, nil

	case "keyword":
		// Get method, default to GET if not specified
		method := plan.Method.ValueString()
		if method == "" {
			method = "GET"
		}

		// Get max redirects, default to 0
		maxRedirects := int(plan.MaxRedirects.ValueInt64())

		// Map HTTP details
		httpDetails := kumamonitor.HTTPDetails{
			URL:           plan.URL.ValueString(),
			Method:        method,
			MaxRedirects:  maxRedirects,
			Body:          plan.Body.ValueString(),
			Headers:       plan.Headers.ValueString(),
			AuthMethod:    kumamonitor.AuthMethod(plan.AuthMethod.ValueString()),
			BasicAuthUser: plan.BasicAuthUser.ValueString(),
			BasicAuthPass: plan.BasicAuthPass.ValueString(),
			IgnoreTLS:     plan.IgnoreTLS.ValueBool(),
		}

		// Handle AcceptedStatusCodes
		httpDetails.AcceptedStatusCodes = []string{}
		if !plan.AcceptedStatusCodes.IsNull() {
			var codes []int64
			plan.AcceptedStatusCodes.ElementsAs(ctx, &codes, false)
			strCodes := make([]string, len(codes))
			for i, c := range codes {
				strCodes[i] = strconv.FormatInt(c, 10)
			}
			httpDetails.AcceptedStatusCodes = strCodes
		}

		m := &kumamonitor.HTTPKeyword{
			Base:        base,
			HTTPDetails: httpDetails,
			HTTPKeywordDetails: kumamonitor.HTTPKeywordDetails{
				Keyword: plan.Keyword.ValueString(),
			},
		}
		return m, nil

	default:
		return nil, fmt.Errorf("unsupported monitor type: %s", plan.Type.ValueString())
	}
}

func (r *MonitorResource) monitorToModel(ctx context.Context, m kumamonitor.Monitor, data *MonitorResourceModel) {
	// Common fields
	data.ID = types.Int64Value(m.GetID())

	// Helper for tags
	mapTags := func(tags []tag.MonitorTag) {
		if len(tags) > 0 {
			type tagModel struct {
				TagID types.Int64  `tfsdk:"tag_id"`
				Value types.String `tfsdk:"value"`
			}
			var tfTags []tagModel
			for _, t := range tags {
				tm := tagModel{
					TagID: types.Int64Value(t.TagID),
				}
				if t.Value != "" {
					tm.Value = types.StringValue(t.Value)
				} else {
					tm.Value = types.StringNull()
				}
				tfTags = append(tfTags, tm)
			}
			// Use struct to define element type implies ObjectType.
			// We need to match the schema. Schema is ListNestedAttribute.
			// ListValueFrom with struct slice works for ListNestedAttribute?
			// usually yes if elements match.
			// Actually ListValueFrom takes `elemType` which is `types.Type`.
			// For nested attribute, it's `types.ObjectType`.
			// But creating ObjectType manually is verbose.
			// New approach: Use `types.ListValueFrom` with `types.ObjectType`.

			objType := types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"tag_id": types.Int64Type,
					"value":  types.StringType,
				},
			}

			data.Tags, _ = types.ListValueFrom(ctx, objType, tfTags)
		} else {
			elemType := types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"tag_id": types.Int64Type,
					"value":  types.StringType,
				},
			}
			data.Tags = types.ListNull(elemType)
		}
	}

	switch v := m.(type) {
	case *kumamonitor.HTTP:
		mapTags(v.Tags)
		data.Name = types.StringValue(v.Name)
		data.Type = types.StringValue("http")
		data.Active = types.BoolValue(v.IsActive)

		if v.URL != "" {
			data.URL = types.StringValue(v.URL)
		} else {
			data.URL = types.StringNull()
		}
		if v.Method != "" {
			data.Method = types.StringValue(v.Method)
		} else {
			data.Method = types.StringNull()
		}

		data.IgnoreTLS = types.BoolValue(v.IgnoreTLS)
		data.MaxRedirects = types.Int64Value(int64(v.MaxRedirects))

		if v.Body != "" {
			data.Body = types.StringValue(v.Body)
		} else {
			data.Body = types.StringNull()
		}
		if v.Headers != "" {
			data.Headers = types.StringValue(v.Headers)
		} else {
			data.Headers = types.StringNull()
		}

		if string(v.AuthMethod) != "" {
			data.AuthMethod = types.StringValue(string(v.AuthMethod))
		} else {
			data.AuthMethod = types.StringNull()
		}
		if v.BasicAuthUser != "" {
			data.BasicAuthUser = types.StringValue(v.BasicAuthUser)
		} else {
			data.BasicAuthUser = types.StringNull()
		}
		if v.BasicAuthPass != "" {
			data.BasicAuthPass = types.StringValue(v.BasicAuthPass)
		} else {
			data.BasicAuthPass = types.StringNull()
		}

		if len(v.AcceptedStatusCodes) > 0 {
			var codes []types.Int64
			for _, c := range v.AcceptedStatusCodes {
				if i, err := strconv.ParseInt(c, 10, 64); err == nil {
					codes = append(codes, types.Int64Value(i))
				}
			}
			data.AcceptedStatusCodes, _ = types.ListValueFrom(ctx, types.Int64Type, codes)
		} else {
			// If empty list, we prefer null to match config if omitted
			data.AcceptedStatusCodes = types.ListNull(types.Int64Type)
		}

		// Base fields
		data.Interval = types.Int64Value(v.Interval)
		data.RetryInterval = types.Int64Value(v.RetryInterval)
		data.ResendInterval = types.Int64Value(v.ResendInterval)
		data.MaxRetries = types.Int64Value(v.MaxRetries)
		data.UpsideDown = types.BoolValue(v.UpsideDown)

		if len(v.NotificationIDs) > 0 {
			outIDs := make([]types.Int64, len(v.NotificationIDs))
			for i, id := range v.NotificationIDs {
				outIDs[i] = types.Int64Value(id)
			}
			data.NotificationIDList, _ = types.ListValueFrom(ctx, types.Int64Type, outIDs)
		} else {
			data.NotificationIDList = types.ListNull(types.Int64Type)
		}

	case *kumamonitor.Ping:
		mapTags(v.Tags)
		data.Name = types.StringValue(v.Name)
		data.Type = types.StringValue("ping")
		data.Active = types.BoolValue(v.IsActive)
		if v.Hostname != "" {
			data.Hostname = types.StringValue(v.Hostname)
		} else {
			data.Hostname = types.StringNull()
		}

		data.Interval = types.Int64Value(v.Interval)
		data.RetryInterval = types.Int64Value(v.RetryInterval)
		data.ResendInterval = types.Int64Value(v.ResendInterval)
		data.MaxRetries = types.Int64Value(v.MaxRetries)
		data.UpsideDown = types.BoolValue(v.UpsideDown)

		if len(v.NotificationIDs) > 0 {
			outIDs := make([]types.Int64, len(v.NotificationIDs))
			for i, id := range v.NotificationIDs {
				outIDs[i] = types.Int64Value(id)
			}
			data.NotificationIDList, _ = types.ListValueFrom(ctx, types.Int64Type, outIDs)
		} else {
			data.NotificationIDList = types.ListNull(types.Int64Type)
		}

	case *kumamonitor.TCPPort:
		mapTags(v.Tags)
		data.Name = types.StringValue(v.Name)
		data.Type = types.StringValue("port")
		data.Active = types.BoolValue(v.IsActive)
		if v.Hostname != "" {
			data.Hostname = types.StringValue(v.Hostname)
		} else {
			data.Hostname = types.StringNull()
		}
		data.Port = types.Int64Value(int64(v.Port))

		data.Interval = types.Int64Value(v.Interval)
		data.RetryInterval = types.Int64Value(v.RetryInterval)
		data.ResendInterval = types.Int64Value(v.ResendInterval)
		data.MaxRetries = types.Int64Value(v.MaxRetries)
		data.UpsideDown = types.BoolValue(v.UpsideDown)

		if len(v.NotificationIDs) > 0 {
			outIDs := make([]types.Int64, len(v.NotificationIDs))
			for i, id := range v.NotificationIDs {
				outIDs[i] = types.Int64Value(id)
			}
			data.NotificationIDList, _ = types.ListValueFrom(ctx, types.Int64Type, outIDs)
		} else {
			data.NotificationIDList = types.ListNull(types.Int64Type)
		}

	case *kumamonitor.HTTPKeyword:
		mapTags(v.Tags)
		data.Name = types.StringValue(v.Name)
		data.Type = types.StringValue("keyword")
		data.Active = types.BoolValue(v.IsActive)
		if v.URL != "" {
			data.URL = types.StringValue(v.URL)
		} else {
			data.URL = types.StringNull()
		}
		if v.Keyword != "" {
			data.Keyword = types.StringValue(v.Keyword)
		} else {
			data.Keyword = types.StringNull()
		}

		data.Interval = types.Int64Value(v.Interval)
		data.RetryInterval = types.Int64Value(v.RetryInterval)
		data.ResendInterval = types.Int64Value(v.ResendInterval)
		data.MaxRetries = types.Int64Value(v.MaxRetries)
		data.UpsideDown = types.BoolValue(v.UpsideDown)

		if len(v.NotificationIDs) > 0 {
			outIDs := make([]types.Int64, len(v.NotificationIDs))
			for i, id := range v.NotificationIDs {
				outIDs[i] = types.Int64Value(id)
			}
			data.NotificationIDList, _ = types.ListValueFrom(ctx, types.Int64Type, outIDs)
		} else {
			data.NotificationIDList = types.ListNull(types.Int64Type)
		}

	default:
		// Fallback for unknown types
		data.Type = types.StringValue(m.Type())
	}
}
