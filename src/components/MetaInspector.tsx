import React, { PureComponent } from 'react';
import { MetadataInspectorProps, DataFrame } from '@grafana/data';
import { DataSource } from '../DataSource';
import { SheetsQuery, SheetResponseMeta, SheetsSourceOptions } from '../types';

export type Props = MetadataInspectorProps<DataSource, SheetsQuery, SheetsSourceOptions>;

export class MetaInspector extends PureComponent<Props> {
  state = { index: 0 };

  renderInfo = (frame: DataFrame) => {
    const meta = frame.meta?.custom as SheetResponseMeta;
    if (!meta) {
      return null;
    }

    return (
      <div>
        <h3>Info</h3>
        <pre>{JSON.stringify(meta, null, 2)}</pre>
      </div>
    );
  };

  render() {
    const { data } = this.props;
    if (!data || !data.length) {
      return <div>No Data</div>;
    }
    return (
      <div>
        <h3>Google Sheets Metadata</h3>
        {data.map(frame => {
          return this.renderInfo(frame);
        })}
      </div>
    );
  }
}
