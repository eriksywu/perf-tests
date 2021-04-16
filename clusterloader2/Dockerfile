FROM ubuntu:18.04 as curler
RUN apt update && apt install -y curl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl 

FROM golang as gobuilder

COPY . /clusterloader2
RUN cd /clusterloader2/cmd && go build -o /clusterloader

FROM alpine:latest 

ARG ARTIFACTSDIR=_artifacts
ARG CONFIGDIR=_configs
WORKDIR /

RUN mkdir -p ${CONFIGDIR}
COPY --chown=nonroot:nonroot  testing/node-throughput/ ${CONFIGDIR}
COPY --from=gobuilder /clusterloader .
COPY --from=gobuilder /clusterloader2/run-e2e.sh .
RUN chmod +x ./run-e2e.sh
COPY --from=curler /kubectl /usr/bin/kubectl

ENTRYPOINT [ "./run-e2e.sh" ]