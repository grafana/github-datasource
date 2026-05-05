package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/mcp"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/mcp/fromschema"

	"github.com/grafana/github-datasource/pkg/plugin"
)

const dsID = "grafana-github-datasource"

func main() {
	schema, err := loadSchema()
	if err != nil {
		backend.Logger.Error("schema load failed", "err", err)
		os.Exit(1)
	}

	mcpServer := mcp.NewServer(mcp.ServerOpts{
		Name:    dsID,
		Version: "1.0.0",
	})

	// build the instance manager once, bind to MCP, then pass to Manage.
	im := datasource.NewInstanceManager(plugin.NewDataSourceInstance)
	mgr := datasource.NewAutomanagementHandler(im)

	mcpServer.BindQueryDataHandler(mgr)
	mcpServer.BindCallResourceHandler(mgr)
	mcpServer.BindCheckHealthHandler(mgr)

	fromschema.RegisterQueryTools(mcpServer, schema)
	fromschema.RegisterRouteTools(mcpServer, schema)
	fromschema.RegisterQueryExamples(mcpServer, schema)
	fromschema.RegisterHealthCheckTool(mcpServer)

	mcpServer.RegisterPrompt(mcp.Prompt{
		Name:        "investigate-pull-requests",
		Description: "Walk through investigating recent pull requests in a repository",
		Template:    "List the most recent pull requests for the configured repository, then summarise patterns in review activity over the last 7 days.",
	})

	if err := datasource.Manage(dsID, plugin.NewDataSourceInstance, datasource.ManageOpts{
		MCPServer: mcpServer,
	}); err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
