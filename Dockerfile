
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/soul-soldiers/image-resizer
COPY . .

RUN go get -d -v

RUN go build -o /go/bin/resize-image
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/resize-image /go/bin/resize-image
# Run the resize-image binary.
ENTRYPOINT ["/go/bin/resize-image"]