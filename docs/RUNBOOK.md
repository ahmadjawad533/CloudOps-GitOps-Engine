# Operations Runbook: GitOps & Kubernetes Delivery Pipeline

This runbook covers operational troubleshooting, rollout verification, rollbacks, and incident response procedures.

---

## 1. Verifying Application Health

### Cluster Workload Status
```bash
# Check pod rollout status in target namespace
kubectl rollout status deployment/prod-sample-app -n sample-app-production

# Check active pods and restart counts
kubectl get pods -n sample-app-production -l app.kubernetes.io/name=sample-app -o wide
```

### Checking Application Endpoints
```bash
# Verify health probe
curl -I http://<INGRESS_IP>/healthz

# Check readiness probe
curl -I http://<INGRESS_IP>/ready

# Inspect Prometheus metrics
curl http://<INGRESS_IP>/metrics
```

---

## 2. ArgoCD Operations

### Synchronizing Manually
If automated sync is disabled or out-of-sync:
```bash
# Using ArgoCD CLI
argocd app sync sample-app-production --prune

# Or check sync status
argocd app get sample-app-production
```

### Investigating Sync Failures & Drift
```bash
# View detailed diff between Git and Cluster
argocd app diff sample-app-production

# Check controller logs
kubectl logs -n argocd -l app.kubernetes.io/name=argocd-application-controller --tail=100 -f
```

---

## 3. Emergency Rollback Procedure

In GitOps, **Git is the single source of truth**. Never modify resources directly in the cluster (`kubectl edit` will be immediately reverted by ArgoCD self-healing).

### Standard Rollback via Git (Recommended)
1. Revert the commit in Git:
   ```bash
   git revert HEAD
   git push origin main
   ```
2. ArgoCD will detect the change within seconds and roll back the pods to the previous stable image tag.

### Emergency Instant Rollback via ArgoCD UI / CLI
If immediate action is required before Git CI finishes:
```bash
# Disable self-heal temporarily
argocd app set sample-app-production --self-heal=false

# Rollback to the previous deployment revision
argocd app rollback sample-app-production <PREVIOUS_REVISION_ID>
```
*Note: Make sure to commit the fix to Git afterwards so self-healing can be safely re-enabled.*

---

## 4. Common Alerts & Diagnostics

| Symptom | Probable Cause | Action |
|---|---|---|
| `CrashLoopBackOff` | Configuration error or failed health check | Inspect `kubectl logs <pod> -c app --previous` |
| `OOMKilled` | Memory request/limit exceeded | Check metrics and bump memory limits in `patches/deployment-patch.yaml` |
| ArgoCD `OutOfSync` | Manual change made to cluster or Git commit waiting sync | Verify Git commit history and trigger `argocd app sync` |
| `NetworkPolicy Block` | Service unable to reach external DNS or DB | Inspect `network-policy.yaml` egress rules |
