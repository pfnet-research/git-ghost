FROM ubuntu:16.04

RUN apt-get update -q && apt-get install -yq --no-install-recommends git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY . /code/
WORKDIR /code
