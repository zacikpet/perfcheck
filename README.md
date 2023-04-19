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
    "host": "localhost:8080",
    "paths": {
        "/": {
            "get": {
                "summary": "Hello",
                "x-perfcheck": {
                    "errorRate": ["rate < 0.1"],
                    "latency": ["avg < 50"]
                }
            }
        }
    },
    "x-perfcheck": {
        "duration": "10s",
        "vus": 5
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
perfcheck test --source openapi --docsUrl http://localhost:8080/swagger/doc.json
```

### Global options

| key | meaning |
| - | - |
| `vus` | [k6 vus](https://k6.io/docs/using-k6/k6-options/reference/#vus) |
| `duration` | [k6 duration](https://k6.io/docs/using-k6/k6-options/reference/#duration) |
| `stages` | [k6 stages](https://k6.io/docs/using-k6/k6-options/reference/#stages) |

### Per-endpoint options

| key | meaning |
| - | - |
| `latency` | array of objectives for the latency in ms |
| `errorRate` | array of objectives for the error rate |

Latency objectives can contain the following metrics:

| metric | meaning |
| - | - |
| `avg_stat` | Statistically significant mean value |
| `avg`, `med` | Average/median value |
| `min`, `max` | Minimum/maximum value |
| `p(N)` | Specific percentile `N` |

For more info, see [k6 trends](https://k6.io/docs/javascript-api/k6-metrics/trend/).

Error rate objectives can contain the following metrics:

| metric | meaning |
| - | - |
| `rate` | ratio of errors to all requests |

## How to use with Google Cloud SLOs

1. Create a [Google Cloud SLO](https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring/ui/create-slo)



### Support table

Currently supported SLI metrics are `Availability` and `Latency`.

| Metric | Request-based SLI | Window-based SLI |
| - | - | - |
| Availability | ✅ | ❌ |
| Latency | ✅ | ❌ |
| Response size | ✅ | ❌ |
| Custom | ❌ | ❌ |

2. Set your default google credentials using the Google Cloud CLI

```bash
gcloud auth login
```

3. Install and run `perfcheck`

```bash
go install github.com/zacikpet/perfcheck
```

```bash
perfcheck test \
    --source gcloud \
    --gcloudProjectId [your-project-id] \
    --gcloudServiceId [id-of-your-monitoring-service] \
    --gcloudServiceUrl [url-of-your-service]
```

**Warning: Using with google cloud only tests your root path!**


