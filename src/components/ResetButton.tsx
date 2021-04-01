import React from 'react';
import { Button } from '@grafana/ui';

interface Props {
  onClick: (event: any) => void;
  disabled: boolean;
}

const ResetButton = (props: Props) => {
  return (
    <Button variant="secondary" type="button" {...props}>
      Reset
    </Button>
  );
};

export default ResetButton;
