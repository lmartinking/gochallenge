FROM golang:1.11.1 as build

WORKDIR /go/src/gochallenge
COPY *.go ./

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gochallenge -v .

FROM scratch
WORKDIR /
COPY etc .
COPY --from=build /go/src/gochallenge/gochallenge .
EXPOSE 8000
CMD ["./gochallenge", "-port", "8000", "-pubkey", "key.pub", "-privkey", "key.priv"]