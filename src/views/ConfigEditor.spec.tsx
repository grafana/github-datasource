import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConfigEditor from './ConfigEditor';

describe('Config Editor', () => {
  it('should select basic license type as checked by default', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: {}, secureJsonFields: {} } as any;
    render(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Additional Settings')).toBeInTheDocument());
    onOptionsChange.mockClear(); // this is called on component render, so we need to clear it to avoid false positives
    expect(screen.getByLabelText('Basic')).toBeChecked();
    expect(screen.getByLabelText('Enterprise')).not.toBeChecked();
    expect(screen.queryByText('GitHub Enterprise URL')).not.toBeInTheDocument();
    expect(onOptionsChange).toHaveBeenCalledTimes(0);
    await userEvent.click(screen.getByLabelText('Enterprise'));
    expect(onOptionsChange).toHaveBeenCalledTimes(0);
    expect(screen.queryByText('GitHub Enterprise URL')).toBeInTheDocument();
  });
  it('should select enterprise license type as checked when the url is not empty', async () => {
    const onOptionsChange = jest.fn();
    const options = { jsonData: { githubUrl: 'https://foo.bar' }, secureJsonFields: {} } as any;
    render(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
    await waitFor(() => expect(screen.getByText('Additional Settings')).toBeInTheDocument());
    onOptionsChange.mockClear();
    expect(screen.getByLabelText('Basic')).not.toBeChecked();
    expect(screen.getByLabelText('Enterprise')).toBeChecked();
    expect(screen.queryByText('GitHub Enterprise URL')).toBeInTheDocument();
    expect(onOptionsChange).toHaveBeenCalledTimes(0);
    await userEvent.click(screen.getByLabelText('Basic'));
    expect(onOptionsChange).toHaveBeenNthCalledWith(1, { jsonData: { githubUrl: '' }, secureJsonFields: {} });
    expect(screen.queryByText('GitHub Enterprise URL')).not.toBeInTheDocument();
  });
});
