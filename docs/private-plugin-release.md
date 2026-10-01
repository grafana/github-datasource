# Release and deploy the private GitHub data source

This tutorial covers the manual setup and release steps for the `rwe-github-datasource` private Grafana plugin. The release workflow builds, tests, signs, packages, and publishes the plugin. It signs one immutable artifact for both sphere Grafana instances.

## Values used by this repository

| Item                          | Value                                           |
| ----------------------------- | ----------------------------------------------- |
| Grafana organization          | `https://grafana.com/orgs/rwe`                  |
| Plugin ID                     | `rwe-github-datasource`                         |
| GitHub environment            | `plugin-signing`                                |
| GitHub secret                 | `GRAFANA_ACCESS_POLICY_TOKEN`                   |
| Test Grafana `root_url`       | `https://grafana.sre.test.aws-eu1.energy.local` |
| Production Grafana `root_url` | `https://grafana.sre.aws-eu1.energy.local`      |
| Workflow                      | **Private plugin release**                      |

The signing URLs in `.github/workflows/publish.yaml` must exactly match the `root_url` values configured in sphere. A hostname, scheme, port, or path difference makes the private signature invalid for that Grafana instance.

## 1. Verify the Grafana organization slug

Private plugin signing binds the plugin ID to a Grafana Cloud organization. The first segment of the plugin ID is `rwe`, matching the confirmed organization at `https://grafana.com/orgs/rwe`.

1. Sign in to the Grafana Cloud portal.
2. Open **My Account** and select the organization that will own the plugin.
3. Confirm its organization slug is `rwe`.
4. Stop if the URL does not end in `/orgs/rwe`. A token from another organization cannot sign `rwe-github-datasource`.

A Grafana Cloud account is required to obtain the signing token even when the target Grafana instances are self-hosted.

## 2. Create the Grafana signing token

Grafana's current signing flow uses a Cloud Access Policy token. Do not look for a
legacy API key or create a Grafana service account token; neither credential signs
plugins.

