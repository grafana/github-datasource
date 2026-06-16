import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConfigEditor from './ConfigEditor';

describe('Config Editor', () => {
  it('should select basic license type as checked by default', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: {}, secureJsonFields: {} } as any;
    render(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Connection')).toBeInTheDocument());
    onOptionsChange.mockClear(); // this is called on component render, so we need to clear it to avoid false positives
    expect(screen.getByLabelText('Free, Pro & Team')).toBeChecked();
    expect(screen.getByLabelText('Enterprise Server')).not.toBeChecked();
    expect(screen.queryByText('GitHub Enterprise Server URL')).not.toBeInTheDocument();
    expect(onOptionsChange).toHaveBeenCalledTimes(0);
    await userEvent.click(screen.getByLabelText('Enterprise Server'));
    expect(onOptionsChange).toHaveBeenCalledTimes(1);
    expect(screen.queryByText('GitHub Enterprise Server URL')).toBeInTheDocument();
  });
  it('should select enterprise license type as checked when the url is not empty', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: { githubUrl: 'https://foo.bar' }, secureJsonFields: {} } as any;
    render(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Connection')).toBeInTheDocument());
    onOptionsChange.mockClear();
    expect(screen.getByLabelText('Free, Pro & Team')).not.toBeChecked();
    expect(screen.getByLabelText('Enterprise Server')).toBeChecked();
    expect(screen.queryByText('GitHub Enterprise Server URL')).toBeInTheDocument();
    expect(onOptionsChange).toHaveBeenCalledTimes(0);
    await userEvent.click(screen.getByLabelText('Free, Pro & Team'));
    expect(onOptionsChange).toHaveBeenNthCalledWith(1, {
      jsonData: { githubUrl: '', githubPlan: 'github-basic' },
      secureJsonFields: {},
    });
    expect(screen.queryByText('GitHub Enterprise URL')).not.toBeInTheDocument();
    await userEvent.click(screen.getByLabelText('Enterprise Cloud'));
    expect(onOptionsChange).toHaveBeenNthCalledWith(2, {
      jsonData: { githubUrl: '', githubPlan: 'github-enterprise-cloud' },
      secureJsonFields: {},
    });
    await userEvent.click(screen.getByLabelText('Enterprise Server'));
    expect(screen.queryByText('GitHub Enterprise Server URL')).toBeInTheDocument();
  });
});
