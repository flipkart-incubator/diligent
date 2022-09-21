FROM golang:1.18-alpine as diligent-builder
RUN apk add --no-cache make
WORKDIR /diligent
COPY go.mod go.sum Makefile ./
RUN go mod download
COPY apps ./apps
COPY pkg ./pkg
RUN make boss

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /diligent
COPY --from=diligent-builder /diligent/build/boss ./
CMD ["/diligent/boss"]

