FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ ./pkg/
COPY services/ingestion/ ./services/ingestion/

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ingestion ./services/ingestion/cmd

FROM alpine:3.20
COPY --from=build /out/ingestion /ingestion
CMD ["/ingestion"]