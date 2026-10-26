# Architecture Decision Record (ADR): GitOps & Kubernetes Delivery Pipeline

## Status
Accepted

## Context
Delivering software reliably at scale requires a clear boundary between CI (Continuous Integration: testing, building, scanning artifacts) and CD (Continuous Delivery: reconciling desired cluster state). Traditional push-based deployment models (where CI directly connects to Kubernetes clusters with high-privilege credentials) introduce serious security risks, potential credential leaks, and state drift.

## Architectural Decisions

### 1. GitOps Model via Pull-Based Reconciler (ArgoCD)
- **Decision**: The cluster reconciler (ArgoCD) continuously pulls the declarative state from Git and applies it to the cluster.
- **Benefits**:
  - Zero Kubernetes cluster admin credentials stored in GitHub Actions runners.
  - Automatic drift detection and self-healing when external factors alter cluster resources.
  - Complete auditable version control of all infrastructure and application states.

### 2. Multi-Environment Orchestration with Kustomize
- **Decision**: Use pure declarative Kustomize `base` and `overlays/{dev,staging,production}` instead of Helm templating for application resource definitions.
- **Benefits**:
  - Template-free, native Kubernetes YAML structures that are easy to validate with standard linters.
  - Clear inheritance: base definitions contain mandatory security contexts, labels, probes, and network policies; environment overlays only override replica counts, hostnames, and resources.

### 3. Hardened Container Architecture
- **Decision**: Multi-stage minimal container build using non-root UID 10001, statically compiled binaries, and dropped Linux capabilities.
- **Benefits**:
  - Extremely small image footprint (<20MB).
  - Immunity to standard privilege escalation attacks; read-only root filesystem guarantees immutability at runtime.

### 4. Zero-Trust Network & Pod Disruption Policies
- **Decision**: Every service must be accompanied by a Kubernetes `NetworkPolicy` and `PodDisruptionBudget`.
- **Benefits**:
  - Default deny for unauthorized pod-to-pod lateral movements.
  - Guarantees minimum available pods during rolling upgrades and node draining.
