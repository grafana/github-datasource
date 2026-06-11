import React from 'react';
import { QueryEditorBranches } from './QueryEditorBranches';
import { render, screen } from '@testing-library/react';

describe('QueryEditorBranches', () => {
  it('renders a Filter input field', () => {
    const props = { onChange: jest.fn() };
    render(<QueryEditorBranches {...props} />);
    expect(screen.getByPlaceholderText('release/')).toBeInTheDocument();
  });

  it('shows existing query value', () => {
    const onChange = jest.fn();
    render(<QueryEditorBranches query="release/" onChange={onChange} />);
    const input = screen.getByPlaceholderText('release/');
    expect(input).toHaveValue('release/');
  });
});
