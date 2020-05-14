FROM ubuntu:bionic

ARG version

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates

RUN \
  # Download
  curl -o  /tmp/odyssey.tar.gz -L https://github.com/yandex/odyssey/releases/download/$version/odyssey.linux-amd64.b7bcb86.tar.gz && \
  cd /tmp && \
  # Unpack
  tar xvfz /tmp/odyssey.tar.gz && \
  # Manual install
  cp odyssey /usr/bin && \
  mkdir -p /etc/odyssey && \
  # Cleanup
  rm -rf /tmp/odyssey*

CMD ["sh", "-c", "/usr/bin/odyssey /etc/odyssey/odyssey.conf"]