FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
  wget \
  gnupg \
  lsb-release \
  software-properties-common \
  curl \
  build-essential \
  cmake \
  ninja-build \
  git \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /tmp/* \
  && rm -rf /var/log/*

RUN curl -OL https://apt.llvm.org/llvm.sh && \
  chmod +x llvm.sh && \
  ./llvm.sh 19 all && \
  apt-get install -y --no-install-recommends clang-19 llvm-19 lldb-19 lld-19 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /tmp/* \
  && rm -rf /var/log/*

ENV GOLANG_VERSION=1.22.4
ENV PATH="/usr/local/go/bin:${PATH}"

RUN curl -OL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && rm go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && rm -rf /tmp/* \
  && rm -rf /var/log/*

WORKDIR /tmp/go-build

COPY lang/go.mod .
COPY lang/go.sum .

RUN go mod download && go mod tidy
