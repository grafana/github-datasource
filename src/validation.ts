import { isEmpty } from 'lodash';
import { GitHubQuery, QueryType } from './types';

export const isValid = (query: GitHubQuery): boolean => {
  // The current requirement is that the query has a querytype
  // TODO: have each option implement a validation function
  if (query.queryType === QueryType.Projects) {
    if (isEmpty(query.options?.organization) && isEmpty(query.options?.user)) {
      return false;
    }
  }
  return !!query.queryType;
};
