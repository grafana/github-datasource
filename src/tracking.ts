import { DataQueryRequest } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { GitHubQuery, IssueTimeField, PullRequestTimeField, WorkflowsTimeField } from 'types';

export const trackRequest = (request: DataQueryRequest<GitHubQuery>) => {
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
