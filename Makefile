.PHONY: all test build-docker validate-manifests kind-up kind-down argocd-apps help

IMAGE_NAME ?= ghcr.io/cloudicraft/gitops-sample-service
TAG ?= 1.0.0

help:
	@echo "Available Make targets:"
	@echo "  make test               - Run microservice unit tests"
	@echo "  make build-docker       - Build hardened multi-stage Docker container"
	@echo "  make validate-manifests - Lint and validate Kustomize manifests"
	@echo "  make kind-up            - Launch local KinD cluster with NGINX & ArgoCD"
	@echo "  make kind-down          - Teardown local KinD cluster"
	@echo "  make deploy-argocd-apps - Apply declarative ArgoCD application definitions"

test:
	@echo "==> Running microservice unit tests..."
	cd app && go test -v -race ./...

build-docker:
	@echo "==> Building container $(IMAGE_NAME):$(TAG)..."
	docker build -t $(IMAGE_NAME):$(TAG) ./app

validate-manifests:
	@echo "==> Validating Kubernetes manifests..."
	./scripts/validate-manifests.sh

kind-up:
	@echo "==> Setting up local KinD cluster..."
	./scripts/setup-local-kind.sh

kind-down:
	@echo "==> Tearing down local cluster..."
	./scripts/teardown.sh

deploy-argocd-apps:
	@echo "==> Applying ArgoCD declarative AppProject & Applications..."
	kubectl apply -f gitops/argocd/app-project.yaml
	kubectl apply -f gitops/argocd/app-dev.yaml
	kubectl apply -f gitops/argocd/app-staging.yaml
	kubectl apply -f gitops/argocd/app-production.yaml
