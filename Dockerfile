FROM golang AS build-env

WORKDIR /src

COPY go.mod go.sum /src/

RUN go mod download

COPY . /src

RUN go build -v -o servirus

FROM alpine

WORKDIR /app

COPY --from=build-env /src/servirus /app/

RUN chmod +x ./servirus && chown root:root ./servirus

CMD ["./servirus"]

EXPOSE 50051
