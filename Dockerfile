FROM golang:1.18 AS deps

WORKDIR /go/src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

FROM deps as build

COPY cmd cmd
COPY pkg pkg

RUN go build -o /go/bin/app cmd/sellerpayout/sellerpayout.go

FROM gcr.io/distroless/base-debian10

COPY --from=build /go/bin/app /

ENTRYPOINT ["/app"]