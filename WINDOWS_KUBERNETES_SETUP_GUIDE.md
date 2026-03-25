# InSavein Platform - Windows Kubernetes Setup Guide

Complete step-by-step guide for setting up and deploying the InSavein Platform to Kubernetes on Windows using Command Prompt.

## Table of Contents
- [Prerequisites Installation](#prerequisites-installation)
- [Kubernetes Cluster Options](#kubernetes-cluster-options)
- [Local Development Setup](#local-development-setup)
- [Cloud Deployment Setup](#cloud-deployment-setup)
- [Deploying the Application](#deploying-the-application)
- [Verification & Testing](#verification--testing)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites Installation

### 1. Install kubectl (Kubernetes CLI)

**Method A: Using Chocolatey (Recommended)**

First, install Chocolatey (run Command Prompt as Administrator):
```cmd
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "[System.Net.ServicePointManager]::SecurityProtocol = 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
```

Then install kubectl:
```cmd
choco install kubernetes-cli
```

**Method B: Manual Installation**

1. Download kubectl from: https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/
2. Download the latest release:
   ```cmd
   curl.exe -LO "https://dl.k8s.io/release/v1.28.0/bin/windows/amd64/kubectl.exe"
   ```
3. Move to a directory in your PATH:
   ```cmd
   move kubectl.exe C:\Windows\System32\
   ```

**Verify Installation:**
```cmd
kubectl version --client
```

Expected output: `Client Version: v1.28.0` or higher

---

### 2. Install OpenSSL (For Generating Secrets)

**Using Chocolatey:**
```cmd
choco install openssl
```

**Manual Installation:**
1. Download from: https://slproweb.com/products/Win32OpenSSL.html
2. Install "Win64 OpenSSL v3.x.x"
3. Add to PATH via System Properties > Environment Variables
   - Add `C:\Program Files\OpenSSL-Win64\bin` to PATH

**Verify Installation:**
```cmd
openssl version
```

---

### 3. Install Helm (Kubernetes Package Manager)

**Using Chocolatey:**
```cmd
choco install kubernetes-helm
```

**Manual Installation:**
1. Download from: https://github.com/helm/helm/releases
2. Extract `helm.exe` to a directory in PATH
3. Verify:
   ```cmd
   helm version
   ```

---

### 4. Install Git for Windows

```cmd
choco install git
```

Or download from: https://git-scm.com/download/win

---

## Kubernetes Cluster Options

Choose one of the following options based on your needs:

### Option 1: Docker Desktop (Easiest for Local Development)

**Requirements:**
- Windows 10/11 Pro, Enterprise, or Education
- WSL 2 enabled
- 8GB RAM minimum (16GB recommended)

**Setup Steps:**

1. **Install Docker Desktop**
   - Download from: https://www.docker.com/products/docker-desktop/
   - Run installer
   - Restart computer

2. **Enable Kubernetes**
   - Open Docker Desktop
   - Go to Settings → Kubernetes
   - Check "Enable Kubernetes"
   - Click "Apply & Restart"
   - Wait 5-10 minutes for Kubernetes to start

3. **Verify Installation**
   ```cmd
   kubectl cluster-info
   kubectl get nodes
   ```

**Pros:**
- Easy setup
- Good for local development
- Integrated with Docker

**Cons:**
- Single-node cluster
- Limited resources
- Not suitable for production

---

### Option 2: Minikube (Flexible Local Development)

**Requirements:**
- Hyper-V or VirtualBox
- 4GB RAM minimum (8GB recommended)

**Setup Steps:**

1. **Install Minikube**
   ```cmd
   choco install minikube
   ```

2. **Start Minikube**
   ```cmd
   REM Using Hyper-V (Windows Pro/Enterprise)
   minikube start --driver=hyperv --memory=8192 --cpus=4

   REM Or using VirtualBox
   minikube start --driver=virtualbox --memory=8192 --cpus=4

   REM Or using Docker
   minikube start --driver=docker --memory=8192 --cpus=4
   ```

3. **Verify Installation**
   ```cmd
   kubectl cluster-info
   kubectl get nodes
   ```

4. **Enable Addons**
   ```cmd
   minikube addons enable ingress
   minikube addons enable metrics-server
   minikube addons enable dashboard
   ```

**Pros:**
- More flexible than Docker Desktop
- Supports multiple drivers
- Good addon ecosystem

**Cons:**
- Requires virtualization
- Single-node cluster
- Not for production

---

### Option 3: Amazon EKS (Production - AWS)

**Requirements:**
- AWS Account
- AWS CLI installed
- eksctl tool

**Setup Steps:**

1. **Install AWS CLI**
   ```cmd
   choco install awscli
   ```

2. **Configure AWS Credentials**
   ```cmd
   aws configure
   REM Enter: Access Key ID, Secret Access Key, Region, Output format
   ```

3. **Install eksctl**
   ```cmd
   choco install eksctl
   ```

4. **Create EKS Cluster**
   ```cmd
   eksctl create cluster ^
     --name insavein-cluster ^
     --region us-east-1 ^
     --nodegroup-name standard-workers ^
     --node-type t3.medium ^
     --nodes 3 ^
     --nodes-min 2 ^
     --nodes-max 5 ^
     --managed
   ```

   This takes 15-20 minutes.

5. **Verify Cluster**
   ```cmd
   kubectl get nodes
   kubectl get svc
   ```

**Pros:**
- Production-ready
- Managed service
- Auto-scaling
- High availability

**Cons:**
- Costs money
- More complex setup
- Requires AWS knowledge

---

### Option 4: Azure AKS (Production - Azure)

**Requirements:**
- Azure Account
- Azure CLI installed

**Setup Steps:**

1. **Install Azure CLI**
   ```cmd
   choco install azure-cli
   ```

2. **Login to Azure**
   ```cmd
   az login
   ```

3. **Create Resource Group**
   ```cmd
   az group create ^
     --name insavein-rg ^
     --location eastus
   ```

4. **Create AKS Cluster**
   ```cmd
   az aks create ^
     --resource-group insavein-rg ^
     --name insavein-cluster ^
     --node-count 3 ^
     --node-vm-size Standard_D2s_v3 ^
     --enable-addons monitoring ^
     --generate-ssh-keys
   ```

5. **Get Credentials**
   ```cmd
   az aks get-credentials ^
     --resource-group insavein-rg ^
     --name insavein-cluster
   ```

6. **Verify Cluster**
   ```cmd
   kubectl get nodes
   ```

**Pros:**
- Production-ready
- Managed service
- Azure integration
- Good monitoring

**Cons:**
- Costs money
- Requires Azure knowledge

---

### Option 5: Google GKE (Production - Google Cloud)

**Requirements:**
- Google Cloud Account
- gcloud CLI installed

**Setup Steps:**

1. **Install gcloud CLI**
   - Download from: https://cloud.google.com/sdk/docs/install
   - Run installer
   - Initialize: `gcloud init`

2. **Create GKE Cluster**
   ```cmd
   gcloud container clusters create insavein-cluster ^
     --zone us-central1-a ^
     --num-nodes 3 ^
     --machine-type n1-standard-2 ^
     --enable-autoscaling ^
     --min-nodes 2 ^
     --max-nodes 5
   ```

3. **Get Credentials**
   ```cmd
   gcloud container clusters get-credentials insavein-cluster ^
     --zone us-central1-a
   ```

4. **Verify Cluster**
   ```cmd
   kubectl get nodes
   ```

**Pros:**
- Production-ready
- Excellent performance
- Good integration with Google services

**Cons:**
- Costs money
- Requires GCP knowledge

---

## Local Development Setup

This section covers deploying to Docker Desktop or Minikube for local development.

### Step 1: Clone Repository

```cmd
REM Navigate to your projects folder
cd C:\Users\%USERNAME%\Projects

REM Clone repository
git clone <repository-url> insavein-platform
cd insavein-platform
```

### Step 2: Build Docker Images

```cmd
REM Build all service images
docker build -t insavein/auth-service:latest ./auth-service
docker build -t insavein/user-service:latest ./user-service
docker build -t insavein/savings-service:latest ./savings-service
docker build -t insavein/budget-service:latest ./budget-service
docker build -t insavein/goal-service:latest ./goal-service
docker build -t insavein/education-service:latest ./education-service
docker build -t insavein/notification-service:latest ./notification-service
docker build -t insavein/analytics-service:latest ./analytics-service
docker build -t insavein/frontend:latest ./frontend
```

**Or use the batch script:**
```cmd
docker-build-test.bat
```

### Step 3: Generate Secrets

```cmd
cd k8s

REM Generate secure random values
echo JWT Secret:
openssl rand -base64 64

echo.
echo Data Encryption Key:
openssl rand -base64 32

echo.
echo Session Encryption Key:
openssl rand -base64 32

echo.
echo Database Password:
openssl rand -base64 32

echo.
echo Replica Password:
openssl rand -base64 32

echo.
echo Postgres Password:
openssl rand -base64 32

echo.
echo Replication Password:
openssl rand -base64 32
```

### Step 4: Update Secrets File

1. Open `k8s/secrets.yaml` in a text editor
2. Replace ALL `CHANGE_ME_*` placeholders with generated values
3. Save the file
4. **IMPORTANT:** Do NOT commit this file to Git!

### Step 5: Deploy to Kubernetes

```cmd
cd k8s

REM Create namespace
kubectl apply -f namespace.yaml

REM Create priority classes
kubectl apply -f priority-class.yaml

REM Create ConfigMap
kubectl apply -f configmap.yaml

REM Create Secrets
kubectl apply -f secrets.yaml

REM Apply resource quotas
kubectl apply -f resource-quota.yaml

REM Apply network policies
kubectl apply -f network-policy.yaml

REM Deploy PostgreSQL
kubectl apply -f postgres-statefulset.yaml

REM Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n insavein --timeout=300s

REM Deploy backend services
kubectl apply -f auth-service-deployment.yaml
kubectl apply -f user-service-deployment.yaml
kubectl apply -f savings-service-deployment.yaml
kubectl apply -f budget-service-deployment.yaml
kubectl apply -f goal-service-deployment.yaml
kubectl apply -f education-service-deployment.yaml
kubectl apply -f notification-service-deployment.yaml
kubectl apply -f analytics-service-deployment.yaml

REM Deploy frontend
kubectl apply -f frontend-deployment.yaml

REM Deploy ingress (optional for local)
kubectl apply -f ingress.yaml
```

### Step 6: Access the Application

**For Docker Desktop:**
```cmd
REM Port forward to access services
kubectl port-forward -n insavein svc/frontend 3000:3000
kubectl port-forward -n insavein svc/auth-service 8080:8080
```

Access at: http://localhost:3000

**For Minikube:**
```cmd
REM Get Minikube IP
minikube ip

REM Access via NodePort or use tunnel
minikube service frontend -n insavein
```

---

## Cloud Deployment Setup

This section covers deploying to EKS, AKS, or GKE for production.

### Step 1: Prepare Container Registry

**For AWS ECR:**
```cmd
REM Create ECR repositories
aws ecr create-repository --repository-name insavein/auth-service
aws ecr create-repository --repository-name insavein/user-service
REM ... repeat for all services

REM Login to ECR (replace <account-id> and region)
aws ecr get-login-password --region us-east-1 > ecr-password.txt
type ecr-password.txt | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
del ecr-password.txt

REM Tag and push images
docker tag insavein/auth-service:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/insavein/auth-service:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/insavein/auth-service:latest
REM ... repeat for all services
```

**For Azure ACR:**
```cmd
REM Create ACR
az acr create --resource-group insavein-rg --name insaveinacr --sku Basic

REM Login to ACR
az acr login --name insaveinacr

REM Tag and push images
docker tag insavein/auth-service:latest insaveinacr.azurecr.io/insavein/auth-service:latest
docker push insaveinacr.azurecr.io/insavein/auth-service:latest
REM ... repeat for all services
```

**For Google GCR:**
```cmd
REM Configure Docker for GCR
gcloud auth configure-docker

REM Tag and push images
docker tag insavein/auth-service:latest gcr.io/<project-id>/insavein/auth-service:latest
docker push gcr.io/<project-id>/insavein/auth-service:latest
REM ... repeat for all services
```

### Step 2: Update Deployment Files

Update image references in all deployment YAML files:

```yaml
# Change from:
image: insavein/auth-service:latest

# To (AWS):
image: <account-id>.dkr.ecr.us-east-1.amazonaws.com/insavein/auth-service:latest

# Or (Azure):
image: insaveinacr.azurecr.io/insavein/auth-service:latest

# Or (GCP):
image: gcr.io/<project-id>/insavein/auth-service:latest
```

### Step 3: Setup External Database (Recommended for Production)

**For AWS RDS:**
```cmd
aws rds create-db-instance ^
  --db-instance-identifier insavein-db ^
  --db-instance-class db.t3.medium ^
  --engine postgres ^
  --engine-version 15.3 ^
  --master-username postgres ^
  --master-user-password <secure-password> ^
  --allocated-storage 100 ^
  --storage-type gp3 ^
  --backup-retention-period 7 ^
  --multi-az
```

**For Azure Database for PostgreSQL:**
```cmd
az postgres flexible-server create ^
  --resource-group insavein-rg ^
  --name insavein-db ^
  --location eastus ^
  --admin-user postgres ^
  --admin-password <secure-password> ^
  --sku-name Standard_D2s_v3 ^
  --tier GeneralPurpose ^
  --storage-size 128 ^
  --version 15
```

**For Google Cloud SQL:**
```cmd
gcloud sql instances create insavein-db ^
  --database-version=POSTGRES_15 ^
  --tier=db-n1-standard-2 ^
  --region=us-central1 ^
  --root-password=<secure-password> ^
  --storage-size=100GB ^
  --storage-type=SSD ^
  --backup
```

Update ConfigMap with external database endpoint.

### Step 4: Setup Ingress Controller

**Install NGINX Ingress Controller:**

```cmd
REM Add Helm repository
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

REM Install ingress controller
helm install nginx-ingress ingress-nginx/ingress-nginx ^
  --namespace ingress-nginx ^
  --create-namespace ^
  --set controller.service.type=LoadBalancer
```

**Wait for external IP:**
```cmd
kubectl get svc -n ingress-nginx
```

### Step 5: Setup TLS Certificates

**Install cert-manager:**
```cmd
REM Add Helm repository
helm repo add jetstack https://charts.jetstack.io
helm repo update

REM Install cert-manager
helm install cert-manager jetstack/cert-manager ^
  --namespace cert-manager ^
  --create-namespace ^
  --set installCRDs=true
```

**Apply certificate issuer:**
```cmd
kubectl apply -f k8s/cert-manager-issuer.yaml
```

### Step 6: Deploy Application

```cmd
cd k8s

REM Apply all configurations
kubectl apply -f namespace.yaml
kubectl apply -f priority-class.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
kubectl apply -f resource-quota.yaml
kubectl apply -f network-policy.yaml

REM Deploy services (skip postgres if using external DB)
kubectl apply -f auth-service-deployment.yaml
kubectl apply -f user-service-deployment.yaml
kubectl apply -f savings-service-deployment.yaml
kubectl apply -f budget-service-deployment.yaml
kubectl apply -f goal-service-deployment.yaml
kubectl apply -f education-service-deployment.yaml
kubectl apply -f notification-service-deployment.yaml
kubectl apply -f analytics-service-deployment.yaml
kubectl apply -f frontend-deployment.yaml

REM Deploy ingress
kubectl apply -f ingress.yaml
```

### Step 7: Configure DNS

Point your domain to the ingress controller's external IP:

```
api.insavein.com → <ingress-external-ip>
```

---

## Verification & Testing

### 1. Check Cluster Status

```cmd
REM View all resources
kubectl get all -n insavein

REM Check nodes
kubectl get nodes

REM Check resource usage
kubectl top nodes
kubectl top pods -n insavein
```

### 2. Check Deployments

```cmd
REM View deployments
kubectl get deployments -n insavein

REM Check deployment status
kubectl rollout status deployment/auth-service -n insavein

REM View pods
kubectl get pods -n insavein

REM Check pod logs
kubectl logs -f <pod-name> -n insavein
```

### 3. Test Service Health

```cmd
REM Port forward to test locally
kubectl port-forward -n insavein svc/auth-service 8080:8080

REM In another terminal, test health endpoint
curl http://localhost:8080/health
```

### 4. Test Database Connection

```cmd
REM Connect to PostgreSQL pod
kubectl exec -it postgres-0 -n insavein -- psql -U insavein_user -d insavein

REM Run test query (in psql shell)
REM SELECT version();
REM \dt
REM \q
```

### 5. Check Ingress

```cmd
REM View ingress
kubectl get ingress -n insavein

REM Describe ingress
kubectl describe ingress insavein-ingress -n insavein

REM Test ingress (if DNS configured)
curl https://api.insavein.com/api/auth/health
```

### 6. Monitor Logs

```cmd
REM View logs for all pods
kubectl logs -l app=auth-service -n insavein --tail=100

REM Follow logs
kubectl logs -f deployment/auth-service -n insavein

REM View logs from all containers
kubectl logs -l tier=backend -n insavein --all-containers=true
```

### 7. Check Resource Quotas

```cmd
REM View resource quota usage
kubectl describe resourcequota insavein-resource-quota -n insavein

REM View limit ranges
kubectl describe limitrange -n insavein
```

---

## Troubleshooting

### Issue: Pods Not Starting

**Check pod status:**
```cmd
kubectl get pods -n insavein
kubectl describe pod <pod-name> -n insavein
```

**Common causes:**
- Image pull errors (check image name and registry access)
- Resource limits exceeded (check resource quotas)
- Missing secrets or configmaps
- Node resource constraints

**Solutions:**
```cmd
REM Check events
kubectl get events -n insavein --sort-by=.lastTimestamp

REM Check resource usage
kubectl top nodes
kubectl describe resourcequota -n insavein

REM Check image pull
kubectl describe pod <pod-name> -n insavein | findstr /C:"Events:"
```

### Issue: Cannot Connect to Services

**Check service endpoints:**
```cmd
kubectl get svc -n insavein
kubectl get endpoints -n insavein
```

**Test connectivity:**
```cmd
REM Create debug pod
kubectl run debug --image=nicolaka/netshoot -n insavein -it --rm

REM Inside debug pod:
REM nslookup auth-service.insavein.svc.cluster.local
REM curl http://auth-service.insavein.svc.cluster.local:8080/health
```

**Check network policies:**
```cmd
kubectl get networkpolicies -n insavein
kubectl describe networkpolicy <policy-name> -n insavein
```

### Issue: Database Connection Failed

**Check PostgreSQL pod:**
```cmd
kubectl get pods -l app=postgres -n insavein
kubectl logs postgres-0 -n insavein
```

**Test database connection:**
```cmd
kubectl exec -it postgres-0 -n insavein -- psql -U postgres -c "SELECT 1"
```

**Check secrets:**
```cmd
kubectl get secret insavein-secrets -n insavein
kubectl describe secret insavein-secrets -n insavein
```

### Issue: Ingress Not Working

**Check ingress controller:**
```cmd
kubectl get pods -n ingress-nginx
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller
```

**Check ingress resource:**
```cmd
kubectl get ingress -n insavein
kubectl describe ingress insavein-ingress -n insavein
```

**Test without ingress:**
```cmd
kubectl port-forward -n insavein svc/auth-service 8080:8080
curl http://localhost:8080/health
```

### Issue: TLS Certificate Not Working

**Check cert-manager:**
```cmd
kubectl get pods -n cert-manager
kubectl logs -n cert-manager -l app=cert-manager
```

**Check certificate:**
```cmd
kubectl get certificate -n insavein
kubectl describe certificate insavein-tls -n insavein
```

**Check certificate request:**
```cmd
kubectl get certificaterequest -n insavein
kubectl describe certificaterequest <request-name> -n insavein
```

### Issue: High Resource Usage

**Check resource consumption:**
```cmd
kubectl top pods -n insavein --sort-by=memory
kubectl top pods -n insavein --sort-by=cpu
```

**Check HPA status:**
```cmd
kubectl get hpa -n insavein
kubectl describe hpa auth-service-hpa -n insavein
```

**Adjust resource limits:**
Edit deployment files and update resource requests/limits.

### Issue: Secrets Contain Placeholders

**Error:** "CHANGE_ME" values in secrets

**Solution:**
```cmd
REM Generate new secrets
cd k8s
openssl rand -base64 64
openssl rand -base64 32
openssl rand -base64 32

REM Update secrets.yaml
REM Delete and recreate secret
kubectl delete secret insavein-secrets -n insavein
kubectl apply -f secrets.yaml
```

### Issue: Minikube Not Starting

**Common solutions:**
```cmd
REM Delete and recreate cluster
minikube delete
minikube start --driver=hyperv --memory=8192 --cpus=4

REM Check Hyper-V is enabled (run as Administrator)
dism /Online /Get-FeatureInfo /FeatureName:Microsoft-Hyper-V

REM Enable Hyper-V if needed (requires restart)
dism /Online /Enable-Feature /All /FeatureName:Microsoft-Hyper-V
```

---

## Useful Commands Reference

### Cluster Management
```cmd
REM View cluster info
kubectl cluster-info

REM View nodes
kubectl get nodes

REM View all resources in namespace
kubectl get all -n insavein

REM View resource usage
kubectl top nodes
kubectl top pods -n insavein
```

### Pod Management
```cmd
REM List pods
kubectl get pods -n insavein

REM Describe pod
kubectl describe pod <pod-name> -n insavein

REM View logs
kubectl logs <pod-name> -n insavein
kubectl logs -f <pod-name> -n insavein

REM Execute command in pod
kubectl exec -it <pod-name> -n insavein -- /bin/sh

REM Delete pod (will be recreated by deployment)
kubectl delete pod <pod-name> -n insavein
```

### Deployment Management
```cmd
REM List deployments
kubectl get deployments -n insavein

REM Scale deployment
kubectl scale deployment auth-service --replicas=5 -n insavein

REM Update image
kubectl set image deployment/auth-service auth-service=insavein/auth-service:v2 -n insavein

REM Rollout status
kubectl rollout status deployment/auth-service -n insavein

REM Rollback deployment
kubectl rollout undo deployment/auth-service -n insavein

REM View rollout history
kubectl rollout history deployment/auth-service -n insavein
```

### Service Management
```cmd
REM List services
kubectl get svc -n insavein

REM Describe service
kubectl describe svc auth-service -n insavein

REM Port forward
kubectl port-forward svc/auth-service 8080:8080 -n insavein
```

### ConfigMap & Secrets
```cmd
REM View ConfigMaps
kubectl get configmap -n insavein

REM View ConfigMap data
kubectl get configmap insavein-config -n insavein -o yaml

REM View Secrets
kubectl get secrets -n insavein

REM Decode secret (requires base64 decoding tool)
kubectl get secret insavein-secrets -n insavein -o jsonpath="{.data.JWT_SECRET_KEY}"
```

### Debugging
```cmd
REM View events
kubectl get events -n insavein --sort-by=.lastTimestamp

REM Describe resource
kubectl describe <resource-type> <resource-name> -n insavein

REM Run debug pod
kubectl run debug --image=nicolaka/netshoot -n insavein -it --rm

REM Copy files from pod
kubectl cp <pod-name>:/path/to/file ./local-file -n insavein

REM Copy files to pod
kubectl cp ./local-file <pod-name>:/path/to/file -n insavein
```

---

## Production Checklist

Before deploying to production:

- [ ] Use external managed database (RDS, Cloud SQL, Azure Database)
- [ ] Configure automatic backups
- [ ] Set up monitoring (Prometheus, Grafana)
- [ ] Configure log aggregation (ELK, CloudWatch, Stackdriver)
- [ ] Enable TLS/SSL with valid certificates
- [ ] Configure DNS with your domain
- [ ] Set up CI/CD pipeline
- [ ] Configure resource limits and quotas
- [ ] Enable network policies
- [ ] Use external secret management (Vault, AWS Secrets Manager)
- [ ] Configure horizontal pod autoscaling
- [ ] Set up alerting (PagerDuty, Slack)
- [ ] Perform load testing
- [ ] Document disaster recovery procedures
- [ ] Configure backup and restore procedures
- [ ] Set up staging environment
- [ ] Implement blue-green or canary deployments
- [ ] Configure rate limiting
- [ ] Enable audit logging
- [ ] Perform security scanning
- [ ] Set up cost monitoring and alerts

---

## Additional Resources

- **Kubernetes Documentation**: https://kubernetes.io/docs/
- **kubectl Cheat Sheet**: https://kubernetes.io/docs/reference/kubectl/cheatsheet/
- **Docker Desktop Kubernetes**: https://docs.docker.com/desktop/kubernetes/
- **Minikube Documentation**: https://minikube.sigs.k8s.io/docs/
- **AWS EKS Documentation**: https://docs.aws.amazon.com/eks/
- **Azure AKS Documentation**: https://docs.microsoft.com/en-us/azure/aks/
- **Google GKE Documentation**: https://cloud.google.com/kubernetes-engine/docs
- **Helm Documentation**: https://helm.sh/docs/
- **NGINX Ingress Controller**: https://kubernetes.github.io/ingress-nginx/
- **cert-manager Documentation**: https://cert-manager.io/docs/

---

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review Kubernetes logs: `kubectl logs -n insavein <pod-name>`
3. Check cluster events: `kubectl get events -n insavein`
4. Review the main deployment guide: `k8s/DEPLOYMENT_GUIDE.md`
5. Contact your DevOps team

---

**Happy Deploying! 🚀**
