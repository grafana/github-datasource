import { isEmpty } from 'lodash';
import { GitHubQuery, ProjectQueryType, QueryType } from './types';

export const isValid = (query: GitHubQuery): boolean => {
  if (query.queryType === QueryType.Repositories) {
    if (isEmpty(query.owner)) {
      return false;
    }
  }
  if (
    query.queryType === QueryType.Commits ||
    query.queryType === QueryType.Contributors ||
    query.queryType === QueryType.Tags ||
    query.queryType === QueryType.Releases ||
    query.queryType === QueryType.Labels ||
    query.queryType === QueryType.Milestones ||
    query.queryType === QueryType.Vulnerabilities ||
    query.queryType === QueryType.Stargazers
  ) {
    if (isEmpty(query.owner) || isEmpty(query.repository)) {
      return false;
    }
  }
  if (query.queryType === QueryType.Projects) {
    if (isEmpty(query.options?.user) && query.options?.kind === ProjectQueryType.USER) {
      return false;
    }
    if (isEmpty(query.options?.organization)) {
      return false;
    }
  }
  return !!query.queryType;
};
