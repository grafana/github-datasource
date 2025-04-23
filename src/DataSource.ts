import {
  AnnotationEvent,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { replaceVariables, GithubVariableSupport } from './variables';
import { isValid } from './validation';
import { getAnnotationsFromFrame } from 'common/annotationsFromDataFrame';
import { prepareAnnotation } from 'migrations';
import { Observable } from 'rxjs';
import { trackRequest } from 'tracking';
import type { GitHubQuery } from './types/query';
import type { GitHubDataSourceOptions } from './types/config';

export class GitHubDataSource extends DataSourceWithBackend<GitHubQuery, GitHubDataSourceOptions> {
  templateSrv = getTemplateSrv();

  constructor(instanceSettings: DataSourceInstanceSettings<GitHubDataSourceOptions>) {
    super(instanceSettings);
    this.annotations = {
      prepareAnnotation,
    };
    this.variables = new GithubVariableSupport(this);
  }

  // Required by DataSourceApi. It executes queries based on the provided DataQueryRequest.
  query(request: DataQueryRequest<GitHubQuery>): Observable<DataQueryResponse> {
    trackRequest(request);
    return super.query(request);
  }

  // Implemented as a part of DataSourceApi
  // Only execute queries that have a query type
  filterQuery = (query: GitHubQuery) => {
    return isValid(query) && !query.hide;
  };

  // Implemented as a part of DataSourceApi. Interpolates variables and adds ad hoc filters to a list of GitHub queries.
  applyTemplateVariables(query: GitHubQuery, scoped: ScopedVars): GitHubQuery {
    return replaceVariables(this.templateSrv, query, scoped);
  }

  /**
   * Implemented as part of the DataSourceAPI. It allows the datasource to serve as a source of annotations for a dashboard.
   * @returns A promise that resolves to an array of AnnotationEvent objects representing the annotations for the dashboard.
   * @todo This is deprecated and it is recommended to use the `AnnotationSupport` feature for annotations.
   */
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
        time: annotation.timeField, // or first time field
        text: annotation.field || 'name',
      },
    });
  }

  // Used in VariableQueryEditor to get the choices for variables
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
}
