package solace

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "solace_application",
		Description: "Get a list of applications that match the given parameters",
		List: &plugin.ListConfig{
			Hydrate: listApplications,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getApplication,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "applicationType", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "brokerType", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "applicationDomainId", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "numberOfVersions", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "customAttributes", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "type", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "createdBy", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "createdTime", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "changedBy", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "updatedTime", Type: proto.ColumnType_TIMESTAMP, Description: ""},
		},
	}
}

func listApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	LogQueryContext("listApplications", ctx, d, h)
	plugin.Logger(ctx).Trace("DEBUGGING solace_application.QueryData", "QueryData", fmt.Sprintf("%+v", d.FetchType))

	client, err := NewSolaceClient(d.Connection)
	if err != nil {
		return nil, err
	}

	tlp := client.NewApplicationListPaginator(nil)
	pagesLeft := true
	count := 0
	for pagesLeft {
		applications, meta, err := tlp.NextPage()
		plugin.Logger(ctx).Trace("DEBUGGING solace_application.listApplications", "applications", applications)
		if err != nil {
			plugin.Logger(ctx).Error("solace_application.listApplications", "request_error", err)
			pagesLeft = false
			// return nil, err
		} else {
			count += meta.Pagination.Count
			plugin.Logger(ctx).Trace("RECORDS FETCHED - ", count)
		}

		// stream results
		for _, i := range applications {
			d.StreamListItem(ctx, i)

			if d.RowsRemaining(ctx) <= 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("DEBUGGING solace_application.QueryData", "QueryData", fmt.Sprintf("%+v", d.FetchType))
	c, err := NewSolaceClient(d.Connection)
	if err != nil {
		return nil, err
	}
	id := d.EqualsQualString("id")
	plugin.Logger(ctx).Trace("DEBUGGING solace_application.getApplication - ID", id)

	application, err := c.GetApplication(id)
	plugin.Logger(ctx).Trace("DEBUGGING solace_application.getApplication", "application", fmt.Sprintf("%+v", application))
	if err != nil {
		plugin.Logger(ctx).Error("DEBUGGING solace_application.getApplication", "request_error", err)
		return nil, err
	}

	return application, nil
}
