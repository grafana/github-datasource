import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { GithubDataSourceOptions, Label } from './types';
import { GitHubQuery } from './query';

export class DataSource extends DataSourceWithBackend<GitHubQuery, GithubDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GithubDataSourceOptions>) {
    super(instanceSettings);
  }

  async getLabels(repository: string, owner: string, query?: string): Promise<Label[]> {
    return this.getResource('labels', {
      repository,
      owner,
      query,
    });
  }
}
