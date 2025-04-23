import { DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings, toDataFrame } from '@grafana/data';
import { lastValueFrom, of } from 'rxjs';
import { GithubVariableSupport } from 'variables';
import { GitHubDataSource } from 'DataSource';
import type { GitHubVariableQuery } from 'types/query';

describe('DataSource', () => {
  describe('GithubVariableSupport', () => {
    const SAMPLE_RESPONSE_WITH_MULTIPLE_FIELDS = [
      toDataFrame({
        fields: [
          { name: 'test', values: ['value1', 'value2'] },
          { name: 'foo', values: ['foo1', 'foo2'] },
          { name: 'bar', values: ['bar1', 'bar2'] },
        ],
      }),
    ];
    it('should return empty array if data in response is empty array', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = {} as GitHubVariableQuery;
      jest.spyOn(ds, 'query').mockReturnValue(of({ data: [] }));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual([]);
      expect(res?.data.map((d) => d.text)).toEqual([]);
    });
    it('should return empty array if no data in response', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = {} as GitHubVariableQuery;
      jest.spyOn(ds, 'query').mockReturnValue(of({} as DataQueryResponse));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual([]);
      expect(res?.data.map((d) => d.text)).toEqual([]);
    });
    it('should return array with values if response has data', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = { key: 'test', field: 'test' } as GitHubVariableQuery;
      const data = [toDataFrame({ fields: [{ name: 'test', values: ['value1', 'value2'] }] })];
      jest.spyOn(ds, 'query').mockReturnValue(of({ data }));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual(['value1', 'value2']);
      expect(res?.data.map((d) => d.text)).toEqual(['value1', 'value2']);
    });
    it('mapping of key', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = { key: 'foo' } as GitHubVariableQuery;
      const data = SAMPLE_RESPONSE_WITH_MULTIPLE_FIELDS;
      jest.spyOn(ds, 'query').mockReturnValue(of({ data }));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual(['foo1', 'foo2']);
      expect(res?.data.map((d) => d.text)).toEqual(['foo1', 'foo2']);
    });
    it('mapping of key and field', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = { key: 'bar', field: 'foo' } as GitHubVariableQuery;
      const data = SAMPLE_RESPONSE_WITH_MULTIPLE_FIELDS;
      jest.spyOn(ds, 'query').mockReturnValue(of({ data }));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual(['bar1', 'bar2']);
      expect(res?.data.map((d) => d.text)).toEqual(['foo1', 'foo2']);
    });
    it('mapping of field', async () => {
      const ds = new GitHubDataSource({} as DataSourceInstanceSettings);
      const vs = new GithubVariableSupport(ds);
      const query = { field: 'foo' } as GitHubVariableQuery;
      const data = SAMPLE_RESPONSE_WITH_MULTIPLE_FIELDS;
      jest.spyOn(ds, 'query').mockReturnValue(of({ data }));
      const res = await lastValueFrom(vs.query({ targets: [query] } as DataQueryRequest));
      expect(res?.data.map((d) => d.value)).toEqual(['foo1', 'foo2']);
      expect(res?.data.map((d) => d.text)).toEqual(['foo1', 'foo2']);
    });
  });
});
