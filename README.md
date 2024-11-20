## Description
**kubectl multiclone** is a **kubectl** plugin that empowers users to clone Kubernetes resources across multiple clusters and projects managed by Rancher. It simplifies the process of duplicating resources from one cluster to another or within different namespaces and projects, with optional on-the-fly modifications. This tool enhances multi-cluster resource management, making it invaluable for environments where Rancher orchestrates numerous Kubernetes clusters.

## Goals
1. **Seamless Multi-Cluster Cloning**
   - Clone Kubernetes resources across clusters/projects with one command.
   - Simplifies management, reduces operational effort.

## Resources
1. **Rancher & Kubernetes Docs**
   - Rancher API, Cluster Management, Kubernetes client libraries.

2. **Development Tools**
   - Kubectl plugin docs, Go programming resources.

### **Building and Installing the Plugin**

1. **Set Environment Variables:**
   Export the Rancher URL and API token:

   - `export RANCHER_URL="https://rancher.example.com"`
   - `export RANCHER_TOKEN="token-xxxxx:xxxxxxxxxxxxxxxxxxxx"`


2. **Build the Plugin:**
   Compile the Go program:

   - `go build -o kubectl-multiclone main.go`

3. **Install the Plugin:**
   Move the executable to a directory in your `PATH`:
   - `mv kubectl-multiclone /usr/local/bin/`

   Ensure the file is executable:
   - `chmod +x /usr/local/bin/kubectl-multiclone`


4. **Verify the Plugin Installation:**
   Test the plugin by running:

   - `kubectl multiclone --help`

   You should see the usage information for the `kubectl-multiclone` plugin.

### **Usage Examples**
1. **Clone a Deployment from One Cluster to Another:**
   - `kubectl multiclone --source-cluster c-abc123 --type deployment --name nginx-deployment --target-cluster c-def456 --new-name nginx-deployment-clone`


2. **Clone a Service into Another Namespace and Modify Labels:**
   - `kubectl multiclone --source-cluster c-abc123 --type service --name my-service --source-namespace default --target-cluster c-def456 --target-namespace staging --modify "metadata.labels.env=staging"`


3. **Clone a ConfigMap within the Same Cluster but Different Project:**
   - `kubectl multiclone --source-cluster c-abc123 --source-project p-abc123 --type configmap --name my-config --target-cluster c-abc123 --target-project p-def456 --target-namespace dev`

4. **Clone a Secret with a New Name and Modifications:**
   - `kubectl multiclone --source-cluster c-abc123  --type secret --name my-secret --target-cluster c-def456 --new-name my-secret-copy --modify "metadata.annotations.description=Cloned Secret"`
