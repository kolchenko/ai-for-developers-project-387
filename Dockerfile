# syntax=docker/dockerfile:1

FROM node:22-alpine AS frontend
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM golang:1.22-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/callcalendar ./cmd/server

FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=backend /out/callcalendar /app/callcalendar
COPY --from=frontend /app/dist /app/web
RUN chown -R app:app /app
USER app
ENV PORT=8080
ENV WEB_DIR=/app/web
EXPOSE 8080
CMD ["/app/callcalendar"]
