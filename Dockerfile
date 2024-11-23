# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine
# FROM gcr.io/distroless/static:nonroot
COPY gss /gss
COPY --from=nginx:stable-alpine --chown=65534:65534 /usr/share/nginx/html /site
# nobody
WORKDIR /site
USER 65534:65534
ENTRYPOINT ["/gss"]
