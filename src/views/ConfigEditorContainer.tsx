import React, { PureComponent } from 'react';
import ConfigEditor from './ConfigEditor';
import { ConfigEditorProps } from '../types';

export default class extends PureComponent<ConfigEditorProps> {
  render() {
    return (
      <>
        <ConfigEditor {...this.props} />
      </>
    );
  }
}
