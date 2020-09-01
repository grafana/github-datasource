import './views/AnnotationQueryEditor';
import { GitHubAnnotationQuery, DefaultQueryType } from './types';
import { DataSource } from './DataSource';
import { AnnotationQueryRequest } from '@grafana/data';
import { defaultsDeep } from 'lodash';

const defaultQuery: GitHubAnnotationQuery = {
  queryType: DefaultQueryType,
  refId: '',
};

export default class AnnotationCtrl {
  // @ts-ignore
  annotation: AnnotationQueryRequest<GitHubAnnotationQuery>;

  // @ts-ignore
  private datasource?: DataSource;

  static templateUrl = 'partials/annotations.editor.html';

  /** @ngInject */
  constructor() {
    // @ts-ignore
    this.annotation.annotation = defaultsDeep(this.annotation.annotation, defaultQuery);
    // @ts-ignore
    this.annotation.datasourceId = this.datasource.id;
  }

  onChange = (query: AnnotationQueryRequest<GitHubAnnotationQuery>) => {
    this.annotation.annotation = query.annotation;
  };
}
