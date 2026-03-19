package provider

import (
	"context"
	"fmt"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &notificationPublisherDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationPublisherDataSource{}
)

type (
	notificationPublisherDataSource struct {
		client *dtrack.Client
		semver *Semver
	}

	notificationPublisherDataSourceModel struct {
		ID               types.String `tfsdk:"id"`
		Name             types.String `tfsdk:"name"`
		Description      types.String `tfsdk:"description"`
		PublisherClass   types.String `tfsdk:"publisher_class"`
		DefaultPublisher types.Bool   `tfsdk:"default_publisher"`
	}
)

func NewNotificationPublisherDataSource() datasource.DataSource {
	return &notificationPublisherDataSource{}
}

func (*notificationPublisherDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_publisher"
}

func (*notificationPublisherDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Notification Publisher by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Notification Publisher to look up. " +
					"Built-in publishers include: Slack, Microsoft Teams, Mattermost, Email, Console, Outbound Webhook, Cisco Webex, Jira.",
				Required: true,
			},
			"id": schema.StringAttribute{
				Description: "UUID of the Notification Publisher.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the Notification Publisher.",
				Computed:    true,
			},
			"publisher_class": schema.StringAttribute{
				Description: "Fully-qualified class name of the publisher implementation.",
				Computed:    true,
			},
			"default_publisher": schema.BoolAttribute{
				Description: "Whether this is a default built-in publisher.",
				Computed:    true,
			},
		},
	}
}

func (d *notificationPublisherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notificationPublisherDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	tflog.Debug(ctx, "Reading Notification Publisher", map[string]any{
		"name": name,
	})

	publishers, err := d.client.Notification.GetAllPublishers(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch notification publishers",
			"Error from: "+err.Error(),
		)
		return
	}

	publisher, err := Find(publishers, func(p dtrack.NotificationPublisher) bool {
		return p.Name == name
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find notification publisher",
			fmt.Sprintf("No publisher found with name '%s'. Error: %s", name, err.Error()),
		)
		return
	}

	state = notificationPublisherDataSourceModel{
		ID:               types.StringValue(publisher.UUID.String()),
		Name:             types.StringValue(publisher.Name),
		Description:      types.StringValue(publisher.Description),
		PublisherClass:   types.StringValue(publisher.PublisherClass),
		DefaultPublisher: types.BoolValue(publisher.DefaultPublisher),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Read Notification Publisher", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (d *notificationPublisherDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = clientInfoData.client
	d.semver = clientInfoData.semver
}
