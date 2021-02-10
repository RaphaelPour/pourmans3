FROM alpine:latest

# Copy static website content to caddy's www directory
COPY ci-build/pourmans3 /usr/bin/pourmans3
ENTRYPOINT [ "/usr/bin/pourmans3", "--port=8004" ]
