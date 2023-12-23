# GoLang Backend Think Laundry

---

## IDE Setup

### VSCode

* gofumpt. see [vscode gfumpt setup](https://github.com/mvdan/gofumpt#visual-studio-code)

* golangci-lint

    ```json
    "go.lintTool":"golangci-lint",
    "go.lintFlags": [
    "--fast"
    ]
    ```

### GoLand

* gofumpt. see [goland gofumpt setup](https://github.com/mvdan/gofumpt#goland)

* golangci-lint

  1. Install `File Watchers` plugin from Intellij Plugin Marketplace.

  2. Settings -> Tools -> File Watchers -> Click + -> `golangci-lint`

  3. Arguments : `run --fix --fast $FileDir$`

## Local development

* Copy file `config.example.yml` to `config.yml`
* Setup database (local) to create database schema equal with config `(schema: "think_laundry-dev")`
* Setup Makefile param for migration connection database ex: `PG_DB_URL=postgresql://postgres:postgres@127.0.0.1:5439/think_laundry-dev?sslmode=disable`
* run mod=vendor dependency with `make deps`
* Up / Run migrate with `make run pg.migrate.up` (Use WSL for OS Windows)
* run with `make run-service-local`

## API Docs
### [Postman API Docs]

## License
[© 2023 thinkIT]
