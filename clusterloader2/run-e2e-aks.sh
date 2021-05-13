#!/bin/bash

./clusterloader --alsologtostderr --v=2 --run-from-cluster=true --provider=aks --testconfig=/_configs/node-throughput/config.yaml --testoverrides=/_configs/overrides/node_containerd.yaml --report-dir=/results "$@"
echo "done clusterloader tests"
touch /results/done #sentinel file
while sleep 3600; do :; done