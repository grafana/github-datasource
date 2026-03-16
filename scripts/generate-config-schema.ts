import { z } from 'zod';
import * as fs from 'fs';
import * as path from 'path';

import { GitHubDataSourceOptionsSchema, GitHubSecureJsonDataSchema } from '../src/types/config';

// Guard: ensure config.ts only imports from external packages (no local file imports)
const configSource = fs.readFileSync(path.resolve(__dirname, '..', 'src', 'types', 'config.ts'), 'utf-8');
const localImportPattern = /(?:^|\n)\s*import\s.*from\s+['"]\.\.?\//;
if (localImportPattern.test(configSource)) {
  console.error('Error: src/types/config.ts must not import from local files.');
  console.error('All config types must be self-contained in a single file to ensure schema generation is reliable.');
  process.exit(1);
}

const configSchema = {
  $schema: 'https://json-schema.org/draft/2020-12/schema',
  title: 'GitHubDataSourceConfig',
  description: 'Configuration schema for the Grafana GitHub data source plugin',
  type: 'object' as const,
  properties: {
    jsonData: z.toJSONSchema(GitHubDataSourceOptionsSchema),
    secureJsonData: z.toJSONSchema(GitHubSecureJsonDataSchema),
  },
};

const schemaJSON = JSON.stringify(configSchema, null, 2) + '\n';

// Write to src/schema/ for frontend consumers
const frontendPath = path.resolve(__dirname, '..', 'src', 'schema', 'config.json');
fs.mkdirSync(path.dirname(frontendPath), { recursive: true });
fs.writeFileSync(frontendPath, schemaJSON);
console.log(`Config JSON schema written to ${frontendPath}`);

// Write to pkg/configschema/ for Go backend embedding
const backendPath = path.resolve(__dirname, '..', 'pkg', 'configschema', 'config.json');
fs.mkdirSync(path.dirname(backendPath), { recursive: true });
fs.writeFileSync(backendPath, schemaJSON);
console.log(`Config JSON schema written to ${backendPath}`);
