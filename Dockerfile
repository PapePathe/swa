FROM silkeh/clang

RUN apt-get update && apt-get install -y curl

# Install Go
ENV GOLANG_VERSION=1.22.4
RUN curl -OL https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz \
  && rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /tmp/go-build
COPY lang/go.mod .
COPY lang/go.sum .

RUN go mod tidy 
