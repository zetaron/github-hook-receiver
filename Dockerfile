FROM alpine:3.4

EXPOSE 80
ENTRYPOINT ["/usr/bin/secret-wrapper", "/usr/bin/github-hook-receiver"]
LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/zetaron/github-hook-receiver" \
      org.label-schema.vcs-url="https://github.com/zetaron/github-hook-receiver" \
      org.label-schema.name="github-hook-receiver" \
      org.label-schema.docker.cmd="docker run -d -p 80 -v ${SECRETS_VOLUME_NAME:-github-hook-receiver-secrets}:/var/cache/secrets --name github-hook-receiver zetaron/github-hook-receiver"

COPY secret-wrapper /usr/bin/secret-wrapper
COPY github-hook-receiver /usr/bin/github-hook-receiver
