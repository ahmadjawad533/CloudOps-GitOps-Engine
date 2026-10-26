#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="gitops-lab"

echo "=== [1/4] Checking Prerequisites ==="
for cmd in docker; do
  if ! command -v "$cmd" &> /dev/null; then
    echo "Error: $cmd is required but not installed." >&2
    exit 1
  fi
done

if ! command -v kind &> /dev/null; then
  echo "kind not found. Downloading kind binary to /tmp..."
  KIND_VERSION="v0.22.0"
  curl -Lo /tmp/kind "https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64"
  chmod +x /tmp/kind
  KIND_CMD="/tmp/kind"
else
  KIND_CMD="kind"
fi

echo "=== [2/4] Provisioning KinD Cluster with Port Forwarding ==="
cat <<EOF > /tmp/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: ${CLUSTER_NAME}
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
    protocol: TCP
  - containerPort: 443
    hostPort: 8443
    protocol: TCP
EOF

if ${KIND_CMD} get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "Cluster ${CLUSTER_NAME} already exists."
else
  ${KIND_CMD} create cluster --config /tmp/kind-config.yaml
fi

echo "=== [3/4] Installing Ingress-NGINX Controller ==="
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml || true

echo "=== [4/4] Installing ArgoCD ==="
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

echo "Waiting for ArgoCD components to be ready..."
echo "To port-forward ArgoCD Web UI run:"
echo "  kubectl port-forward svc/argocd-server -n argocd 8081:443"
echo "Default username: admin"
echo "To retrieve admin initial password:"
echo "  kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d && echo"
echo ""
echo "=== KinD GitOps Cluster Ready! ==="
