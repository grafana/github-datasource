import { CoreApp, DataQueryRequest, DataSourcePluginMeta, type TestDataSourceResponse } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { GitHubQuery, IssueTimeField, PullRequestTimeField, WorkflowsTimeField } from 'types';

export const trackRequest = (request: DataQueryRequest<GitHubQuery>) => {
  if (request.app === CoreApp.Dashboard || request.app === CoreApp.PanelViewer) {
    return;
  }

  request.targets.forEach((target) => {
    let properties: Partial<GitHubQuery> = { app: request.app, queryType: target.queryType };

    if (target.queryType === 'Issues') {
      properties.timeField = IssueTimeField[target.options?.timeField ?? 0];
    }

    if (target.queryType === 'Pull_Requests') {
      properties.timeField = PullRequestTimeField[target.options?.timeField ?? 0];
    }

    if (target.queryType === 'Workflows') {
      properties.timeField = WorkflowsTimeField[target.options?.timeField ?? 0];
    }

    if (target.queryType === 'Packages') {
      properties.timeField = target?.options?.packageType ?? 'NPM';
    }

    reportInteraction('grafana_github_query_executed', properties);
  });
};

export const trackHealthCheck = (res: TestDataSourceResponse, meta: DataSourcePluginMeta) => {
  let properties: Record<string, any> = {
    'plugin.id': meta?.id || 'unknown',
    'plugin.version': meta?.info?.version || 'unknown',
    'datasource.healthcheck.status': res.status || 'unknown',
    'datasource.healthcheck.message': res.message || 'unknown',
  };
  if ((res?.status || '').toLowerCase() !== 'success') {
    console.error(`Health check failed. ${res.message}.`, JSON.stringify({ res }));
    reportInteraction('plugin_health_check_completed', properties);
  }
};
