FROM golang:1.18 as builder

# Set Environment Variables
ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN go build -o app ./cmd/api/

FROM alpine:latest as service

COPY --from=builder /app/app .
COPY --from=builder /app/static static/.
COPY --from=builder /app/openapi openapi/.
COPY --from=builder /app/configs configs/.
COPY --from=builder /app/data/data_dump.csv data/data_dump.csv

#EXPOSE 4000
ENTRYPOINT ["./app"]