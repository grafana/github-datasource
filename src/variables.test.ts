import { replaceVariables } from './variables';
import type { GitHubQuery } from './types/query';

describe('variables', () => {
  it('should not interpolate refId', () => {
    const t = {
      replace: jest.fn((a) => a),
      getVariables: jest.fn(),
      updateTimeRange: jest.fn(),
      containsTemplate: jest.fn(),
    };
    const query: GitHubQuery = { refId: 'A', foo: 'bar', foo1: 123.45, query: 'myQuery', foo2: 'bar2' };
    expect(replaceVariables(t, query, {})).toStrictEqual({ ...query });
    expect(t.replace).toHaveBeenCalledTimes(3);
    expect(t.replace).toHaveBeenNthCalledWith(1, 'bar', {}, undefined);
    expect(t.replace).toHaveBeenNthCalledWith(2, 'myQuery', {}, 'csv');
    expect(t.replace).toHaveBeenNthCalledWith(3, 'bar2', {}, undefined);
  });

  it('should interpolate query with options', () => {
    const t = {
      replace: jest.fn((a) => a),
      getVariables: jest.fn(),
      updateTimeRange: jest.fn(),
      containsTemplate: jest.fn(),
    };
    const query: GitHubQuery = {
      refId: 'A',
      foo: 'bar',
      foo1: 123.45,
      query: 'myQuery',
      foo2: 'bar2',
      options: { foo: 'options_bar', foo1: 123.45, query: 'options_myQuery', foo2: 'options_bar2' },
    };
    expect(replaceVariables(t, query, {})).toStrictEqual({ ...query });
    expect(t.replace).toHaveBeenCalledTimes(6);
    expect(t.replace).toHaveBeenNthCalledWith(1, 'bar', {}, undefined);
    expect(t.replace).toHaveBeenNthCalledWith(2, 'myQuery', {}, 'csv');
    expect(t.replace).toHaveBeenNthCalledWith(3, 'bar2', {}, undefined);
    expect(t.replace).toHaveBeenNthCalledWith(4, 'options_bar', {}, undefined);
    expect(t.replace).toHaveBeenNthCalledWith(5, 'options_myQuery', {}, 'csv');
    expect(t.replace).toHaveBeenNthCalledWith(6, 'options_bar2', {}, undefined);
  });

  it('should interpolate variables in options as well', () => {
    const t = {
      replace: jest.fn((a: string) => a.slice(1, a.length)),
      getVariables: jest.fn(),
      updateTimeRange: jest.fn(),
      containsTemplate: jest.fn(),
    };
    const query: GitHubQuery = {
      refId: 'A',
      foo: '$bar',
      query: 'myQuery',
      options: { foo: '$options_bar' },
    };
    const result = replaceVariables(t, query, {});
    expect(result['foo']).toBe('bar');
    expect(result.options?.['foo']).toBe('options_bar');
  });
});
