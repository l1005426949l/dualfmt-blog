# Stage 1: Build the application
# We use a specific Rust version for reproducibility.
FROM rust:1.78 as builder

# Install build dependencies. protobuf-compiler is needed for tonic-build.
RUN apt-get update && apt-get install -y protobuf-compiler

# Create a new empty project in the workspace to cache dependencies.
# This is faster than copying the whole project.
WORKDIR /usr/src/app
RUN cargo init

# Copy the dependency manifests
COPY typst-parser/Cargo.toml ./Cargo.toml

# Fetch and build the dependencies. This layer is cached as long as
# Cargo.toml does not change.
RUN cargo build --release

# Now copy the actual source code
COPY typst-parser/src ./src
COPY typst-parser/proto ./proto
COPY typst-parser/build.rs ./build.rs

# Build the application with the source code.
# This will be fast because dependencies are already cached.
RUN cargo build --release

# Stage 2: Create the runtime image
# Use a minimal image for a smaller footprint and better security.
FROM debian:slim as runtime

# Copy the compiled binary from the builder stage.
COPY --from=builder /usr/src/app/target/release/typst-parser /usr/local/bin/typst-parser

# Expose the gRPC port the service listens on.
EXPOSE 50051

# Set the command to run the binary.
CMD ["typst-parser"]
