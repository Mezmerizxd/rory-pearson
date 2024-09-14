# Use an official Node image as a base
FROM node:18 AS ui-build

# Install Go
RUN apt-get update && apt-get install -y golang-go

# Set working directory for UI
WORKDIR /app/ui

# Copy UI project files
COPY ./ui/ /app/ui/

# Use a minimal base image for final deployment
FROM golang:1.23-alpine AS go-build

# Copy UI build artifacts from the previous stage
COPY --from=ui-build /app/ui/build /app/ui/build

# Set working directory for Go
WORKDIR /app

# Copy Go project files
COPY . .

# Run your Go app (assuming main.go exists)
CMD ["go", "run", "cmd/main/main.go"]

# Expose port 3000
EXPOSE 3000