# Use an official Node image as a base for building the UI
FROM node:18 AS ui-build

# Set working directory for UI
WORKDIR /app/ui

# Copy only the package.json and yarn.lock first to leverage Docker caching
COPY ./ui/package.json ./ui/yarn.lock ./

# Install UI dependencies
RUN yarn install

# Now copy the rest of the UI files
COPY ./ui/ ./

# Build the UI
RUN yarn build

# Use an official Golang image for the Go backend
FROM golang:1.23-alpine AS go-build

# Set working directory for Go
WORKDIR /app

# Copy Go project files
COPY . .

# Copy the UI build from the previous stage
COPY --from=ui-build /app/ui/build /app/ui/build

# Build the Go app (assuming there's a Go module in place)
RUN go build -o /app/server ./cmd/main

# Final stage: Use Python as the base image
FROM python:3.12.5

# Set working directory
WORKDIR /app

# Copy the Go binary from the go-build stage
COPY --from=go-build /app/server /app/server

# Copy the UI build files (if needed at runtime by your Go app)
COPY --from=go-build /app/ui/build /app/ui/build

# Expose the port the Go app will run on
EXPOSE 3000

# Install any Python dependencies (if needed)
# RUN pip install -r requirements.txt  # Uncomment if you have Python dependencies

# Run the Go app
CMD ["/app/server"]
