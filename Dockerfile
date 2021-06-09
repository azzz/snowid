FROM golang:1.16.3 as build

WORKDIR /build

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . .
RUN mkdir /build/bin
RUN go build -ldflags="-s -w" -o /build/bin/snowid ./cmd/snowid

FROM scratch

COPY --from=build /build/bin/snowid /app/snowid

ENV MACHINE_ID=1
ENV LISTEN=":8080"
ENV LOG_LEVEL=info
ENV EPOCH="20210413001805"

ENTRYPOINT ["/app/snowid"]
