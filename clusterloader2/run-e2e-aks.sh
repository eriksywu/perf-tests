#!/bin/bash

./clusterloader --alsologtostderr --v=2 --run-from-cluster=true --provider=aks --testconfig=/_configs/node-throughput/config.yaml --testoverrides=/_configs/overrides/node_containerd.yaml --report-dir=/results "$@"
if [ $? -ne 0 ]; then
    status="failed"
else
    status="done"
fi

echo "done clusterloader tests with status ${status}"
curl -X POST http://localhost:8080/results?status=${status}
while sleep 3600; do :; done