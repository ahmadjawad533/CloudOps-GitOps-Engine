#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "=== Validating GitOps Kustomize Manifests ==="

validate_overlay() {
  local target="$1"
  echo "Checking: ${target}..."
  if command -v kustomize &> /dev/null; then
    kustomize build "${ROOT_DIR}/${target}" > /dev/null
    echo "  [OK] Kustomize build succeeded for ${target}"
  elif command -v kubectl &> /dev/null; then
    kubectl kustomize "${ROOT_DIR}/${target}" > /dev/null
    echo "  [OK] kubectl kustomize succeeded for ${target}"
  else
    echo "  [NOTICE] Neither kustomize nor kubectl found in PATH. Checking YAML syntax with python..."
    python3 -c "
import os, sys, yaml
for root, _, files in os.walk('${ROOT_DIR}/${target}'):
  for f in files:
    if f.endswith(('.yaml', '.yml')):
      p = os.path.join(root, f)
      try:
        list(yaml.safe_load_all(open(p)))
      except Exception as e:
        print(f'Error in {p}: {e}', file=sys.stderr)
        sys.exit(1)
"
    echo "  [OK] YAML syntax validation passed for ${target}"
  fi
}

validate_overlay "gitops/kustomize/base"
validate_overlay "gitops/kustomize/overlays/dev"
validate_overlay "gitops/kustomize/overlays/staging"
validate_overlay "gitops/kustomize/overlays/production"

echo "=== All Manifests Validated Successfully! ==="
