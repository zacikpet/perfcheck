FROM golang:latest

RUN go install github.com/zacikpet/perf-check@latest

ENV SOURCE=openapi
ENV DOCS_URL=http://host.docker.internal:8080

ENTRYPOINT ["perf-check"]