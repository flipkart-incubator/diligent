FROM golang:1.16-alpine as diligent-builder
WORKDIR /diligent
COPY go.mod go.sum ./
RUN go mod download
COPY apps ./apps
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o boss github.com/flipkart-incubator/diligent/apps/boss

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /diligent
COPY --from=diligent-builder /diligent/boss ./
CMD ["/diligent/boss"]

