package provider

import (
	"context"
	"fmt"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &notificationRuleTagResource{}
	_ resource.ResourceWithConfigure = &notificationRuleTagResource{}
)

type (
	notificationRuleTagResource struct {
		client *dtrack.Client
		semver *Semver
	}

	notificationRuleTagResourceModel struct {
		RuleID types.String `tfsdk:"rule"`
		Tag    types.String `tfsdk:"tag"`
	}
)

func NewNotificationRuleTagResource() resource.Resource {
	return &notificationRuleTagResource{}
}

func (*notificationRuleTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule_tag"
}

func (*notificationRuleTagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an association of a Notification Rule to a Tag. Only applicable for PORTFOLIO-scoped rules. Requires API version >= 4.12.",
		Attributes: map[string]schema.Attribute{
			"rule": schema.StringAttribute{
				Description: "UUID of the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tag": schema.StringAttribute{
				Description: "Name of the Tag to associate with the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *notificationRuleTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRuleTagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleCreate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	tagName := plan.Tag.ValueString()

	tflog.Debug(ctx, "Creating Notification Rule Tag Mapping", map[string]any{
		"rule": ruleID.String(),
		"tag":  tagName,
	})

	err := r.client.Tag.TagNotificationRules(ctx, tagName, []uuid.UUID{ruleID})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification rule tag mapping",
			"Error from: "+err.Error(),
		)
		return
	}
	plan = notificationRuleTagResourceModel{
		RuleID: types.StringValue(ruleID.String()),
		Tag:    types.StringValue(tagName),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created Notification Rule Tag Mapping", map[string]any{
		"rule": plan.RuleID.ValueString(),
		"tag":  plan.Tag.ValueString(),
	})
}

func (r *notificationRuleTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRuleTagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleRead, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	tagName := state.Tag.ValueString()

	tflog.Debug(ctx, "Reading Notification Rule Tag Mapping", map[string]any{
		"rule": ruleID.String(),
		"tag":  tagName,
	})

	allRules, err := dtrack.FetchAll(func(po dtrack.PageOptions) (dtrack.Page[dtrack.NotificationRule], error) {
		return r.client.Notification.GetAllRules(ctx, po, dtrack.SortOptions{}, dtrack.GetAllRulesFilterOptions{})
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch notification rules",
			"Error from: "+err.Error(),
		)
		return
	}
	rule, err := Find(allRules, func(rule dtrack.NotificationRule) bool {
		return rule.UUID == ruleID
	})
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	_, err = Find(rule.Tags, func(tag dtrack.Tag) bool {
		return tag.Name == tagName
	})
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state = notificationRuleTagResourceModel{
		RuleID: types.StringValue(rule.UUID.String()),
		Tag:    types.StringValue(tagName),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Read Notification Rule Tag Mapping", map[string]any{
		"rule": state.RuleID.ValueString(),
		"tag":  state.Tag.ValueString(),
	})
}

func (*notificationRuleTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRuleTagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleUpdate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	tagName := plan.Tag.ValueString()

	tflog.Debug(ctx, "Updating Notification Rule Tag Mapping", map[string]any{
		"rule": ruleID.String(),
		"tag":  tagName,
	})

	plan = notificationRuleTagResourceModel{
		RuleID: types.StringValue(ruleID.String()),
		Tag:    types.StringValue(tagName),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated Notification Rule Tag Mapping", map[string]any{
		"rule": plan.RuleID.ValueString(),
		"tag":  plan.Tag.ValueString(),
	})
}

func (r *notificationRuleTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRuleTagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleDelete, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
		return
	}
	tagName := state.Tag.ValueString()

	tflog.Debug(ctx, "Deleting Notification Rule Tag Mapping", map[string]any{
		"rule": ruleID.String(),
		"tag":  tagName,
	})

	err := r.client.Tag.UntagNotificationRules(ctx, tagName, []uuid.UUID{ruleID})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete notification rule tag mapping",
			"Error from: "+err.Error(),
		)
	}
	tflog.Debug(ctx, "Deleted Notification Rule Tag Mapping", map[string]any{
		"rule": state.RuleID.ValueString(),
		"tag":  state.Tag.ValueString(),
	})
}

func (r *notificationRuleTagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientInfoData, ok := req.ProviderData.(clientInfo)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected provider.clientInfo, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = clientInfoData.client
	r.semver = clientInfoData.semver
}
