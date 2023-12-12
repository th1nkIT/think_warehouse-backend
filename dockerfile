FROM golang:1.17.6-alpine

WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
# RUN go mod tidy
COPY . .
# Put the binary inside root path
# cmd/bin used for manual build outside container only
RUN go build -o main cmd/api/application.go
EXPOSE 6969

# Let's rock :)
CMD ["./main"]
