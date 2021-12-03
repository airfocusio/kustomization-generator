FROM alpine:latest as certs
RUN apk add --update --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bin/kustomization-helm"]
COPY kustomization-helm /bin/kustomization-helm
WORKDIR /workdir
