FROM alpine:3.15.2 as git-ghost-cli
RUN apk add --no-cache git ca-certificates openssh-client
COPY git-ghost /usr/local/bin/
