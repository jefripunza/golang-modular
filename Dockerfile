FROM golang:1.22 as builder 
MAINTAINER Jefri Herdi Triyanto, jefriherditriyanto@gmail.com

# Setup Environment
ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64
ENV CGO 0

# Install dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    musl-dev \
    wget

WORKDIR /build
COPY . .

# Install Go modules
# RUN go mod download
RUN go mod tidy

# Build the application
RUN go build -o ./run

# Configure Environment (change target env)
# RUN sed -i 's#localhost#host.docker.internal#g' .env

# Finishing
FROM ubuntu:latest as runner
WORKDIR /app

# Add the community repository to get ffmpeg
RUN apt-get update && apt-get install -y \
    openssl curl nano ffmpeg \
    xvfb libfontconfig wkhtmltopdf

# Copy the built application and environment configuration
# COPY --from=builder /build/.env /app/.env
COPY --from=builder /build/run /app/run
COPY --from=builder /build/views /app/views

RUN chmod +x /app/run

ENTRYPOINT ["/app/run"]
CMD ["run"]
