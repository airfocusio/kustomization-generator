# kustomization-helm

## Usage

1. Prepare folder with definition:
    ```yaml
    # kustomization-helm.yaml
    registry: https://charts.jetstack.io
    chart: cert-manager
    version: v1.6.1
    name: cert-manager
    namespace: cert-manager-system
    args:
      - --set
      - installCRDs=true
    values:
      webhook:
        timeoutSeconds: 4
    ```
2. Run `kustomization-helm` in that folder
3. Find a ready to use kustomization with the content of that chart
    ```
    ├── kustomization-helm.yaml
    ├── kustomization.yaml
    └── templates
        ├── cainjector-deployment.yaml
        ├── ...
        └── webhook-validating-webhook.yaml
    ```



## Installation

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/choffmeister/kustomization-helm:latest
docker run --rm -v $PWD:/workdir ghcr.io/choffmeister/kustomization-helm:latest
```
