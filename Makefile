default: build

build_container:
	docker run --rm -v "$(CURDIR)":/go/src/github.com/kcq/poc-prometheus-metrics -w /go/src/github.com/kcq/poc-prometheus-metrics golang:1.11.4 make build
	docker build -t poc/prom-metrics-service -f $(CURDIR)/deployments/Dockerfile $(CURDIR)

build_container_only:
	'$(CURDIR)/scripts/build_container.sh'

build:
	'$(CURDIR)/scripts/build.sh'

build_with_local_gopath:
	'$(CURDIR)/scripts/local_gopath.sh'
	. '$(CURDIR)/scripts/env.sh'; '$(CURDIR)/scripts/build.sh'

build_linux:
	'$(CURDIR)/scripts/build_linux_local.sh'

clean:
	'$(CURDIR)/scripts/clean.sh'

fmt:
	'$(CURDIR)/scripts/fmt.sh'

inspect:
	'$(CURDIR)/scripts/inspect.sh'

run_container:
	docker run -it --rm --name="prom-metrics-service" -p 7000:7000 poc/prom-metrics-service

run:
	$(CURDIR)/bin/prom-metrics-service

up:
	docker-compose -p prom-metrics-service -f $(CURDIR)/deployments/docker-compose.yaml up -d

down:
	docker-compose -p prom-metrics-service -f $(CURDIR)/deployments/docker-compose.yaml down

ps:
	docker-compose -p prom-metrics-service -f $(CURDIR)/deployments/docker-compose.yaml ps

tail:
	docker-compose -p prom-metrics-service -f $(CURDIR)/deployments/docker-compose.yaml logs -f

.PHONY: default build_container build_with_local_gopath build clean fmt up down