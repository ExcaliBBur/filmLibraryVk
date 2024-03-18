FROM golang:1.21

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o film-library ./cmd/main.go

CMD ["./film-library"]