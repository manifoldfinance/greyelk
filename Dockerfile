FROM golang:1.13-alpine AS builder

COPY . /go/src/github.com/manifoldfinance/greyelk

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates bash git openssh build-base && \
    cd /go/src/github.com/manifoldfinance/greyelk && \
    make

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=builder /go/src/github.com/manifoldfinance/greyelk/greyelk* ./

ENTRYPOINT [ "/root/greyelk" ]
CMD [ "-c", "/root/greyelk.yml", "-e"]