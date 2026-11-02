import { CoreApp, DataQueryRequest } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { trackRequest } from './tracking';
import type { GitHubQuery } from './types/query';

jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
  reportInteraction: jest.fn(),
}));

const pullRequestsRequest = (app: string, options: object) =>
  ({ app, targets: [{ refId: 'A', queryType: 'Pull_Requests', options }] }) as unknown as DataQueryRequest<GitHubQuery>;

describe('trackRequest', () => {
  beforeEach(() => jest.mocked(reportInteraction).mockClear());

  it('reports all fields when fields is unset', () => {
    trackRequest(pullRequestsRequest(CoreApp.Explore, {}));
    expect(reportInteraction).toHaveBeenCalledWith(
      'grafana_github_query_executed',
      expect.objectContaining({ fields: 'all', excludedFields: 'none' })
    );
  });

  it('reports the selected and excluded fields in a fixed order', () => {
    trackRequest(pullRequestsRequest(CoreApp.Explore, { fields: ['labels', 'author'] }));
    expect(reportInteraction).toHaveBeenCalledWith(
      'grafana_github_query_executed',
      expect.objectContaining({ fields: 'author,labels', excludedFields: 'diffstat,mergeable,merged_by' })
    );
  });

  it('reports no fields when every field is cleared', () => {
    trackRequest(pullRequestsRequest(CoreApp.Explore, { fields: [] }));
    expect(reportInteraction).toHaveBeenCalledWith(
      'grafana_github_query_executed',
      expect.objectContaining({ fields: 'none', excludedFields: 'diffstat,mergeable,author,merged_by,labels' })
    );
  });

  it('does not report queries run by dashboards', () => {
    trackRequest(pullRequestsRequest(CoreApp.Dashboard, { fields: [] }));
    expect(reportInteraction).not.toHaveBeenCalled();
  });
});
