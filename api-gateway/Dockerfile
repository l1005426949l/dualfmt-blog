# --- Build Stage ---
# Use the official Golang image as a builder.
# Using alpine for a smaller image size.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for creating a static binary
# -o /app/server specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# --- Final Stage ---
# Use a minimal, non-root base image for the final container.
# scratch is the most minimal, containing nothing.
FROM scratch

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server /server

# Copy the configuration file(s)
COPY --from=builder /app/configs /configs

# Expose the port the app runs on
# The default HTTP port is 8000 and gRPC is 9000 from our config
EXPOSE 8000
EXPOSE 9000

# The command to run the application
ENTRYPOINT ["/server"]
