import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { DataSourceOptions, Label } from './types';
import { GitHubQuery } from './query';

export class DataSource extends DataSourceWithBackend<GitHubQuery, DataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
  }

  async getLabels(repository: string, owner: string, query?: string): Promise<Label[]> {
    return this.getResource('/labels', {
      repository,
      owner,
      query,
    });
  }
}
