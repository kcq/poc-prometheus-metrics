FROM scratch

WORKDIR /opt/poc/bin
COPY bin .

EXPOSE 7000

ENTRYPOINT ["/opt/poc/bin/prom-metrics-service"]