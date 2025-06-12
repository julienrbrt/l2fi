FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache \
    build-base \
    go \
    npm \
    make

# Copy source code
COPY . .

# Build the application
RUN make build

EXPOSE 8080

CMD ["./l2fi"]
