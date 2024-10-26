# STAGE 1
FROM golang:1.23 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

# STAGE 2
FROM scratch AS final
COPY --from=build /app /app
COPY migrations migrations

ENTRYPOINT ["/app"]
