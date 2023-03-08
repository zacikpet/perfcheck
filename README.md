# perf-check

This repository contains:
- `perf-check`: a tool that generates benchmarks from a swagger documentation
- a testing REST app that serves swagger documentation

## How to run

### 1. Generate docs for `test-app` (optional)

Inside the `test-app` subdirectory, run:

    swag init

to generate swagger documentation from the source code. To install `swag`, see https://github.com/swaggo/swag.

### 2. Run `test-app`

Inside the `test-app` subdirectory, run:

    go run main.go

to launch a testing HTTP server. The docs will be served at:

- http://localhost:8080/swagger/index.html (HTML)
- http://localhost:8080/swagger/doc.json (JSON)

### 3. Add environment variables

Create a `.env` file with the following contents:

    DOCS_URL=http://localhost:8080/swagger/doc.json

This will tell `perf-check` the location and description of the service.

### 4. Generate and run benchmark file

Inside the root directory, run:

    go run *.go

to generate a `benchmarks/benchmark.js` file. The file is automatically ran using the `k6` binary in your path. If you do not have `k6` installed, this operation will fail.