# perf-check


This repository contains:
- `perf-check`: a tool that generates benchmarks from OpenApi documentation
- a testing REST app that serves OpenAPI documentation

## How to use

### 1. Environment variables

Add the following variables to your environment (`.env`)

    SOURCE=openapi
    DOCS_URL={URL of your OpenAPI documentation in the JSON format}

### 4. Generate and run benchmark file

Inside the root directory, run:

    go run *.go

to generate a `benchmarks/benchmark.js` file. The file is automatically ran using the `k6` binary in your path. If you do not have `k6` installed, this operation will fail.
