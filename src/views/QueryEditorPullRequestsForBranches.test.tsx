import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { QueryEditorPullRequestsForBranches } from './QueryEditorPullRequestsForBranches';

describe('QueryEditorPullRequestsForBranches', () => {
  it('shows existing branches', () => {
    render(<QueryEditorPullRequestsForBranches branches={['feat/a', 'fix/b']} onChange={jest.fn()} />);
    expect(screen.getByRole('textbox')).toHaveValue('feat/a, fix/b');
  });

  it('splits the input into trimmed branch names on blur', () => {
    const onChange = jest.fn();
    render(<QueryEditorPullRequestsForBranches onChange={onChange} />);
    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: ' feat/a ,, fix/b,' } });
    fireEvent.blur(input);
    expect(onChange).toHaveBeenCalledWith({ branches: ['feat/a', 'fix/b'] });
  });
});
