#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="gitops-lab"

echo "=== Tearing Down GitOps Lab Cluster ==="
if command -v kind &> /dev/null; then
  kind delete cluster --name "${CLUSTER_NAME}" || true
elif [ -f /tmp/kind ]; then
  /tmp/kind delete cluster --name "${CLUSTER_NAME}" || true
else
  echo "kind binary not found. Skipping cluster deletion."
fi

echo "Teardown complete."
