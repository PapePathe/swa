FROM ghcr.io/papepathe/swa-base-container-builder:master

WORKDIR /

COPY . .

RUN go build -o swa main.go && rm -rf ./*.go

CMD ["./swa"]
