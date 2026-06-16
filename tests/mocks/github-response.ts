export const githubResponse = {
  results: {
    A: {
      status: 200,
      frames: [
        {
          schema: {
            name: 'releases',
            refId: 'A',
            fields: [
              {
                name: 'name',
                type: 'string',
                typeInfo: {
                  frame: 'string',
                },
              },
              {
                name: 'created_by',
                type: 'string',
                typeInfo: {
                  frame: 'string',
                },
              },
              {
                name: 'is_draft',
                type: 'boolean',
                typeInfo: {
                  frame: 'bool',
                },
              },
              {
                name: 'is_prerelease',
                type: 'boolean',
                typeInfo: {
                  frame: 'bool',
                },
              },
              {
                name: 'tag',
                type: 'string',
                typeInfo: {
                  frame: 'string',
                },
              },
              {
                name: 'url',
                type: 'string',
                typeInfo: {
                  frame: 'string',
                },
              },
              {
                name: 'created_at',
                type: 'time',
                typeInfo: {
                  frame: 'time.Time',
                },
              },
              {
                name: 'published_at',
                type: 'time',
                typeInfo: {
                  frame: 'time.Time',
                  nullable: true,
                },
              },
            ],
          },
          data: {
            values: [
              ['grafana-github-datasource v1.5.7'],
              ['grafanabot'],
              [false],
              [false],
              ['v1.5.7'],
              ['https://github.com/grafana/github-datasource/releases/tag/v1.5.7'],
              [1713519461000],
              [1713519486000],
            ],
          },
        },
      ],
    },
  },
};
