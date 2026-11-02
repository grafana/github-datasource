import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryEditorPullRequests } from './QueryEditorPullRequests';

let multiComboboxProps: any;

jest.mock('@grafana/ui', () => {
  const actual = jest.requireActual('@grafana/ui');
  return {
    ...actual,
    MultiCombobox: (props: any) => {
      multiComboboxProps = props;
      return <div data-testid="fields-picker">{props.value.join(',')}</div>;
    },
  };
});

const optionsFor = (values: string[]) => values.map((value) => ({ value }));

describe('QueryEditorPullRequests', () => {
  it('summarizes all fields when fields is unset', () => {
    render(<QueryEditorPullRequests onChange={jest.fn()} />);
    expect(screen.getByText('Options')).toBeInTheDocument();
    expect(screen.getByText('Additional fields: All')).toBeInTheDocument();
  });

  it('summarizes the selected fields', () => {
    render(<QueryEditorPullRequests fields={['labels', 'author']} onChange={jest.fn()} />);
    expect(screen.getByText('Additional fields: Author, Labels')).toBeInTheDocument();
  });

  it('does not summarize a list with duplicates as all fields', () => {
    render(
      <QueryEditorPullRequests fields={['labels', 'labels', 'author', 'author', 'diffstat']} onChange={jest.fn()} />
    );
    expect(screen.getByText('Additional fields: Additions / deletions, Author, Labels')).toBeInTheDocument();
  });

  it('summarizes no additional fields', () => {
    render(<QueryEditorPullRequests fields={[]} onChange={jest.fn()} />);
    expect(screen.getByText('Additional fields: None')).toBeInTheDocument();
  });

  it('selects every field in the picker when fields is unset', async () => {
    render(<QueryEditorPullRequests onChange={jest.fn()} />);
    await userEvent.click(screen.getByText('Options'));
    expect(screen.getByTestId('fields-picker')).toHaveTextContent('diffstat,mergeable,author,merged_by,labels');
  });

  it('stores the selected fields, or leaves fields unset when every field is selected', async () => {
    const onChange = jest.fn();
    render(<QueryEditorPullRequests query="is:open" onChange={onChange} />);
    await userEvent.click(screen.getByText('Options'));

    multiComboboxProps.onChange(optionsFor(['author', 'labels']));
    expect(onChange).toHaveBeenLastCalledWith(
      expect.objectContaining({ query: 'is:open', fields: ['author', 'labels'] })
    );

    multiComboboxProps.onChange(optionsFor([]));
    expect(onChange).toHaveBeenLastCalledWith(expect.objectContaining({ fields: [] }));

    multiComboboxProps.onChange(optionsFor(['diffstat', 'mergeable', 'author', 'merged_by', 'labels']));
    expect(onChange).toHaveBeenLastCalledWith(expect.objectContaining({ fields: undefined }));
  });
});
