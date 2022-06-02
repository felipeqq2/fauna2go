FROM golang:1.18.3-buster AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .

RUN go build -o /fauna2go


FROM fauna/faunadb:4.13.0
WORKDIR /faunadb

COPY --from=build /fauna2go /fauna2go

EXPOSE 1000
EXPOSE 8443
EXPOSE 8084
EXPOSE 8085

CMD faunadb & /fauna2go