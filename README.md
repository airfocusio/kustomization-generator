# kustomization-generator

This tool allows to fully embed external kubernetes resources like Helm Charts, raw URLs or remote Kustomizations into your repository. The general usage is:

1. Prepare a folder with an configuration (see below), for example `vendors/cert-manager/kustomization-generator.yaml`
2. Run `kustomization-generator --dir=vendors/cert-manager`
3. Find a ready to use Kustomization locally stored:
    ```
    ├── vendors
    │   ├── cert-manager
    │   │   ├── kustomization-generator.yaml
    │   │   ├── kustomization.yaml
    │   │   ├── crds
    │   │   │   ├── kustomization.yaml
    │   │   │   └── ...
    │   │   ├── namespaces
    │   │   │   ├── kustomization.yaml
    │   │   │   └── ...
    │   │   └── resources
    │   │   │   ├── kustomization.yaml
    │   │       └── ...
    │   └── ...
    └── ...
    ```

## Usage helm

This generator allows you to convert a hosted helm chart into locally stored resource definitions.

```yaml
# kustomization-generator.yaml
type: helm
registry: https://charts.jetstack.io
chart: cert-manager
version: v1.6.1
name: cert-manager
namespace: cert-manager-system
args:
  - --include-crds
values:
  some: value
```

## Usage kustomize

This generator allows you to convert a remote kustomization into a locally stored resource definitions.

```yaml
# kustomization-generator.yaml
type: kustomize
url: github.com/CrunchyData/postgres-operator-examples/kustomize/install?ref=main
args:
  - --reorder
  - legacy
```

## Usage download

```yaml
# kustomization-generator.yaml
type: download
url: https://raw.githubusercontent.com/longhorn/longhorn/v1.2.2/deploy/longhorn.yaml
```

## Installation

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/choffmeister/kustomization-generator:latest
docker run --rm -v $PWD:/workdir ghcr.io/choffmeister/kustomization-generator:latest
```
