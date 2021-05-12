FROM ubuntu:18.04 as curler
RUN apt update && apt install -y curl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl 

FROM golang as gobuilder

COPY . /clusterloader2
RUN cd /clusterloader2/cmd && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /clusterloader
RUN cd /clusterloader2/servlet && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /servlet
RUN cd 

FROM alpine:3.7 

ARG ARTIFACTSDIR=_artifacts
ARG CONFIGDIR=_configs
WORKDIR /

RUN mkdir -p ${CONFIGDIR}
COPY --chown=nonroot:nonroot  testing/ ${CONFIGDIR}
COPY --from=gobuilder /clusterloader .
COPY --from=gobuilder /servlet .
COPY --from=gobuilder /clusterloader2/pkg /pkg
COPY --from=gobuilder /clusterloader2/run-e2e.sh .
RUN chmod +x ./run-e2e.sh
COPY --from=curler /kubectl /usr/bin/kubectl

RUN apk add --no-cache bash
RUN apk add --no-cache \
    libc6-compat
ENTRYPOINT [ "/bin/bash" ]
CMD ["-c" ,"/run-e2e.sh", "--run-from-cluster=true"]