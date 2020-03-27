# Google Sheets Datasource

⚠️ Work in Progress ⚠️ 

This aims to be the first production ready plugin using the new https://github.com/grafana/grafana-plugin-sdk-go and will become the best practice example for grafana 6.7+

In particular, the code style and tooling needs work before it should be widely replicated.

## Development

You need to install the following first:

* [Mage](https://magefile.org/)
* [Yarn](https://yarnpkg.com/)
* [Docker Compose](https://docs.docker.com/compose/)

```
mage watch
```

In another terminal
```
docker-compose up
```

To restart after backend changes:
`./scripts/restart-plugin.sh`
