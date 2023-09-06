import {
  AnnotationEvent,
  DataFrame,
  DataFrameView,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  MetricFindValue,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv, reportInteraction } from '@grafana/runtime';
import {
  GithubDataSourceOptions,
  GitHubQuery,
  GitHubVariableQuery,
  IssueTimeField,
  Label,
  PullRequestTimeField,
  WorkflowsTimeField,
} from './types';
import { replaceVariables } from './variables';
import { isValid } from './validation';
import { getAnnotationsFromFrame } from 'common/annotationsFromDataFrame';
import { prepareAnnotation } from 'migrations';
import { Observable } from 'rxjs';

export class GithubDataSource extends DataSourceWithBackend<GitHubQuery, GithubDataSourceOptions> {
  templateSrv = getTemplateSrv();

  constructor(instanceSettings: DataSourceInstanceSettings<GithubDataSourceOptions>) {
    super(instanceSettings);
    this.annotations = {
      prepareAnnotation,
    };
  }

  query(request: DataQueryRequest<GitHubQuery>): Observable<DataQueryResponse> {
    request.targets.forEach((target) => {
      let properties: Partial<GitHubQuery> = { app: request.app, queryType: target.queryType };

      if (target.queryType === 'Issues') {
        properties.timeField = IssueTimeField[target.options?.timeField ?? 0];
      }

      if (target.queryType === 'Pull_Requests') {
        properties.timeField = PullRequestTimeField[target.options?.timeField ?? 0];
      }

      if (target.queryType === 'Workflows') {
        properties.timeField = WorkflowsTimeField[target.options?.timeField ?? 0];
      }

      if (target.queryType === 'Packages') {
        properties.timeField = target?.options?.packageType ?? 'NPM';
      }

      reportInteraction('grafana_github_query_executed', properties);
    });

    return super.query(request);
  }

  // Only execute queries that have a query type
  filterQuery = (query: GitHubQuery) => {
    return isValid(query) && !query.hide;
  };

  applyTemplateVariables(query: GitHubQuery, scoped: ScopedVars): Record<string, any> {
    return replaceVariables(this.templateSrv, query, scoped);
  }

  async getLabels(repository: string, owner: string, query?: string): Promise<Label[]> {
    return this.getResource('labels', {
      repository,
      owner,
      query,
    });
  }

  async annotationQuery(request: any): Promise<AnnotationEvent[]> {
    const { annotation } = request.annotation;

    const query = {
      targets: [
        {
          ...annotation,
          datasourceId: this.id,
          refId: this.name,
        },
      ],
      range: request.range,
      interval: request.interval,
    } as DataQueryRequest<GitHubQuery>;

    const res = await this.query(query).toPromise();

    if (!res?.data?.length) {
      return [];
    }
    return getAnnotationsFromFrame(res.data[0], {
      field: {
        // title: `${request.annotation.name} - ${annotation.queryType}`,
        time: annotation.timeField, // or first time field
        text: annotation.field || 'name',
      },
    });
  }

  async getChoices(query: GitHubQuery): Promise<string[]> {
    const request = {
      targets: [
        {
          ...query,
          refId: 'metricFindQuery',
        },
      ],
      range: {
        to: {},
        from: {},
      },
    } as DataQueryRequest;

    try {
      const res = await this.query(request).toPromise();
      const columns = (res?.data[0]?.fields || []).map((f: any) => f.name) || [];
      return columns;
    } catch (err) {
      return Promise.reject(err);
    }
  }

  async metricFindQuery(query: GitHubVariableQuery, options: any): Promise<MetricFindValue[]> {
    const request = {
      targets: [
        {
          ...query,
          refId: 'metricFindQuery',
        },
      ],
      range: options.range,
      rangeRaw: options.rangeRaw,
    } as DataQueryRequest;
    try {
      let res = await this.query(request).toPromise();
      if (!res || !res.data || res.data.length < 0) {
        return [];
      }
      const view = new DataFrameView(res.data[0] as DataFrame);
      return view.map((item) => {
        const value = item[query.key || ''] || item[query.field || 'name'];
        return {
          value,
          text: item[query.field || 'name'],
        };
      });
    } catch (ex) {
      return Promise.reject(ex);
    }
  }
}
