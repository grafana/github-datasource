import { z } from 'zod';
import * as fs from 'fs';
import * as path from 'path';

import { GitHubDataSourceOptionsSchema, GitHubSecureJsonDataSchema } from '../src/types/config';

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

// Write to src/schemas/ for frontend consumers
const frontendPath = path.resolve(__dirname, '..', 'src', 'schemas', 'config-schema.json');
fs.mkdirSync(path.dirname(frontendPath), { recursive: true });
fs.writeFileSync(frontendPath, schemaJSON);
console.log(`Config JSON schema written to ${frontendPath}`);

// Write to pkg/configschema/ for Go backend embedding
const backendPath = path.resolve(__dirname, '..', 'pkg', 'configschema', 'config_schema.json');
fs.mkdirSync(path.dirname(backendPath), { recursive: true });
fs.writeFileSync(backendPath, schemaJSON);
console.log(`Config JSON schema written to ${backendPath}`);
