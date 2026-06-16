import React from 'react';
import { QueryEditorIssues } from './QueryEditorIssues';
import { render, screen } from '@testing-library/react';

const comboboxOnChangeSpy = jest.fn();

jest.mock('@grafana/ui', () => {
  const actual = jest.requireActual('@grafana/ui');
  return {
    ...actual,
    Combobox: (props: any) => {
      comboboxOnChangeSpy.mockImplementation(props.onChange);
      return (
        <select
          data-testid={props['data-testid']}
          value={props.value}
          onChange={(e) => {
            const opt = props.options.find((o: any) => String(o.value) === e.target.value);
            if (opt) {
              props.onChange(opt);
            }
          }}
        >
          {props.options?.map((o: any) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      );
    },
  };
});

describe('QueryEditorIssues', () => {
  it('should have CreatedAt, ClosedAt and UpdatedAt time field option', () => {
    const props = {
      onChange: jest.fn(),
    };
    render(<QueryEditorIssues {...props} />);
    expect(screen.getByText('Time Field')).toBeInTheDocument();
    expect(screen.getByText('CreatedAt')).toBeInTheDocument();
    expect(screen.getByText('ClosedAt')).toBeInTheDocument();
    expect(screen.getByText('UpdatedAt')).toBeInTheDocument();
  });
});
