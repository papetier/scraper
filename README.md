# Papetier's scraper

Scrapes the [arXiv API](https://arxiv.org/help/api/basics) to index the papers!

## Configuration

The configuration is read from the following sources (in order of priority):

1. Environment variables explicitly set
2. `.env` file: placed in the root folder and named after the environment. Defaults to `local.env` (recommended)
3. Default values as defined in [`defaults.go`](./pkg/config/defaults.go)

To define an environment file, you can duplicate the [`example.env`](./example.env) file and adapt it to your settings.

If you want to keep several environment settings, you can save them under different names (i.e. `local.env`, `dev.env`, `prod.env`). When running a command, you can pass the corresponding environment name with the `SCRAPER_ENVIRONMENT` environment variable:

```shell
SCRAPER_ENVIRONMENT="prod" go run cmd/server/main.go
```

If `SCRAPER_ENVIRONMENT` isn't set, the config will read the settings from `local.env` by default.

> All `.env` (except `example.env`) files are ignored by git to avoid exposing any credentials.

## Commands

The repository exposes 1 command defined in the `cmd` folder.

### `scraper`

To launch the scraper.

---

## Development

### Build

The scraper uses [Mage](https://magefile.org/) for common development tasks, including building the binary. Mage is a
pure Go library with no dependencies which offers a much easier API and syntax than the classic Makefile.

To install `mage`, run the following command:

```shell
go install github.com/magefile/mage
```

#### Prod build

For production, you can use the `prod` target:

````shell
mage build:prod
````

#### Local build

To build a binary compatible with your local machine, you can use the `local` target:

````shell
mage build:local
````
