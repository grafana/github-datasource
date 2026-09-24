# Contributing

## Signed commits are required

> [!IMPORTANT]
> All commits must be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) (GPG, SSH, or S/MIME) to be merged into this repository. Pull requests with unsigned commits will need to be re-committed with signatures before they can be merged.

Thank you for your interest in contributing to this repository. We are glad you want help us to improve the project and join our community. Feel free to [browse the open issues](https://github.com/grafana/github-datasource/issues). If you wanna more straightforward tasks to complete, [we have some](https://github.com/grafana/github-datasource/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22). For more details about how you can help, please take a look at [Grafana's Contributing Guide](https://github.com/grafana/grafana/blob/main/CONTRIBUTING.md).

## Next steps

- [Set up your development environment](./docs/developer-guide.md)

## Data source configuration schema

[`pkg/schema/dsconfig.json`](./pkg/schema/dsconfig.json) describes the plugin's configuration fields, storage locations, validation rules, UI hints, and instructions for automation. Edit this file when changing the configuration schema.

The format and shared tooling come from [`grafana/dsconfig`](https://github.com/grafana/dsconfig/tree/main/dsconfig):

- [README](https://github.com/grafana/dsconfig/tree/main/dsconfig#readme): concepts and examples.
- [Schema reference](https://github.com/grafana/dsconfig/blob/main/dsconfig/schema.md): field properties and validation rules.
- [Publishing guide](https://github.com/grafana/dsconfig/blob/main/dsconfig/PUBLISH-SCHEMA.md): generation and packaging setup.

### Source and generated files

| File | Purpose |
| --- | --- |
| `pkg/schema/dsconfig.json` | Hand-authored configuration schema. |
| `pkg/models/settings.go` | Backend settings model, parsing, and `SecureJsonDataKeys`. |
| `pkg/schema/dsconfig_test.go` | Shared conformance test setup and provisioning examples. |
| `pkg/schema/schema.gen.json` | Generated plugin schema. |
| `pkg/schema/settings.gen.json` | Generated settings schema. |
| `pkg/schema/settings.examples.gen.json` | Generated provisioning examples. |

Do not edit `.gen.json` files by hand. The generator uses the version of `github.com/grafana/dsconfig/schema` pinned in `go.mod`. The `$schema` URL in `dsconfig.json` identifies the JSON Schema used to validate the authoring format; review that pin when upgrading dsconfig.

### Updating configuration

1. Update `pkg/schema/dsconfig.json`, including any relevant groups, validation rules, and instructions. Keep stored fields consistent with the configuration editor and backend. A virtual field represents a UI control rather than a saved setting.
2. Update `models.Settings` and its JSON tags where needed. For secrets, update `models.SecureJsonDataKeys` and the code that reads `DecryptedSecureJSONData` instead of adding a normal JSON setting. Preserve support for existing configurations when changing types, including string and numeric GitHub App IDs.
3. Update `SettingsExamples` in `pkg/schema/dsconfig_test.go` when the configuration examples change. The existing examples cover personal access tokens, GitHub Apps, Enterprise Cloud, and Enterprise Server. Use placeholders, never real credentials.
4. Regenerate the artifacts from the repository root:

   ```bash
   go generate ./pkg/schema/...
   ```

5. Verify schema consistency and backend parsing:

   ```bash
   go test ./pkg/schema/... ./pkg/models/... -count=1
   ```

6. Review and commit any generated changes alongside the source changes. An instructions-only edit may leave the generated artifacts unchanged.

The `go:generate` directive in `dsconfig_test.go` runs `go test -run TestPlugin -generateArtifacts`. Normal test runs check the artifacts without regenerating them. A `SchemaArtifactInSync` failure means regeneration is needed; key, type, or secret consistency failures require correcting the schema or Go declarations.

### Building the published files

Run `npm run build` to copy the source schema and generated artifacts into the plugin's `dist/schema/` directory via `webpack.config.ts`. This packages the files; it does not replace the Go generation step. Grafana serves the source schema at `/public/plugins/grafana-github-datasource/schema/dsconfig.json` once that build is installed.
