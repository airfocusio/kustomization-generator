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
        ├── cert-manager-cainjector-clusterrole.yaml
        └── ...
    ```

## Usage kustomize

This generator allows you to convert a remote kustomization into a locally stored resource definitions.

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
        ├── pgo-deployment.yaml
        └── ...
    ```

## Usage download

This generator allows you to convert a remote resource file into a locally stored resource definitions.

1. Prepare folder with definition:
    ```yaml
    # kustomization-generator.yaml
    type: download
    namespace: longhorn-system
    url: https://raw.githubusercontent.com/longhorn/longhorn/v1.2.2/deploy/longhorn.yaml
    ```
2. Run `kustomization-generator` in that folder
3. Find a ready to use kustomization with the content of that chart
    ```
    ├── kustomization-generator.yaml
    ├── kustomization.yaml
    └── generated
        ├── backingimagedatasources-longhorn-io-customresourcedefinition.yaml
        └── ...
    ```

## Installation

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/choffmeister/kustomization-generator:latest
docker run --rm -v $PWD:/workdir ghcr.io/choffmeister/kustomization-generator:latest
```
