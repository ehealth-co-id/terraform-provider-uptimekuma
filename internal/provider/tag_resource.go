// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	kumatag "github.com/breml/go-uptime-kuma-client/tag"
	"github.com/ehealth-co-id/terraform-provider-uptimekuma/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the resource implementation.
type TagResource struct {
	client *client.Client
}

// TagResourceModel describes the resource data model.
type TagResourceModel struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Uptime Kuma Tag resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Tag identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Tag name",
				Required:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Tag color in hex format (e.g., #FF0000, #00FF00)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^#([0-9A-Fa-f]{3}){1,2}$`),
						"must be a valid hex color code (e.g., #FFF or #FFFFFF)",
					),
				},
			},
		},
	}
}

func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TagResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag := kumatag.Tag{
		Name:  data.Name.ValueString(),
		Color: data.Color.ValueString(),
	}

	// Create the tag
	id, err := r.client.Kuma.CreateTag(ctx, tag)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tag: %s", err))
		return
	}

	// Update Terraform state
	data.ID = types.Int64Value(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TagResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagID := data.ID.ValueInt64()

	// Read the tag from the API
	tag, err := r.client.Kuma.GetTag(ctx, tagID)
	if err != nil {
		// Go client checks err == ErrNotFound could be useful but text check is fallback
		// If error indicates not found...
		// Library GetTag returns specific error wrapped.
		// For now simple error check.
		// If "not found" in error string or ID is 0?
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read tag %d: %s", tagID, err),
		)
		return
	}

	// If ID is zero? (Shouldn't happen with valid GetTag return)

	data.Name = types.StringValue(tag.Name)
	if tag.Color != "" {
		data.Color = types.StringValue(tag.Color)
	} else {
		data.Color = types.StringNull()
	}
	data.ID = types.Int64Value(tag.ID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TagResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag := kumatag.Tag{
		ID:    data.ID.ValueInt64(),
		Name:  data.Name.ValueString(),
		Color: data.Color.ValueString(),
	}

	// Update the tag
	if err := r.client.Kuma.UpdateTag(ctx, tag); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tag %d: %s", tag.ID, err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TagResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagID := data.ID.ValueInt64()

	// Delete the tag
	if err := r.client.Kuma.DeleteTag(ctx, tagID); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tag %d: %s", tagID, err))
		return
	}
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Convert import ID (string) to int64
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tag ID",
			fmt.Sprintf("Tag ID must be a number, got: %s", req.ID),
		)
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
