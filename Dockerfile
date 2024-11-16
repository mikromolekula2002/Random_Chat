# Stage 1: Build
FROM golang:1.23-alpine AS build
WORKDIR /app
RUN apk add --no-cache make
COPY . .
RUN go mod download
RUN make build

# Stage 2: Production
FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/cmd/random-chat/main ./main
COPY --from=build /app/.env ./.env
EXPOSE 7994
CMD ["./main"]