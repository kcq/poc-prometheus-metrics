# PoC - Go Service Metrics Instrumentation with Prometheus

The PoC shows how to instrument a Go service and how to create custom Prom registry and HTTP metrics handlers.

The Docker Compose config for the PoC starts Prometheus, Grafana and the instrumented PoC service.

The PoC services exposes a dummy API on port 7000. The Prometheus metrics is exposed using the `/metrics` endpoint.

## Demo Steps

* Build the PoC service and its container image (`make build_container` or you can build them separately with `make build_linux` and `make build_container_only`).
* Start Docker Compose (either with `make up` or `all_up.command` in the script directory).
* Make HTTP calls to the PoC service (`curl http://localhost:7000`).
* Use Prom UI (`http://localhost:9090`) or Grafana (`http://localhost:3000`) to see the generated metrics
* Access the `/metrics` endpoint (`curl http://localhost:7000/metrics`) if you want to see the raw metrics exposed by the PoC service.

Note:

The Grafana dashboard requires authentication. The PoC Grafana setup uses the default credentials (user: admin , password: admin). Grafana will prompt you to change the password on first login (you can skip this step). The PoC Grafana is not configured, so you'll need to add the Prometheus data source (set url to `http://prometheus:9090` and type to `prometheus`) and create your own dashboard (future PoC enhancement).

### Run Standalone Service in Docker (no Prometheus and Grafana)

* Build the PoC service: `make build_container`
* Run the PoC service: `make run_container`
* Make HTTP calls to the PoC service (`curl http://localhost:7000`)
* Access the Prometheus metrics exposed by the service (`curl http://localhost:7000/metrics`)

Note:

* The `build_container` make command uses a Dockerized build, so you don't need to have Go installed locally.
* Docker is a requirement.

### Run Standalone Service Natively

* Build the PoC service: `make`
* Run the PoC service: `make run`
* Make HTTP calls to the PoC service (`curl http://localhost:7000`)
* Access the Prometheus metrics exposed by the service (`curl http://localhost:7000/metrics`)

## Overview

* `cmd/prom-metrics-service/main.go` - the PoC code
* `configs` - Prometheus config for the PoC
* `deployments` - Docker and Docker Compose files for the PoC
* `scripts` - helper build and run scripts

## Notes

The PoC also leverages static Go build and Docker containers build from `scratch`, so the images includes only the PoC service binary and nothing else.
