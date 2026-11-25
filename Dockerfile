FROM golang:1.26

WORKDIR /app

ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o forum ./cmd

EXPOSE 8080

CMD ["./forum"]
