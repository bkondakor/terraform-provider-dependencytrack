package provider

import (
	"context"
	"fmt"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &notificationRuleTeamResource{}
	_ resource.ResourceWithConfigure = &notificationRuleTeamResource{}
)

type (
	notificationRuleTeamResource struct {
		client *dtrack.Client
		semver *Semver
	}

	notificationRuleTeamResourceModel struct {
		RuleID types.String `tfsdk:"rule"`
		TeamID types.String `tfsdk:"team"`
	}
)

func NewNotificationRuleTeamResource() resource.Resource {
	return &notificationRuleTeamResource{}
}

func (*notificationRuleTeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule_team"
}

func (*notificationRuleTeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an association of a Notification Rule to a Team. Only applicable for rules using the Email publisher.",
		Attributes: map[string]schema.Attribute{
			"rule": schema.StringAttribute{
				Description: "UUID of the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team": schema.StringAttribute{
				Description: "UUID of the Team to associate with the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *notificationRuleTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRuleTeamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleCreate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	teamID, diag := TryParseUUID(plan.TeamID, LifecycleCreate, path.Root("team"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Notification Rule Team Mapping", map[string]any{
		"rule": ruleID.String(),
		"team": teamID.String(),
	})

	rule, err := r.client.Notification.AddTeamToRule(ctx, ruleID, teamID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification rule team mapping",
			"Error from: "+err.Error(),
		)
		return
	}
	plan = notificationRuleTeamResourceModel{
		RuleID: types.StringValue(rule.UUID.String()),
		TeamID: types.StringValue(teamID.String()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created Notification Rule Team Mapping", map[string]any{
		"rule": plan.RuleID.ValueString(),
		"team": plan.TeamID.ValueString(),
	})
}

func (r *notificationRuleTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRuleTeamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleRead, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	teamID, diag := TryParseUUID(state.TeamID, LifecycleRead, path.Root("team"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Notification Rule Team Mapping", map[string]any{
		"rule": ruleID.String(),
		"team": teamID.String(),
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
		resp.Diagnostics.AddError(
			"Unable to find notification rule",
			"Error from: "+err.Error(),
		)
		return
	}
	_, err = Find(rule.Teams, func(team dtrack.Team) bool {
		return team.UUID == teamID
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find team in notification rule",
			"Error from: "+err.Error(),
		)
		return
	}

	state = notificationRuleTeamResourceModel{
		RuleID: types.StringValue(rule.UUID.String()),
		TeamID: types.StringValue(teamID.String()),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Read Notification Rule Team Mapping", map[string]any{
		"rule": state.RuleID.ValueString(),
		"team": state.TeamID.ValueString(),
	})
}

func (*notificationRuleTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRuleTeamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleUpdate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	teamID, diag := TryParseUUID(plan.TeamID, LifecycleUpdate, path.Root("team"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Notification Rule Team Mapping", map[string]any{
		"rule": ruleID.String(),
		"team": teamID.String(),
	})

	plan = notificationRuleTeamResourceModel{
		RuleID: types.StringValue(ruleID.String()),
		TeamID: types.StringValue(teamID.String()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated Notification Rule Team Mapping", map[string]any{
		"rule": plan.RuleID.ValueString(),
		"team": plan.TeamID.ValueString(),
	})
}

func (r *notificationRuleTeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRuleTeamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleDelete, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	teamID, diag := TryParseUUID(state.TeamID, LifecycleDelete, path.Root("team"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Notification Rule Team Mapping", map[string]any{
		"rule": ruleID.String(),
		"team": teamID.String(),
	})

	_, err := r.client.Notification.RemoveTeamFromRule(ctx, ruleID, teamID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete notification rule team mapping",
			"Error from: "+err.Error(),
		)
	}
	tflog.Debug(ctx, "Deleted Notification Rule Team Mapping", map[string]any{
		"rule": state.RuleID.ValueString(),
		"team": state.TeamID.ValueString(),
	})
}

func (r *notificationRuleTeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
