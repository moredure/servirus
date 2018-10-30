FROM golang AS builder

RUN mkdir /src

COPY go.mod go.sum /src/

WORKDIR /src

RUN go mod download

COPY . /src/

RUN CGO_ENABLED=0 GOOS=linux go build -o servirus

FROM alpine

RUN mkdir /src

COPY --from=builder /src/servirus /src/

CMD ["/src/servirus"]

EXPOSE 50051
