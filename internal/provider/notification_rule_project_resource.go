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
	_ resource.Resource              = &notificationRuleProjectResource{}
	_ resource.ResourceWithConfigure = &notificationRuleProjectResource{}
)

type (
	notificationRuleProjectResource struct {
		client *dtrack.Client
		semver *Semver
	}

	notificationRuleProjectResourceModel struct {
		RuleID    types.String `tfsdk:"rule"`
		ProjectID types.String `tfsdk:"project"`
	}
)

func NewNotificationRuleProjectResource() resource.Resource {
	return &notificationRuleProjectResource{}
}

func (*notificationRuleProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule_project"
}

func (*notificationRuleProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an association of a Notification Rule to a Project. Only applicable for PORTFOLIO-scoped rules.",
		Attributes: map[string]schema.Attribute{
			"rule": schema.StringAttribute{
				Description: "UUID of the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				Description: "UUID of the Project to associate with the Notification Rule.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *notificationRuleProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRuleProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleCreate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	projectID, diag := TryParseUUID(plan.ProjectID, LifecycleCreate, path.Root("project"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Notification Rule Project Mapping", map[string]any{
		"rule":    ruleID.String(),
		"project": projectID.String(),
	})

	rule, err := r.client.Notification.AddProjectToRule(ctx, ruleID, projectID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification rule project mapping",
			"Error from: "+err.Error(),
		)
		return
	}
	plan = notificationRuleProjectResourceModel{
		RuleID:    types.StringValue(rule.UUID.String()),
		ProjectID: types.StringValue(projectID.String()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created Notification Rule Project Mapping", map[string]any{
		"rule":    plan.RuleID.ValueString(),
		"project": plan.ProjectID.ValueString(),
	})
}

func (r *notificationRuleProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRuleProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleRead, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	projectID, diag := TryParseUUID(state.ProjectID, LifecycleRead, path.Root("project"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Notification Rule Project Mapping", map[string]any{
		"rule":    ruleID.String(),
		"project": projectID.String(),
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
	_, err = Find(rule.Projects, func(project dtrack.Project) bool {
		return project.UUID == projectID
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find project in notification rule",
			"Error from: "+err.Error(),
		)
		return
	}

	state = notificationRuleProjectResourceModel{
		RuleID:    types.StringValue(rule.UUID.String()),
		ProjectID: types.StringValue(projectID.String()),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Read Notification Rule Project Mapping", map[string]any{
		"rule":    state.RuleID.ValueString(),
		"project": state.ProjectID.ValueString(),
	})
}

func (*notificationRuleProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRuleProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(plan.RuleID, LifecycleUpdate, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	projectID, diag := TryParseUUID(plan.ProjectID, LifecycleUpdate, path.Root("project"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Notification Rule Project Mapping", map[string]any{
		"rule":    ruleID.String(),
		"project": projectID.String(),
	})

	plan = notificationRuleProjectResourceModel{
		RuleID:    types.StringValue(ruleID.String()),
		ProjectID: types.StringValue(projectID.String()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated Notification Rule Project Mapping", map[string]any{
		"rule":    plan.RuleID.ValueString(),
		"project": plan.ProjectID.ValueString(),
	})
}

func (r *notificationRuleProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRuleProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID, diag := TryParseUUID(state.RuleID, LifecycleDelete, path.Root("rule"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	projectID, diag := TryParseUUID(state.ProjectID, LifecycleDelete, path.Root("project"))
	if diag != nil {
		resp.Diagnostics.Append(diag)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Notification Rule Project Mapping", map[string]any{
		"rule":    ruleID.String(),
		"project": projectID.String(),
	})

	_, err := r.client.Notification.RemoveProjectFromRule(ctx, ruleID, projectID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete notification rule project mapping",
			"Error from: "+err.Error(),
		)
	}
	tflog.Debug(ctx, "Deleted Notification Rule Project Mapping", map[string]any{
		"rule":    state.RuleID.ValueString(),
		"project": state.ProjectID.ValueString(),
	})
}

func (r *notificationRuleProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
