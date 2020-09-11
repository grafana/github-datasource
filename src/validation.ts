import { GitHubQuery } from './types';

export const isValid = (query: GitHubQuery): boolean => {
  // The current requirement is that the query has a querytype
  // TODO: have each option implement a validation function
  return !!query.queryType;
};