1. Sign in at [grafana.com](https://grafana.com/) and click **My Account** to open
   the Cloud Portal. Do not open **Administration > Cloud access policies** inside
   a Grafana stack because that interface creates stack-realm policies only.
2. Select the `rwe` organization in the organization selector at the top of the
   Cloud Portal.
3. In the left navigation, select **Security > Access Policies**.
4. Click **Create access policy** and enter `github-datasource-signing` as its
   display name.
5. Set **Realm** to **rwe (all stacks)**. This organization realm is required for
   plugin signing.
6. At the bottom of the scopes section, click **Add scope**. Search for `plugins`
   and select **Plugins: Create and edit plugins on Grafana.com**
   (`plugins:write`). Grant no other scope.
7. Create the policy, open it, and click **Add token**.
8. Enter `github-datasource-release` as the token display name and set an
   expiration according to the credential-rotation policy.
9. Click **Create**, then **Copy to clipboard**. Grafana displays the token value
   only once; store it in the approved secret manager until it is added to GitHub.

Do not commit the token, put it in a shell history, add it to repository variables, or paste it into an issue.

### If `plugins:write` is unavailable

Do not substitute any of these scopes:

- `provisioned-plugins:write` provisions plugins onto a Grafana Cloud stack.
- `stack-plugins:write` manages catalog plugins installed in a Grafana Cloud
  stack.
- A Grafana service account token authorizes the Grafana HTTP API.

None of them authorizes the plugin signing endpoint. Check these conditions:

1. The browser is in the `grafana.com` Cloud Portal, not the hosted Grafana stack.
2. The selected organization is `rwe`.
3. The realm says **rwe (all stacks)** rather than a stack name.
4. The user has the Cloud Portal **Admin** role. Grafana stack Admin and Cloud
   Portal Admin are independent roles.
5. **Add scope** was used; the dialog does not display every scope by default.

If `plugins:write` is still absent, stop. Open **Support > Open a Ticket** in the
Cloud Portal and send:

```text
Subject: plugins:write scope unavailable for private plugin signing

Organization: rwe (https://grafana.com/orgs/rwe)
Plugin ID: rwe-github-datasource

We need to sign a private Grafana plugin for self-hosted Grafana instances.
The organization-level Access Policy dialog does not offer the plugins:write
scope, although the Grafana plugin-signing documentation requires it. It only
shows provisioned-plugins and stack-plugins scopes. Please enable or restore
private plugin signing and the plugins:write scope for this organization.
```

Cloud API keys are deprecated, which is why the portal may not expose an API key
page. The Cloud Access Policy API is not a bootstrap workaround: it also requires
an existing authorized token. Wait for Grafana Support rather than granting a
different plugin scope.

## 3. Configure the protected GitHub environment

In `rwe-supply-and-trading/github-datasource`:

1. Open **Settings > Environments**.
2. Create an environment named `plugin-signing`.
3. Add the platform maintainers as required reviewers.
4. Restrict deployment branches to `main`.
5. Add an environment secret named `GRAFANA_ACCESS_POLICY_TOKEN`.
6. Paste the Grafana Access Policy token as the secret value.
7. Save the environment.

The workflow receives the token only after environment approval. The two Grafana root URLs are non-secret configuration in the workflow.

## 4. Choose an artifact download location

Grafana must download the ZIP without interactive authentication when the pod starts.

### Public GitHub release

Use the GitHub release URL when the repository and release asset are publicly downloadable:

```text
https://github.com/rwe-supply-and-trading/github-datasource/releases/download/v<VERSION>/rwe-github-datasource-<VERSION>.zip
```

### Internal artifact repository

Use an internal JFrog generic repository when the GitHub repository is private or anonymous GitHub downloads are not acceptable. After the workflow creates the release:

1. Download the ZIP and `.sha256` files from the GitHub release.
2. Verify the checksum.
3. Upload both files to an immutable, versioned path in JFrog.
4. Confirm the Grafana pod can download the ZIP without embedding credentials in the URL.

Do not overwrite a released ZIP. Publish a new plugin version instead.

## 5. Prepare a release version

Use a fork-specific semantic version. For the next release, run locally with the Node version from `.nvmrc`:

```bash
nvm use
npm version 2.9.2-rwe.2 --no-git-tag-version
```

Then:

1. Add a matching entry to `CHANGELOG.md`.
2. Commit `package.json`, `package-lock.json`, and the changelog.
3. Open and merge a pull request to `main`.
4. Confirm the normal **Plugins - CI** workflow succeeds.

The release workflow rejects development versions containing `-dev`. It creates the Git tag itself, so do not create the release tag manually.

## 6. Run the signed release

1. Open **Actions** in `rwe-supply-and-trading/github-datasource`.
2. Select **Private plugin release**.
3. Select the `main` branch.
4. Click **Run workflow**.
5. Approve the `plugin-signing` environment deployment when prompted.
6. Wait for all build, test, signing, and packaging steps to finish.

The workflow performs these operations in order:

1. Installs locked Node dependencies.
2. Runs Jest, TypeScript, ESLint, and all Go tests.
3. Builds frontend assets and Linux AMD64/ARM64 backend binaries.
4. Signs the complete `dist` directory for both sphere root URLs.
5. Verifies `MANIFEST.txt` and executable permissions.
6. Creates a ZIP with a top-level `rwe-github-datasource/` directory.
7. Publishes the ZIP and SHA-256 checksum in a GitHub release tagged with the package version.

Signing must remain after every build step. Any modification to a signed file causes Grafana to report a modified signature.

## 7. Verify the release artifact

Download both release assets and verify them before deployment:

```bash
sha256sum --check rwe-github-datasource-<VERSION>.zip.sha256
zipinfo -1 rwe-github-datasource-<VERSION>.zip | head
unzip -p rwe-github-datasource-<VERSION>.zip \
  rwe-github-datasource/MANIFEST.txt | head
```

Confirm that:

- The checksum passes.
- The archive has the top-level `rwe-github-datasource/` directory.
- `MANIFEST.txt` exists inside that directory.
- `plugin.json` identifies `rwe-github-datasource` at the expected version.
- Linux AMD64 and ARM64 backend binaries are present and executable.

## 8. Add the fork to sphere test

Use a two-phase migration. Install the fork alongside `grafana-github-datasource` first; do not remove the upstream plugin until datasource and dashboard migration is complete.

In `gitops-bootstrap/sre-ops-k8s-sphere/components/grafana/base/grafana.yaml`, replace the legacy plugin-install environment variable with `GF_PLUGINS_PREINSTALL_SYNC`. Keep all existing plugins and add the signed fork using its immutable ZIP URL:

```yaml
- name: GF_PLUGINS_PREINSTALL_SYNC
  value: 'aws-datasource-provisioner-app@1.13.11,chaosmeshorg-datasource@3.0.0,grafana-github-datasource@2.9.1,rwe-github-datasource@<VERSION>@<PLUGIN_ZIP_URL>,grafana-advisor-app@1.0.2,yesoreyeram-infinity-datasource@4.0.0,marcusolsson-treemap-panel@2.1.1,grafana-lokiexplore-app@2.5.2,grafana-metricsdrilldown-app@2.5.1,grafana-exploretraces-app@2.2.0,grafana-pyroscope-app@2.3.0'
```

Remove `GF_INSTALL_PLUGINS` and `GF_INSTALL_PLUGINS_FORCE` in the same change. Grafana 13 supports the synchronous preinstall format:

```text
<plugin ID>@<plugin version>@<URL to plugin ZIP>
```

Open a GitOps pull request, render the test overlay, and sync only the sphere test Grafana application first.

## 9. Verify the plugin in sphere test

After Grafana restarts, locate its pod:

```bash
kubectl -n grafana get pods -l app.kubernetes.io/name=grafana
```

Check installation and signature messages:

```bash
kubectl -n grafana logs <GRAFANA_POD> | grep -Ei 'rwe-github-datasource|signature|plugin'
```

Confirm the signed files exist:

```bash
kubectl -n grafana exec <GRAFANA_POD> -- \
  test -s /var/lib/grafana/plugins/rwe-github-datasource/MANIFEST.txt
```

In the Grafana UI:

1. Go to **Administration > Plugins and data > Plugins**.
2. Open **GitHub (RWE)**.
3. Confirm the signature status is **Signed**.
4. Treat **Unsigned**, **Invalid signature**, or **Modified signature** as a failed deployment.

## 10. Configure and test the new datasource

The new plugin ID creates a separate datasource type. Existing `grafana-github-datasource` instances do not migrate automatically.

1. Go to **Connections > Data sources > Add data source**.
2. Select **GitHub (RWE)**.
3. Select **Enterprise Cloud with data residency**.
4. Enter `https://api.rwe.ghe.com` as **GitHub API URL**.
5. Select **GitHub App** authentication.
6. Enter the App ID and Installation ID.
7. Enter the private key directly in Grafana or through the approved secret-provisioning path.
8. Click **Save & test**.
9. Test at least one REST-backed query and one GraphQL-backed query.

Do not paste the GitHub App private key into GitHub issues, workflow inputs, pull requests, or documentation.

## 11. Migrate dashboards and remove upstream

1. Inventory dashboards, alerts, annotations, and variables bound to the old datasource UID.
2. Rebind them to the new `rwe-github-datasource` instance.
3. Verify representative dashboards in sphere test.
4. Promote the same signed ZIP and GitOps change to production.
5. Repeat datasource configuration and dashboard verification in production.
6. Remove `grafana-github-datasource@2.9.1` from `GF_PLUGINS_PREINSTALL_SYNC` only after no references remain.
7. Remove the old datasource instance after a final rollback window.

Do not rebuild or re-sign between test and production. Promote the exact artifact whose checksum was verified in test.

## 12. Rotate the signing token

Before the Access Policy token expires:

1. Create a replacement token on the same `plugins:write` policy.
2. Replace the `GRAFANA_ACCESS_POLICY_TOKEN` environment secret.
3. Run the next normal release.
4. Revoke the old token.

Token rotation does not invalidate already signed plugin releases. A `root_url` change does require a new signed release containing the new exact URL.

## Troubleshooting

### Signing says the organization does not own the plugin ID

The Grafana Cloud organization slug does not match the `rwe` plugin ID prefix. Use the confirmed organization at `https://grafana.com/orgs/rwe` or rename the plugin before release.

### Signing reports that `rootUrls` is required

Confirm the workflow passes the comma-separated `ROOT_URLS` value and that the signing token belongs to the matching organization.

### Grafana reports an invalid private signature

Compare the deployed Grafana `root_url` with the URLs in `.github/workflows/publish.yaml`. They must match exactly.

### Grafana reports a modified signature

The ZIP contents differ from the files that were signed. Rebuild, sign after all builds finish, package without changing `dist`, and publish a new version.

### Grafana cannot download the ZIP

Test the URL from the Grafana pod network. For private GitHub releases, publish the artifact to an approved internal repository that Grafana can access without interactive GitHub authentication.
