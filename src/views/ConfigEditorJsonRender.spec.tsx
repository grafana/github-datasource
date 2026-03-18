import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConfigEditorJsonRender from './ConfigEditorJsonRender';

describe('ConfigEditorJsonRender (json-render POC)', () => {
  it('should render authentication and connection sections', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: {}, secureJsonFields: {} } as any;
    render(<ConfigEditorJsonRender options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Authentication')).toBeInTheDocument());
    expect(screen.getByText('Connection')).toBeInTheDocument();
  });

  it('should default to PAT auth type and basic license', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: {}, secureJsonFields: {} } as any;
    render(<ConfigEditorJsonRender options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Authentication')).toBeInTheDocument());
    expect(screen.getByLabelText('Personal Access Token')).toBeChecked();
    expect(screen.getByLabelText('Free, Pro & Team')).toBeChecked();
  });

  it('should show enterprise URL field when enterprise server is selected', async () => {
    const onOptionsChange = jest.fn();
    const options = {
      jsonData: { githubPlan: 'github-enterprise-server', selectedAuthType: 'personal-access-token' },
      secureJsonFields: {},
    } as any;
    render(<ConfigEditorJsonRender options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Authentication')).toBeInTheDocument());
    expect(screen.getByLabelText('Enterprise Server')).toBeChecked();
    expect(
      screen.getByText('The URL of the GitHub Enterprise Server instance')
    ).toBeInTheDocument();
  });

  it('should show GitHub App fields when github-app auth is selected', async () => {
    const onOptionsChange = jest.fn();
    const options = {
      jsonData: { selectedAuthType: 'github-app' },
      secureJsonFields: {},
    } as any;
    render(<ConfigEditorJsonRender options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Authentication')).toBeInTheDocument());
    expect(screen.getByLabelText('GitHub App')).toBeChecked();
    expect(screen.getByText('The GitHub App ID')).toBeInTheDocument();
    expect(screen.getByText('The GitHub App installation ID')).toBeInTheDocument();
    expect(
      screen.getByText('Private key for GitHub App authentication (PEM format)')
    ).toBeInTheDocument();
  });

  it('should call onOptionsChange when switching auth type', async () => {
    const onOptionsChange = jest.fn();
    const options = {
      jsonData: { selectedAuthType: 'personal-access-token' },
      secureJsonFields: {},
    } as any;
    render(<ConfigEditorJsonRender options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Authentication')).toBeInTheDocument());
    onOptionsChange.mockClear();
    await userEvent.click(screen.getByLabelText('GitHub App'));
    expect(onOptionsChange).toHaveBeenCalledTimes(1);
    expect(onOptionsChange).toHaveBeenCalledWith(
      expect.objectContaining({
        jsonData: expect.objectContaining({ selectedAuthType: 'github-app' }),
      })
    );
  });
});
