import { E2ESelectors } from '@grafana/e2e-selectors';

export const Components = {
  ConfigEditor: {
    AccessToken: {
      input: 'Config editor access token',
    },
  },
  QueryEditor: {
    Owner: {
      input: 'Query editor owner',
    },
    Repository: {
      input: 'Query editor repository',
    },
    Query: {
      input: 'Query editor query',
    },
    Ref: {
      input: 'Query editor ref',
    },
  },
  AnnotationEditor: {
    container: 'Annotation editor container',
  },
};

export const selectors: { components: E2ESelectors<typeof Components> } = {
  components: Components,
};
