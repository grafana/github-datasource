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

const outPath = path.resolve(__dirname, '..', 'pkg', 'schema', 'config.json');
fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, schemaJSON);
console.log(`Config JSON schema written to ${outPath}`);
