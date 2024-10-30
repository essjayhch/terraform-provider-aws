// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//Code generated by tools/tfsdk2fw/main.go. Manual editing is required.

package iot

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	awstypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_iot_billing_group")
func newResourceBillingGroup(context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceBillingGroup{}

	return r, nil
}

const (
	ResNameBillingGroup = "Billing Group"
)

type resourceBillingGroup struct {
	framework.ResourceWithConfigure
}

// Metadata should return the full name of the resource, such as
// examplecloud_thing.
func (r *resourceBillingGroup) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "aws_iot_billing_group"
}

// Schema returns the schema for this resource.
func (r *resourceBillingGroup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"id": framework.IDAttribute(),
			"metadata": schema.ListAttribute{
				CustomType:  fwtypes.NewListNestedObjectTypeOf[metadataModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[metadataModel](ctx),
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// TODO Validate*,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"tags":     tftags.TagsAttribute(),
			"tags_all": tftags.TagsAttributeComputedOnly(),
			"version": schema.Int64Attribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"properties": schema.ListNestedBlock{
				CustomType: fwtypes.NewListNestedObjectTypeOf[propertiesModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}

	response.Schema = s
}

// Create is called when the provider must create a new resource.
// Config and planned state values should be read from the CreateRequest and new state values set on the CreateResponse.
func (r *resourceBillingGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	conn := r.Meta().IoTClient(ctx)
	var data resourceBillingGroupData

	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(data.Name.ValueString())
	input := &iot.CreateBillingGroupInput{
		Tags: getTagsIn(ctx),
	}
	response.Diagnostics.Append(flex.Expand(ctx, data, input, flex.WithFieldNamePrefix("BillingGroup"))...)
	if response.Diagnostics.HasError() {
		return
	}

	out, err := conn.CreateBillingGroup(ctx, input)
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.IoT, create.ErrActionCreating, ResNameBillingGroup, data.Name.String(), err),
			err.Error(),
		)
		return
	}
	if out == nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.IoT, create.ErrActionCreating, ResNameBillingGroup, data.Name.String(), nil),
			errors.New("empty output").Error(),
		)
		return
	}

	findOut, err := findBillingGroupByName(ctx, conn, data.Name.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.IoT, create.ErrActionCreating, ResNameBillingGroup, data.Name.String(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(flex.Flatten(ctx, findOut, &data, flex.WithFieldNamePrefix("BillingGroup"))...)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// Read is called when the provider must read resource values in order to update state.
// Planned state values should be read from the ReadRequest and new state values set on the ReadResponse.
func (r *resourceBillingGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.Meta().IoTClient(ctx)

	var data resourceBillingGroupData

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	out, err := findBillingGroupByName(ctx, conn, data.Name.ValueString())
	if tfresource.NotFound(err) {
		response.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.IoT, create.ErrActionReading, ResNameBillingGroup, data.Name.String(), err),
			err.Error(),
		)
		return
	}

	response.Diagnostics.Append(flex.Flatten(ctx, out, &data, flex.WithFieldNamePrefix("BillingGroup"))...)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource.
// Config, planned state, and prior state values should be read from the UpdateRequest and new state values set on the UpdateResponse.
func (r *resourceBillingGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	conn := r.Meta().IoTClient(ctx)
	var old, new resourceBillingGroupData

	response.Diagnostics.Append(request.State.Get(ctx, &old)...)

	response.Diagnostics.Append(request.Plan.Get(ctx, &new)...)

	if response.Diagnostics.HasError() {
		return
	}

	if !old.Properties.Equal(new.Properties) {
		input := &iot.UpdateBillingGroupInput{}
		response.Diagnostics.Append(flex.Expand(ctx, new, input, flex.WithFieldNamePrefix("BillingGroup"))...)
		if response.Diagnostics.HasError() {
			return
		}

		input.ExpectedVersion = new.Version.ValueInt64Pointer()

		out, err := conn.UpdateBillingGroup(ctx, input)
		if err != nil {
			response.Diagnostics.AddError(
				create.ProblemStandardMessage(names.IoT, create.ErrActionUpdating, ResNameBillingGroup, new.Name.String(), err),
				err.Error(),
			)
			return
		}

		response.Diagnostics.Append(flex.Flatten(ctx, out, &new)...)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &new)...)
}

// Delete is called when the provider must delete the resource.
// Config values may be read from the DeleteRequest.
//
// If execution completes without error, the framework will automatically call DeleteResponse.State.RemoveResource(),
// so it can be omitted from provider logic.
func (r *resourceBillingGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	conn := r.Meta().IoTClient(ctx)

	var data resourceBillingGroupData

	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := conn.DeleteBillingGroup(ctx, &iot.DeleteBillingGroupInput{
		BillingGroupName: aws.String(data.Name.String()),
	})

	if err != nil {
		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
			return
		}
		response.Diagnostics.AddError(
			create.ProblemStandardMessage(names.IoT, create.ErrActionDeleting, ResNameBillingGroup, data.Name.String(), err),
			err.Error(),
		)
		return
	}
}

// ImportState is called when the provider must import the state of a resource instance.
// This method must return enough state so the Read method can properly refresh the full resource.
//
// If setting an attribute with the import identifier, it is recommended to use the ImportStatePassthroughID() call in this method.
func (r *resourceBillingGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

// ModifyPlan is called when the provider has an opportunity to modify
// the plan: once during the plan phase when Terraform is determining
// the diff that should be shown to the user for approval, and once
// during the apply phase with any unknown values from configuration
// filled in with their final values.
//
// The planned new state is represented by
// ModifyPlanResponse.Plan. It must meet the following
// constraints:
// 1. Any non-Computed attribute set in config must preserve the exact
// config value or return the corresponding attribute value from the
// prior state (ModifyPlanRequest.State).
// 2. Any attribute with a known value must not have its value changed
// in subsequent calls to ModifyPlan or Create/Read/Update.
// 3. Any attribute with an unknown value may either remain unknown
// or take on any value of the expected type.
//
// Any errors will prevent further resource-level plan modifications.
func (r *resourceBillingGroup) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	r.SetTagsAll(ctx, request, response)
}

func findBillingGroupByName(ctx context.Context, conn *iot.Client, name string) (*iot.DescribeBillingGroupOutput, error) {
	input := &iot.DescribeBillingGroupInput{
		BillingGroupName: aws.String(name),
	}

	output, err := conn.DescribeBillingGroup(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

type resourceBillingGroupData struct {
	ARN        types.String                                     `tfsdk:"arn"`
	ID         types.String                                     `tfsdk:"id"`
	Metadata   fwtypes.ListNestedObjectValueOf[metadataModel]   `tfsdk:"metadata"`
	Name       types.String                                     `tfsdk:"name"`
	Tags       tftags.Map                                       `tfsdk:"tags"`
	TagsAll    tftags.Map                                       `tfsdk:"tags_all"`
	Version    types.Int64                                      `tfsdk:"version"`
	Properties fwtypes.ListNestedObjectValueOf[propertiesModel] `tfsdk:"properties"`
}

type propertiesModel struct {
	Description types.String `tfsdk:"description"`
}

type metadataModel struct {
	CreationDate types.String `tfsdk:"creation_date"`
}
