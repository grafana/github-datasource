# E2E Mocks (WireMock)

This plugin uses [WireMock](https://wiremock.org/) to record and replay API responses for e2e tests.

## First-time setup

Run the `setup-e2e-mocks` Claude agent to customize the scaffolding for this plugin:

```
/agents/setup-e2e-mocks
```

The agent will detect your API hosts, customize proxy stubs and wire up docker-compose and provisioning env vars.

## Directory structure

```
wiremock/
  mappings/
    proxy-*.json      # proxy stubs - route requests to real APIs during recording
    *.json            # recorded stubs - saved request matchers (committed)
  __files/
    *                 # recorded response bodies (committed)
e2e/
  docker-compose.record.yaml   # WireMock in recording mode
  docker-compose.replay.yaml   # WireMock serving saved stubs
```

## npm scripts

| Script               | Description                                         |
| -------------------- | --------------------------------------------------- |
| `yarn e2e`           | Run Playwright tests (unchanged)                    |
| `yarn server:record` | Clean old recordings, start WireMock in record mode |
| `yarn server:replay` | Start WireMock serving recorded stubs               |

## Recording

```bash
yarn build && mage -v build:linux
yarn server:record
yarn e2e
docker compose -f e2e/docker-compose.record.yaml down
```

This cleans previous recordings (preserving `proxy-*.json`), starts WireMock in record mode and runs the existing e2e tests to capture API traffic as stubs.

**Important:** recording requires real API credentials set via env vars in `e2e/docker-compose.record.yaml`.

## Replaying (locally)

Replay recorded stubs without real credentials:

```bash
yarn server:replay
yarn e2e
docker compose -f e2e/docker-compose.replay.yaml down
```

In CI, the replay docker-compose file is passed to the workflow via `playwright-docker-compose-file`.

## Sanitizing recordings

Review recorded stubs before committing:

- Check `wiremock/mappings/` for request headers containing API keys or tokens
- Check `wiremock/__files/` for response bodies containing PII or credentials
- Remove recordings from endpoints you don't need (health checks, telemetry)
- WireMock records request headers as-is - review recorded stubs for Authorization or Cookie headers
