FROM ubuntu:22.04
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    openssh-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY git-ghost /usr/local/bin
ENTRYPOINT ["git", "ghost"]
