here="$(dirname "$BASH_SOURCE")"
cd $here
docker run -it --rm --name="prom-metrics-service" -p 7000:7000 poc/prom-metrics-service

