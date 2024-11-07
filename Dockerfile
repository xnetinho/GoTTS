# Etapa de construção do aplicativo Go
FROM golang:1.21-alpine AS builder

# Obter a arquitetura de destino
ARG TARGETARCH

# Definir as variáveis de ambiente para a compilação Go
ENV GOOS=linux
ENV GOARCH=$TARGETARCH

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api

# Iniciar uma nova etapa para a imagem final
FROM alpine:latest

# Obter a arquitetura de destino
ARG TARGETARCH

# Instalar dependências
RUN apk add --no-cache libstdc++ bash wget ca-certificates

# Definir o diretório de trabalho
WORKDIR /app

# Copiar o aplicativo Go compilado da etapa anterior
COPY --from=builder /app/main .

# Criar o diretório de vozes
RUN mkdir -p /app/voices
RUN chmod 755 /app/voices

# Determinar o binário apropriado do Piper com base na arquitetura
RUN \
    if [ "$TARGETARCH" = "amd64" ]; then \
        export PIPER_ARCH="amd64"; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
        export PIPER_ARCH="arm64"; \
    else \
        echo "Arquitetura não suportada: $TARGETARCH"; exit 1; \
    fi && \
    wget https://github.com/rhasspy/piper/releases/download/v1.2.0/piper_$PIPER_ARCH.tar.gz && \
    mkdir -p /tmp/piper_install && \
    tar xzf piper_$PIPER_ARCH.tar.gz -C /tmp/piper_install && \
    mv /tmp/piper_install/piper/piper /usr/local/bin/piper && \
    chmod +x /usr/local/bin/piper && \
    mv /tmp/piper_install/piper/lib* /usr/local/lib/ && \
    mv /tmp/piper_install/piper/espeak-ng-data /usr/local/share/ && \
    rm -rf /tmp/piper_install && \
    rm piper_$PIPER_ARCH.tar.gz

# Definir o LD_LIBRARY_PATH
ENV LD_LIBRARY_PATH="/usr/local/lib"

# Expor a porta da aplicação
EXPOSE 8080

# Comando para executar quando o contêiner iniciar
CMD ["./main"]