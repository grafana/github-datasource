import { GitHubQuery } from './types';

export const isValid = (query: GitHubQuery): boolean => {
  if (!query.owner || !query.repository) {
    return false;
  }
  return true;
};
