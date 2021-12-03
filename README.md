# kustomization-generator

## Usage helm

This generator allows you to convert a hosted helm chart into locally stored resource definitions.

1. Prepare folder with definition:
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
2. Run `kustomization-generator` in that folder
3. Find a ready to use kustomization with the content of that chart
    ```
    ├── kustomization-generator.yaml
    ├── kustomization.yaml
    └── generated
        └── crds
            └── ...
        └── templates
            └── ...
    ```

## Usage kustomize

This generator allows you to convert a remote kustomization into a locally stored resource definition.

1. Prepare folder with definition:
    ```yaml
    # kustomization-generator.yaml
    type: kustomize
    namespace: namespace
    url: github.com/CrunchyData/postgres-operator-examples/kustomize/install?ref=main
    args:
      - --reorder
      - legacy
    ```
2. Run `kustomization-generator` in that folder
3. Find a ready to use kustomization with the content of that chart
    ```
    ├── kustomization-generator.yaml
    ├── kustomization.yaml
    └── generated
        └── resources.yaml
    ```

## Usage download

TODO

## Installation

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/choffmeister/kustomization-generator:latest
docker run --rm -v $PWD:/workdir ghcr.io/choffmeister/kustomization-generator:latest
```
