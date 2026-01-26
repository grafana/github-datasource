import React from 'react';
import QueryEditorRepository from './QueryEditorRepository';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { components } from 'components/selectors';

describe('QueryEditorRepository', () => {
  it('should render with initial props', () => {
    const props = {
      owner: 'grafana',
      repository: 'grafana',
      onChange: jest.fn(),
    };
    render(<QueryEditorRepository {...props} />);

    expect(screen.getByLabelText(components.QueryEditor.Owner.input)).toHaveValue('grafana');
    expect(screen.getByLabelText(components.QueryEditor.Repository.input)).toHaveValue('grafana');
  });

  it('should allow typing without calling onChange until blur', async () => {
    const onChange = jest.fn();
    const props = {
      owner: '',
      repository: '',
      onChange,
    };
    render(<QueryEditorRepository {...props} />);

    const ownerInput = screen.getByLabelText(components.QueryEditor.Owner.input);
    await userEvent.type(ownerInput, 'grafana');

    // onChange should not be called while typing
    expect(onChange).not.toHaveBeenCalled();
    // Local state should reflect the typed value
    expect(ownerInput).toHaveValue('grafana');
  });

  it('should call onChange on blur with the typed value', async () => {
    const onChange = jest.fn();
    const props = {
      owner: '',
      repository: '',
      onChange,
    };
    render(<QueryEditorRepository {...props} />);

    const ownerInput = screen.getByLabelText(components.QueryEditor.Owner.input);
    await userEvent.type(ownerInput, 'grafana');
    await userEvent.tab(); // blur

    expect(onChange).toHaveBeenCalledWith(
      expect.objectContaining({
        owner: 'grafana',
      })
    );
  });

  it('should update local state when props change externally', () => {
    const onChange = jest.fn();
    const props = {
      owner: 'initial-owner',
      repository: 'initial-repo',
      onChange,
    };
    const { rerender } = render(<QueryEditorRepository {...props} />);

    expect(screen.getByLabelText(components.QueryEditor.Owner.input)).toHaveValue('initial-owner');

    // Simulate parent updating the props externally
    rerender(<QueryEditorRepository {...props} owner="new-owner" repository="new-repo" />);

    expect(screen.getByLabelText(components.QueryEditor.Owner.input)).toHaveValue('new-owner');
    expect(screen.getByLabelText(components.QueryEditor.Repository.input)).toHaveValue('new-repo');
  });

  it('should NOT reset user input while typing (before blur)', async () => {
    const onChange = jest.fn();
    const props = {
      owner: '',
      repository: '',
      onChange,
    };
    const { rerender } = render(<QueryEditorRepository {...props} />);

    const ownerInput = screen.getByLabelText(components.QueryEditor.Owner.input);
    await userEvent.type(ownerInput, 'grafana');

    // Rerender with same props (simulating parent re-render without prop change)
    rerender(<QueryEditorRepository {...props} />);

    // User's typed value should be preserved, not reset to empty string
    expect(ownerInput).toHaveValue('grafana');
  });
});
