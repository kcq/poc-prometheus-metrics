here="$(dirname "$BASH_SOURCE")"
cd $here/../..
export GOPATH=$HOME/go
go run cmd/prom-metrics-service/main.go