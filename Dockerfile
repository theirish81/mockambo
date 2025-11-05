FROM golang:1.24.8-trixie AS builder

RUN mkdir /usr/local/mockambo
WORKDIR /usr/local/mockambo

COPY . .

RUN go get
RUN go build -o mockambo *.go

FROM debian:trixie
RUN mkdir /usr/local/mockambo
WORKDIR /usr/local/mockambo
COPY --from=builder /usr/local/mockambo/mockambo .

RUN groupadd --gid 1000 mockambo && \
    useradd --create-home --home /usr/local/mockambo --uid 1000 --gid 1000 --shell /bin/bash mockambo

USER mockambo
WORKDIR /usr/local/mockambo
ENTRYPOINT [ "/usr/local/mockambo/mockambo" ]