# frontend builder
FROM node:20-alpine as web

WORKDIR /app

RUN apk add --no-cache git
COPY web/package*.json ./

RUN npm ci

COPY web/ .
RUN npm run build


# go builder
FROM golang:1.26 as build

WORKDIR /make/

COPY --from=web /app/dist /make/web/dist
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /makedotcsh


# final image
# i love binaries <3 37.1MB docker image
FROM gcr.io/distroless/static:nonroot

COPY --from=build /makedotcsh /makedotcsh

EXPOSE 8080
CMD ["/makedotcsh"]

