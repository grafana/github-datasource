import { isEmpty } from 'lodash';
import { GitHubQuery, ProjectQueryType, QueryType } from './types';

export const isValid = (query: GitHubQuery): boolean => {
  // The current requirement is that the query has a querytype
  // TODO: have each option implement a validation function
  if (query.queryType === QueryType.Projects) {
    if (isEmpty(query.options?.organization) && query.options?.kind === ProjectQueryType.ORG) {
      return false;
    }
    if (isEmpty(query.options?.user) && query.options?.kind === ProjectQueryType.USER) {
      return false;
    }
  }
  return !!query.queryType;
};
