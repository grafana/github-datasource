import React from 'react';
import QueryEditorIssues, { DefaultPackageType } from './QueryEditorPackages';
import { render } from '@testing-library/react';
import { PackageType } from './../constants';

describe('QueryEditorPackages', () => {
  it('should update package type to default one if no package is selected', async () => {
    const props = {
      onChange: jest.fn(),
      packageType: undefined,
    };
    render(<QueryEditorIssues {...props} />);
    expect(props.onChange).toHaveBeenCalledTimes(1);
    expect(props.onChange).toHaveBeenCalledWith({ packageType: DefaultPackageType, onChange: props.onChange });
  });
  it('should not update package type to default one if package type is provided ', async () => {
    const props = {
      onChange: jest.fn(),
      packageType: PackageType.DOCKER,
    };
    render(<QueryEditorIssues {...props} />);
    expect(props.onChange).not.toHaveBeenCalled();
  });
});
