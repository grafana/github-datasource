import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryEditorSecurityAdvisories } from './QueryEditorSecurityAdvisories';
import { QueryTypes } from '../constants';
import { isValid } from '../validation';

// jsdom cannot lay out Grafana's virtualized dropdown.
jest.mock('@grafana/ui', () => ({
  ...jest.requireActual('@grafana/ui'),
  Combobox: ({ id, options, value, onChange }: {
    id: string;
    options: Array<{ label: string; value: string }>;
    value: string;
    onChange: (option: { value: string }) => void;
  }) => (
    <select id={id} value={value} onChange={(event) => onChange({ value: event.currentTarget.value })}>
      {options.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
    </select>
  ),
}));

it('offers the query type and requires only an owner', () => {
  expect(QueryTypes).toContain('Security_Advisories');
  expect(isValid({ refId: 'A', queryType: 'Security_Advisories', owner: 'grafana' })).toBe(true);
  expect(isValid({ refId: 'A', queryType: 'Security_Advisories' })).toBe(false);
});

it('selects a state and clears it with All', () => {
  const onChange = jest.fn();
  const { rerender } = render(<QueryEditorSecurityAdvisories onChange={onChange} />);
  const select = screen.getByRole('combobox', { name: 'State' });
  expect(select).toHaveValue('all');
  fireEvent.change(select, { target: { value: 'draft' } });
  expect(onChange).toHaveBeenLastCalledWith({ state: 'draft' });

  rerender(<QueryEditorSecurityAdvisories state="draft" onChange={onChange} />);
  expect(select).toHaveValue('draft');
  fireEvent.change(select, { target: { value: 'all' } });
  expect(onChange).toHaveBeenLastCalledWith({ state: undefined });
});
