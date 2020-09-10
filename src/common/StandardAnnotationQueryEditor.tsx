import React, { PureComponent } from 'react';

import { DataQuery, LoadingState, DataSourceApi, TimeRange } from '@grafana/data';
import { Spinner, Icon, IconName, Button } from '@grafana/ui';

import { cx, css } from 'emotion';
import { standardAnnotationSupport, executeAnnotationQuery } from './standardAnnotationSupport';
import { AnnotationQueryResponse, AnnotationEventMappings, AnnotationQuery } from './types';
import { AnnotationFieldMapper } from './AnnotationResultMapper';
import coreModule from 'grafana/app/core/core_module';
import { getLegacyAngularInjector } from '@grafana/runtime';

interface Props {
  datasource: DataSourceApi;
  annotation: AnnotationQuery<DataQuery>;
  change: (annotation: AnnotationQuery<DataQuery>) => void;
}

interface State {
  running?: boolean;
  response?: AnnotationQueryResponse;
}

export default class StandardAnnotationQueryEditor extends PureComponent<Props, State> {
  state = {} as State;

  componentDidMount() {
    this.verifyDataSource();
  }

  componentDidUpdate(oldProps: Props) {
    if (this.props.annotation !== oldProps.annotation) {
      this.verifyDataSource();
    }
  }

  verifyDataSource() {
    const { annotation } = this.props;

    // Handle any migration issues
    const processor = standardAnnotationSupport;

    const fixed = processor.prepareAnnotation!(annotation);
    if (fixed !== annotation) {
      this.props.change(fixed);
    } else {
      this.onRunQuery();
    }
  }

  onRunQuery = async () => {
    const { datasource, annotation } = this.props;
    this.setState({
      running: true,
    });
    const response = await executeAnnotationQuery(
      {
        range: this.getDashboardTime(),
        panel: {},
        dashboard: {}, //getDashboardSrv().getCurrent(),
      },
      datasource,
      annotation
    ).toPromise();
    this.setState({
      running: false,
      response,
    });
  };

  getDashboardTime = (): TimeRange => {
    const $injector = getLegacyAngularInjector();
    const timeSrv = $injector.get('timeSrv');
    return timeSrv.timeRange();
  };

  onQueryChange = (target: DataQuery) => {
    this.props.change({
      ...this.props.annotation,
      target,
    });
  };

  onMappingChange = (mappings: AnnotationEventMappings) => {
    this.props.change({
      ...this.props.annotation,
      mappings,
    });
  };

  renderStatus() {
    const { response, running } = this.state;
    let rowStyle = 'alert-info';
    let text = '...';
    let icon: IconName | undefined = undefined;

    if (running || response?.panelData?.state === LoadingState.Loading || !response) {
      text = 'loading...';
    } else {
      const { events, panelData, frame } = response;

      if (panelData?.error) {
        rowStyle = 'alert-error';
        icon = 'exclamation-triangle';
        text = panelData.error.message ?? 'error';
      } else if (!events?.length) {
        rowStyle = 'alert-warning';
        icon = 'exclamation-triangle';
        text = 'No events found';
        if (frame?.fields?.length) {
          text += ` (from ${frame?.fields.length} fields)`;
        }
      } else {
        text = `${events.length} events (from ${frame?.fields.length} fields)`;
      }
    }
    return (
      <div
        className={cx(
          rowStyle,
          css`
            margin: 4px 0px;
            padding: 4px;
            display: flex;
            justify-content: space-between;
            align-items: center;
          `
        )}
      >
        <div>
          {icon && (
            <>
              <Icon name={icon} />
              &nbsp;
            </>
          )}
          {text}
        </div>
        <div>
          {running ? (
            <Spinner />
          ) : (
            <Button variant="secondary" size="xs" onClick={this.onRunQuery}>
              TEST
            </Button>
          )}
        </div>
      </div>
    );
  }

  render() {
    const { datasource, annotation } = this.props;
    const { response } = this.state;

    // Find the annotaiton runner
    let QueryEditor = datasource.components?.QueryEditor;
    if (!QueryEditor) {
      return <div>Annotations are not supported. This datasource needs to export a QueryEditor</div>;
    }

    const query = annotation.target ?? { refId: 'Anno' };
    return (
      <>
        <QueryEditor
          key={datasource?.name}
          query={query}
          datasource={datasource}
          onChange={this.onQueryChange}
          onRunQuery={this.onRunQuery}
          data={response?.panelData}
        />
        {this.renderStatus()}

        <AnnotationFieldMapper
          response={response}
          mappings={annotation.mappings}
          change={this.onMappingChange as any}
        />
        <br />
      </>
    );
  }
}

// Careful to use a unique directive name!  many plugins already use "annotationEditor" and have conflicts
coreModule.directive('githubAnnotationEditor', [
  'reactDirective',
  (reactDirective: any) => {
    return reactDirective(StandardAnnotationQueryEditor, ['annotation', 'datasource', 'change']);
  },
]);
