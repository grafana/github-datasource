import { defineConfig } from 'eslint/config';
import baseConfig from './.config/eslint.config.mjs';

export default defineConfig([
  {
    ignores: [
      '**/node_modules/',
      '**/coverage/',
      '**/ci/',
      '**/dist/',
      '**/.idea/',
      'test-results/',
      'playwright-report/',
      'blob-report/',
      'playwright/.cache/',
      'playwright/.auth/',
      '**/.eslintcache',
      '**/__debug_bin',
      '**/mage_output_file.go',
      'provisioning/',
      '**/.pnp.*',
      '.yarn/*',
      '!.yarn/patches',
      '!.yarn/plugins',
      '!.yarn/releases',
      '!.yarn/sdks',
      '!.yarn/versions',
    ],
  },
  ...baseConfig,
]);
