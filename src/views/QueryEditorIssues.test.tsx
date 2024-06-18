import React from 'react';
import QueryEditorIssues from './QueryEditorIssues';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { components } from 'components/selectors';

describe('QueryEditorIssues', () => {
  it('should have CreatedAt, ClosedAt and UpdatedAt time field option', async () => {
    const props = {
      onChange: jest.fn(),
    };
    render(<QueryEditorIssues {...props} />);
    expect(screen.getByText('Time Field')).toBeInTheDocument();
    userEvent.click(screen.getByTestId(components.QueryEditor.Issues.timeFieldInput));
    expect(await screen.findByText('CreatedAt')).toBeInTheDocument();
    expect(await screen.findByText('ClosedAt')).toBeInTheDocument();
    expect(await screen.findByText('UpdatedAt')).toBeInTheDocument();
  });
});
