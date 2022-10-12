FROM debian:bullseye-slim
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    openssh-client \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
COPY git-ghost /usr/local/bin/
