FROM golang:latest as builder

COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -v -o /build/url_shortner ./cmd/url_shortner

FROM scratch

LABEL authors="Ilya Kharev"
COPY --from=builder /build /
ENTRYPOINT ["./url_shortner"]
