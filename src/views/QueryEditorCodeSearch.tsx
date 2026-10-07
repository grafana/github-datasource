import React, { useState } from 'react';
import { Checkbox, Input } from '@grafana/ui';
import { EditorField, EditorRow } from '@grafana/plugin-ui';
import { CodeSearchOptions } from '../types/query';
import { LeftColumnWidth, RightColumnWidth } from './QueryEditor';

interface Props extends CodeSearchOptions {
  onChange: (value: CodeSearchOptions) => void;
}

export const QueryEditorCodeSearch = ({ query = '', exactPath = '', includeTextMatches = false, requireComplete = false, onChange }: Props) => {
  const [searchQuery, setSearchQuery] = useState(query);
  const [path, setPath] = useState(exactPath);
  const [textMatches, setTextMatches] = useState(includeTextMatches);
	const [complete, setComplete] = useState(requireComplete);

  const update = (nextTextMatches = textMatches, nextComplete = complete) =>
    onChange({ query: searchQuery, exactPath: path, includeTextMatches: nextTextMatches, requireComplete: nextComplete });

  return (
    <EditorRow>
      <EditorField label="Search query" tooltip="GitHub code search syntax, including org, repo, path, filename, and content terms." width={LeftColumnWidth}>
        <Input aria-label="Code search query" value={searchQuery} onChange={(event) => setSearchQuery(event.currentTarget.value)} onBlur={() => update()} width={RightColumnWidth} />
      </EditorField>
      <EditorField label="Exact path" tooltip="Optional exact-path filter applied after GitHub search." width={LeftColumnWidth}>
        <Input aria-label="Code search exact path" value={path} onChange={(event) => setPath(event.currentTarget.value)} onBlur={() => update()} width={RightColumnWidth} />
      </EditorField>
      <EditorField label="Text matches" tooltip="Request matching content fragments from GitHub." width={LeftColumnWidth}>
        <Checkbox
          aria-label="Include code search text matches"
          value={textMatches}
          onChange={(event) => {
            const checked = event.currentTarget.checked;
            setTextMatches(checked);
            update(checked);
          }}
        />
      </EditorField>
      <EditorField label="Require complete" tooltip="Fail the query rather than return partial search results." width={LeftColumnWidth}>
        <Checkbox
          aria-label="Require complete code search results"
          value={complete}
          onChange={(event) => {
            const checked = event.currentTarget.checked;
            setComplete(checked);
            update(textMatches, checked);
          }}
        />
      </EditorField>
    </EditorRow>
  );
};
