import './common/StandardAnnotationQueryEditor';
import { DataSource } from './DataSource';
import { AnnotationQuery } from 'common/types';
// import { defaultsDeep } from 'lodash';

// const defaultQuery: GitHubAnnotationQuery = {
//   queryType: DefaultQueryType,
//   refId: '',
// };

export default class AnnotationCtrl {
  // @ts-ignore
  annotation: AnnotationQuery;

  // @ts-ignore
  private datasource?: DataSource;

  static templateUrl = 'partials/annotations.editor.html';

  /** @ngInject */
  constructor() {
    // // @ts-ignore
    // this.annotation.annotation = defaultsDeep(this.annotation.annotation, defaultQuery);
    // // @ts-ignore
    // this.annotation.datasourceId = this.datasource.id;
  }

  onChange = (anno: AnnotationQuery<any>) => {
    this.annotation.target = anno.target;
    this.annotation.mappings = anno.mappings;
  };
}
