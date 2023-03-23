**perfcheck** is a testing automation tool used to generate k6 benchmarks from service level objectives (SLOs)

## How to use with OpenAPI

1. Add custom `x-perfcheck` annotations to your OpenAPI documentation. Example OpenAPI 2.0 JSON schema:

```json
{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "title": "Example API",
    },
    "paths": {
        "/": {
            "get": {
                "summary": "Hello World",
                "x-perfcheck": {
                    "latency": ["avg < 50"]
                }
            }
        }
    },
    "x-perfcheck": {
        "vus": 100,
        "duration": "10s"
    }
}
```
In this schema, the requirement for `GET /` is that the average latency stays below 50ms. The global `x-perfcheck` property defines the duration of the test and number of virtual users. You can also define mulitple test stages (see [k6 stages](https://k6.io/docs/using-k6/k6-options/reference/#stages)).

2. Run the server with OpenAPI docs (for example: http://localhost:8080/swagger/doc.json)

3. Install and run perfcheck

```bash
go install github.com/zacikpet/perfcheck
```

```bash
perfcheck test --source=openapi --docsUrl=http://localhost:8080/swagger/doc.json
```

## How to use with Google Cloud SLOs

TODO
