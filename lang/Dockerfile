FROM golang:1.23.2-bookworm as builder
WORKDIR /swahili
COPY . .
ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all" \
    GOFLAGS="-buildmode=pie"
RUN go build -ldflags "-s -w" -trimpath -o swa main.go

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /swahili/swa /bin/swa
USER 65534
ENTRYPOINT ["/bin/swa"]
