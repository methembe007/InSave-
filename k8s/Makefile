.PHONY: help apply delete verify generate-secrets check-secrets

help:
	@echo "InSavein Kubernetes Configuration"
	@echo ""
	@echo "Available targets:"
	@echo "  apply            - Apply all Kubernetes configurations"
	@echo "  delete           - Delete all Kubernetes resources"
	@echo "  verify           - Verify all resources are created"
	@echo "  generate-secrets - Generate secure random secrets"
	@echo "  check-secrets    - Check if secrets contain placeholder values"
	@echo "  namespace        - Create namespace only"
	@echo "  config           - Apply ConfigMap only"
	@echo "  secrets          - Apply Secrets only"
	@echo "  quotas           - Apply resource quotas only"
	@echo "  network          - Apply network policies only"
	@echo "  priority         - Apply priority classes only"

# Apply all configurations in correct order
apply: check-secrets
	@echo "Applying Kubernetes configurations..."
	kubectl apply -f namespace.yaml
	kubectl apply -f priority-class.yaml
	kubectl apply -f configmap.yaml
	kubectl apply -f secrets.yaml
	kubectl apply -f resource-quota.yaml
	kubectl apply -f network-policy.yaml
	@echo "✓ All configurations applied successfully"
	@echo ""
	@echo "Run 'make verify' to check the deployment"

# Delete all resources
delete:
	@echo "Deleting Kubernetes resources..."
	@read -p "Are you sure you want to delete all resources in the insavein namespace? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		kubectl delete -f network-policy.yaml --ignore-not-found=true; \
		kubectl delete -f resource-quota.yaml --ignore-not-found=true; \
		kubectl delete -f secrets.yaml --ignore-not-found=true; \
		kubectl delete -f configmap.yaml --ignore-not-found=true; \
		kubectl delete -f priority-class.yaml --ignore-not-found=true; \
		kubectl delete -f namespace.yaml --ignore-not-found=true; \
		echo "✓ All resources deleted"; \
	else \
		echo "Deletion cancelled"; \
	fi

# Verify all resources
verify:
	@echo "Verifying Kubernetes resources..."
	@echo ""
	@echo "=== Namespace ==="
	kubectl get namespace insavein
	@echo ""
	@echo "=== Priority Classes ==="
	kubectl get priorityclass | grep insavein
	@echo ""
	@echo "=== ConfigMap ==="
	kubectl get configmap -n insavein
	@echo ""
	@echo "=== Secrets ==="
	kubectl get secrets -n insavein
	@echo ""
	@echo "=== Resource Quotas ==="
	kubectl get resourcequota -n insavein
	@echo ""
	@echo "=== Limit Ranges ==="
	kubectl get limitrange -n insavein
	@echo ""
	@echo "=== Network Policies ==="
	kubectl get networkpolicies -n insavein
	@echo ""
	@echo "✓ Verification complete"

# Generate secure random secrets
generate-secrets:
	@echo "Generating secure random secrets..."
	@echo ""
	@echo "=== JWT Secret Key (64 bytes) ==="
	@openssl rand -base64 64 | tr -d '\n' && echo
	@echo ""
	@echo "=== Data Encryption Key (32 bytes) ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "=== Session Encryption Key (32 bytes) ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "=== Database Password ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "=== Replica Password ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "=== Postgres Superuser Password ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "=== Replication Password ==="
	@openssl rand -base64 32 | tr -d '\n' && echo
	@echo ""
	@echo "Copy these values to secrets.yaml before deploying!"

# Check if secrets contain placeholder values
check-secrets:
	@echo "Checking secrets.yaml for placeholder values..."
	@if grep -q "CHANGE_ME" secrets.yaml; then \
		echo ""; \
		echo "⚠️  WARNING: secrets.yaml contains placeholder values!"; \
		echo ""; \
		echo "Found the following placeholders:"; \
		grep "CHANGE_ME" secrets.yaml | sed 's/^/  /'; \
		echo ""; \
		echo "Please run 'make generate-secrets' to generate secure values"; \
		echo "and update secrets.yaml before deploying."; \
		echo ""; \
		exit 1; \
	else \
		echo "✓ No placeholder values found in secrets.yaml"; \
	fi

# Individual resource targets
namespace:
	kubectl apply -f namespace.yaml

config:
	kubectl apply -f configmap.yaml

secrets: check-secrets
	kubectl apply -f secrets.yaml

quotas:
	kubectl apply -f resource-quota.yaml

network:
	kubectl apply -f network-policy.yaml

priority:
	kubectl apply -f priority-class.yaml
