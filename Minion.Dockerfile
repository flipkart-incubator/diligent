FROM golang:1.16-alpine as diligent-builder
WORKDIR /diligent
COPY go.mod go.sum ./
RUN go mod download
COPY apps ./apps
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o minion github.com/flipkart-incubator/diligent/apps/minion

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /diligent
COPY --from=diligent-builder /diligent/minion ./
CMD ["/diligent/minion"]

