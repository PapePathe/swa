FROM ghcr.io/papepathe/swa-base-container-builder:master as builder

WORKDIR /

COPY . .

RUN go build -o swa main.go && rm -rf ./*.go


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
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /tmp/* \
  && rm -rf /var/log/*

RUN curl -OL https://apt.llvm.org/llvm.sh && \
  chmod +x llvm.sh && \
  ./llvm.sh 18 all && \
  apt-get install -y --no-install-recommends clang-18 llvm-18 lldb-18 lld-18 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && rm -rf /tmp/* \
  && rm -rf /var/log/*

COPY --from=builder /swa /
COPY --from=builder /DockerfileEntrypoint.sh /

ENTRYPOINT ["/DockerfileEntrypoint.sh"]
