import { isValid } from './validation';
import type { GitHubQuery } from './types/query';

const query: GitHubQuery = {
  refId: 'A',
  queryType: 'Merge_Queue_Timings',
  owner: 'github',
  repository: 'docs',
};

describe('Merge Queue Timings query validation', () => {
  it('accepts an owner and repository without extra options', () => {
    expect(isValid(query)).toBe(true);
  });

  it.each(['owner', 'repository'] as const)('requires %s', (key) => {
    expect(isValid({ ...query, [key]: undefined })).toBe(false);
  });
});
