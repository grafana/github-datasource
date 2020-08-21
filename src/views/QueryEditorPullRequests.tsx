import React from 'react';

import QueryEditorRepository from './QueryEditorRepository';
import { PullRequestsOptions } from '../query';

interface Props extends PullRequestsOptions {
  onChange: (value: PullRequestsOptions) => void;
};

export default (props: Props) => {
  return (
    <QueryEditorRepository {...props} />
  );
}
