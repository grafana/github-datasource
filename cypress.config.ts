import conf from '@grafana/e2e/cypress.config.js';
import { defineConfig } from 'cypress';

export default defineConfig({
  ...conf,
  video: false,
  videoCompression: false,
  screenshotOnRunFailure: true,
  reporter: 'spec',
  retries: {
    runMode: 3,
  },

  e2e: {
    ...conf.e2e,
    specPattern: 'cypress/integration/**/*.spec.ts',
  },
});
