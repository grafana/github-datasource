import { DataSourceInstanceSettings } from "@grafana/data";
import { GitHubVariableQuery } from "types";
import { of } from 'rxjs';
import { GithubDataSource } from "DataSource";

describe('DataSource', () => {
  describe('metricFindQuery', () => {
    it('should not throw an error when data in response is empty array', async () => {
      const ds = new GithubDataSource({} as DataSourceInstanceSettings);
      const query = {} as GitHubVariableQuery

      jest.spyOn(ds, 'query').mockReturnValue(of({data: []}))
      const res = await ds.metricFindQuery(query, {})
      expect(res).toEqual([])
    });
    })
  }); 