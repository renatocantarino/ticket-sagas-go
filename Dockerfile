# Etapa 1: Build
FROM golang:1.22-alpine AS builder

# Configurações para binário estático
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Diretório de trabalho
WORKDIR /src

# Copia go.mod e go.sum primeiro (cache eficiente)
COPY go.mod go.sum* ./
RUN go mod download

# Copia todo o código
COPY . .

# Compila o binário
# Ajuste conforme a localização do seu main.go

RUN go build -o ./bin/app ./cmd/main.go

# Etapa 2: Imagem final mínima
FROM gcr.io/distroless/static-debian12 AS final

# Copia o binário compilado
COPY --from=builder /src/bin/app /app

# Porta (ajuste se seu app usa outra)
EXPOSE 8080

# Executa o app
ENTRYPOINT ["/app"]