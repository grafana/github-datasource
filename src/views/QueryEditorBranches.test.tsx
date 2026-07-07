import React from 'react';
import { QueryEditorBranches } from './QueryEditorBranches';
import { render, screen } from '@testing-library/react';

describe('QueryEditorBranches', () => {
  it('renders a Filter input field', () => {
    render(<QueryEditorBranches onChange={jest.fn()} />);
    expect(screen.getByText('Filter')).toBeInTheDocument();
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('shows existing query value', () => {
    render(<QueryEditorBranches query="release/" onChange={jest.fn()} />);
    expect(screen.getByRole('textbox')).toHaveValue('release/');
  });
});
