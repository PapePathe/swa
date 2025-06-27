FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
  wget \
  gnupg \
  lsb-release \
  software-properties-common \
  curl \
  build-essential \
  cmake \
  ninja-build
RUN wget https://apt.llvm.org/llvm.sh && \
  chmod +x llvm.sh && \
  ./llvm.sh 19 all
RUN apt-get install -y clang-19 lldb-19 lld-19
ENV GOLANG_VERSION=1.22.4
RUN curl -OL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
WORKDIR /tmp/go-build
COPY . .
RUN go mod download
RUN go mod tidy
