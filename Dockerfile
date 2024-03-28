FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN apk --no-cache add ca-certificates
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gotest


FROM scratch
COPY --from=builder /app/gotest /gotest
ENV REDIS_URL=$REDIS_URL
ENV DB_URL=$DB_URL
# ARG DB_URL
# ARG REDIS_URL
ENTRYPOINT ["/gotest"]