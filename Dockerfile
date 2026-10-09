FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/document-extractor ./cmd/server

FROM alpine:3.22
RUN adduser -D -H app
USER app
COPY --from=build /out/document-extractor /document-extractor
EXPOSE 8080
ENTRYPOINT ["/document-extractor"]
