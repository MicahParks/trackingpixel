FROM golang:1.15 AS builder

# Get the Golang dependencies for better caching.
WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download

# Copy the code in.
COPY . .

# Build the code.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o trackingpixel


# The actual image being produced.
FROM scratch

COPY --from=builder /app/trackingpixel /trackingpixel
CMD ["/trackingpixel"]
