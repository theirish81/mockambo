FROM golang:1.22.5-bullseye AS builder

RUN mkdir /usr/local/mockambo
WORKDIR /usr/local/mockambo

COPY . .

RUN go get
RUN go build -o mockambo *.go

FROM debian:bullseye
RUN mkdir /usr/local/mockambo
WORKDIR /usr/local/mockambo
COPY --from=builder /usr/local/mockambo/mockambo .

RUN addgroup --gid 1000 mockambo && \
    adduser --home /usr/local/mockambo -u 1000 --gid 1000 mockambo && \
    chown -R mockambo:mockambo /usr/local/mockambo

USER mockambo
WORKDIR /usr/local/mockambo
ENTRYPOINT [ "/usr/local/mockambo/mockambo" ]