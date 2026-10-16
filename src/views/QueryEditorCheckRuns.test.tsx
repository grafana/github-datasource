import React from 'react';
import { QueryEditorCheckRuns } from './QueryEditorCheckRuns';
import { render, screen, fireEvent } from '@testing-library/react';

describe('QueryEditorCheckRuns', () => {
  it('renders the ref and check name fields', () => {
    render(<QueryEditorCheckRuns onChange={jest.fn()} />);
    expect(screen.getByText('Ref')).toBeInTheDocument();
    expect(screen.getByText('Check name')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Filter')).toBeInTheDocument();
  });

  it('shows existing option values', () => {
    render(<QueryEditorCheckRuns gitRef="abc123" checkName="build" onChange={jest.fn()} />);
    expect(screen.getByDisplayValue('abc123')).toBeInTheDocument();
    expect(screen.getByDisplayValue('build')).toBeInTheDocument();
  });

  it('propagates the ref on blur', () => {
    const onChange = jest.fn();
    render(<QueryEditorCheckRuns onChange={onChange} />);

    const refInput = screen.getByLabelText('Query editor ref');
    fireEvent.change(refInput, { target: { value: 'heads/main' } });
    fireEvent.blur(refInput);

    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({ gitRef: 'heads/main' }));
  });
});
