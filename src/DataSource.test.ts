import { DataQueryResponse, DataSourceInstanceSettings, toDataFrame } from '@grafana/data';
import { GitHubVariableQuery } from 'types';
import { of } from 'rxjs';
import { GithubDataSource } from 'DataSource';

describe('DataSource', () => {
  describe('metricFindQuery', () => {
    it('should return empty array if data in response is empty array', async () => {
      const ds = new GithubDataSource({} as DataSourceInstanceSettings);
      const query = {} as GitHubVariableQuery;

      jest.spyOn(ds, 'query').mockReturnValue(of({ data: [] }));
      const res = await ds.metricFindQuery(query, {});
      expect(res).toEqual([]);
    });

    it('should return empty array if no data in response', async () => {
      const ds = new GithubDataSource({} as DataSourceInstanceSettings);
      const query = {} as GitHubVariableQuery;

      jest.spyOn(ds, 'query').mockReturnValue(of({} as DataQueryResponse));
      const res = await ds.metricFindQuery(query, {});
      expect(res).toEqual([]);
    });

    it('should return array with values if response has data', async () => {
      const ds = new GithubDataSource({} as DataSourceInstanceSettings);
      const query = { key: 'test', field: 'test' } as GitHubVariableQuery;

      jest.spyOn(ds, 'query').mockReturnValue(
        of({
          data: [
            toDataFrame({
              fields: [{ name: 'test', values: ['value1', 'value2'] }],
            }),
          ],
        } as DataQueryResponse)
      );
      const res = await ds.metricFindQuery(query, {});
      expect(res).toEqual([
        { value: 'value1', text: 'value1' },
        { value: 'value2', text: 'value2' },
      ]);
    });
  });
});
