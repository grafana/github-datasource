import { E2ESelectors } from '@grafana/e2e-selectors';

export const Components = {
  ConfigEditor: {
    AccessToken: {
      input: 'Config Editor Access Token',
    },
  },
  QueryEditor: {
    Owner: {
      input: 'Query Editor Owner',
    },
    Repository: {
      input: 'Query Editor Repository',
    },
    Query: {
      input: 'Query Editor Query',
    },
    Ref: {
      input: 'Query Editor Ref',
    },
  },
};

export const selectors: { components: E2ESelectors<typeof Components> } = {
  components: Components,
};
