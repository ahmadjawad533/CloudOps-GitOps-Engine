# GitOps & Kubernetes Production Delivery Pipeline

An enterprise-ready, pull-based GitOps deployment platform and cloud-native CI/CD delivery engine powered by **ArgoCD**, **Kustomize**, **GitHub Actions**, **Trivy**, and **Kubernetes**.

---

## 🌟 Key Features

- **Declarative GitOps Engine**: Fully automated sync, drift detection, and self-healing with ArgoCD `Application` & `AppProject` configurations.
- **Multi-Environment Orchestration**: Kustomize `base` with `overlays/dev`, `overlays/staging`, and `overlays/production`.
- **Zero-Trust Hardened Containers**: Multi-stage minimal footprint build running as non-root (`UID 10001`), read-only root filesystems, and dropped capabilities (`ALL`).
- **Resilience & High Availability**: PodDisruptionBudgets (`PDB`), Pod Anti-Affinity, HorizontalPodAutoscalers (`HPA`), and zero-trust `NetworkPolicy` egress/ingress isolation.
- **Enterprise CI/CD & Security Gates**:
  - GitHub Actions matrix testing & Go vet/formatting.
  - Multi-arch Docker Buildx verification (`linux/amd64`, `linux/arm64`).
  - Container vulnerability CVE & Kubernetes manifest misconfiguration scanning with **Trivy**.
  - Strict OpenAPI schema validation with **Kubeconform**.
  - Automated image promotion bot bumping staging overlays via Git pull requests.
- **Local Reproducibility**: One-click KinD bootstrap script with Ingress-NGINX and ArgoCD ready out-of-the-box.

---

## 📁 Repository Structure

```
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Lint, vet, unit test, and multi-arch build
│       ├── security-scan.yml      # Trivy container & K8s config vulnerability scanner
│       ├── k8s-lint.yml           # Kubeconform strict OpenAPI schema verification
│       └── gitops-sync.yml        # Automatic promotion PR bot for staging image tags
├── app/                           # Production-ready Go microservice
│   ├── main.go                    # HTTP handlers with healthz, ready, metrics, and SIGTERM trap
│   ├── main_test.go               # Unit test coverage
│   ├── Dockerfile                 # Multi-stage security hardened container
│   └── go.mod
├── gitops/
│   ├── argocd/                    # Declarative ArgoCD manifests
│   │   ├── app-project.yaml       # Multi-tenant AppProject boundaries
│   │   ├── app-dev.yaml           # Dev sync definition
│   │   ├── app-staging.yaml       # Staging auto-sync definition
│   │   └── app-production.yaml    # Production sync waves & safeguards
│   └── kustomize/                 # Declarative K8s manifest hierarchy
│       ├── base/                  # Deployment, Service, Ingress, HPA, PDB, NetworkPolicy
│       └── overlays/
│           ├── dev/               # Lightweight specs & debug loglevel
│           ├── staging/           # Pre-production test configuration
│           └── production/        # HA replicas, anti-affinity, TLS ingress
├── scripts/
│   ├── setup-local-kind.sh        # One-click KinD cluster setup with Ingress and ArgoCD
│   ├── validate-manifests.sh      # Validates all Kustomize manifests locally
│   └── teardown.sh                # Local cluster teardown
├── docs/
│   ├── ARCHITECTURE.md            # Architecture Decision Records (ADR)
│   └── RUNBOOK.md                 # Operational incident handling and rollback guide
└── Makefile                       # Developer workflow commands
```

---

## 🚀 Quick Start

### 1. Validate Manifests
Run the validation script to test all Kustomize overlays:
```bash
make validate-manifests
```

### 2. Build Container Image
Build the hardened container locally:
```bash
make build-docker
```

### 3. Spin Up Local KinD GitOps Cluster
Provision a local Kubernetes cluster equipped with NGINX ingress and ArgoCD:
```bash
make kind-up
```

### 4. Deploy Declarative ArgoCD Applications
Connect the repository overlays to ArgoCD:
```bash
make deploy-argocd-apps
```

### 5. Access ArgoCD Dashboard
```bash
# Port-forward the ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8081:443

# Retrieve initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo
```
Open [https://localhost:8081](https://localhost:8081) in your browser.

---

## 🔄 GitOps Promotion Flow

1. Developer merges code into `main`.
2. `ci.yml` tests, compiles, and tags container image (e.g., `v1.2.0`).
3. `security-scan.yml` runs Trivy vulnerability scanning.
4. `gitops-sync.yml` updates `gitops/kustomize/overlays/staging/kustomization.yaml` with the new tag and opens a pull request.
5. Once merged, ArgoCD automatically reconciles the change onto the staging namespace with zero downtime.
6. Following qualification, a promotion PR to `overlays/production` rolls out the new version to production.

---

## 📖 Documentation
- [Architecture Decision Record (ADR)](docs/ARCHITECTURE.md)
- [Operations & Incident Runbook](docs/RUNBOOK.md)
