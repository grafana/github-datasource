# Grafana Github datasource

⚠️ Work in Progress ⚠️ 

Purpose of this plugin is to dogfood our backend plugin SDK. This plugin might or might not end up being published. 

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
