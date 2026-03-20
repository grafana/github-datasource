import { isEmpty } from 'lodash';
import { ProjectQueryType } from './constants';
import type { GitHubQuery } from './types/query';

export const isValid = (query: GitHubQuery): boolean => {
  if (query.queryType === "Repositories" || query.queryType === "Code_Scanning") {
    if (isEmpty(query.owner)) {
      return false;
    }
  }
  if (
    query.queryType === "Commits" ||
    query.queryType === "Contributors" ||
    query.queryType === "Tags" ||
    query.queryType === "Releases" ||
    query.queryType === "Labels" ||
    query.queryType === "Milestones" ||
    query.queryType === "Vulnerabilities" ||
    query.queryType === "Stargazers"
  ) {
    if (isEmpty(query.owner) || isEmpty(query.repository)) {
      return false;
    }
  }
  if (query.queryType === "Projects") {
    if (isEmpty(query.options?.user) && query.options?.kind === ProjectQueryType.USER) {
      return false;
    }
    if (isEmpty(query.options?.organization)) {
      return false;
    }
  }
  return !!query.queryType;
};
