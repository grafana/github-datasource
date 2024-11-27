import { CoreApp, DataQueryRequest } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { IssueTimeField, PullRequestTimeField, WorkflowsTimeField } from './constants';
import type { GitHubQuery } from 'types/query';

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
