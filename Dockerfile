
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/soul-soldiers/image-resizer
COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -v -o image-resizer
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/src/github.com/soul-soldiers/image-resizer/image-resizer /image-resizer
# Run the resize-image binary.
ENTRYPOINT ["/image-resizer"]