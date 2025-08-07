# Stage 1: Build the application
FROM rust:1.78 as builder

# Install protobuf compiler
RUN apt-get update && apt-get install -y protobuf-compiler

WORKDIR /usr/src/typst-parser

# Copy the project files
COPY . .

# Build the project
RUN cargo build --release

# Stage 2: Create the runtime image
FROM debian:slim as runtime

# Copy the compiled binary from the builder stage
COPY --from=builder /usr/src/typst-parser/target/release/typst-parser /usr/local/bin/typst-parser

# Expose the gRPC port
EXPOSE 50051

# Set the command to run the binary
CMD ["typst-parser"]
