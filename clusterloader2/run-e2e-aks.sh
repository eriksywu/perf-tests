#!/bin/bash

./clusterloader --alsologtostderr --v=2 --run-from-cluster=true --provider=aks --testconfig=/_configs/node-throughput/config.yaml --testoverrides=/_configs/overrides/node_containerd.yaml --report-dir=/results "$@" || exit 1
echo "done clusterloader tests"
while sleep 3600; do :; done