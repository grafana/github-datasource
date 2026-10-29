import { CoreApp, DataQueryRequest } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { IssueTimeField, PullRequestFields, PullRequestTimeField, WorkflowsTimeField } from './constants';
import type { GitHubQuery, Options } from 'types/query';

export const trackRequest = (request: DataQueryRequest<GitHubQuery>) => {
  if (request.app === CoreApp.Dashboard || request.app === CoreApp.PanelViewer) {
    return;
  }

  request.targets.forEach((target) => {
    let properties: Options = { app: request.app, queryType: target.queryType };

    if (target.queryType === 'Issues') {
      properties.timeField = IssueTimeField[target.options?.timeField ?? 0];
    }

    if (target.queryType === 'Pull_Requests') {
      properties.timeField = PullRequestTimeField[target.options?.timeField ?? 0];
      // Sent with the existing query event rather than on every change in the editor, so it adds no events
      const fields = target.options?.fields;
      const included = PullRequestFields.filter((f) => fields === undefined || fields.includes(f));
      properties.fields = included.length === PullRequestFields.length ? 'all' : included.join(',') || 'none';
      properties.excludedFields = PullRequestFields.filter((f) => !included.includes(f)).join(',') || 'none';
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
