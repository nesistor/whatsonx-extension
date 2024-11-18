# Wybieramy oficjalny obraz Go jako bazowy
FROM golang:1.20-alpine as builder

# Tworzymy katalog roboczy w kontenerze
WORKDIR /app

# Kopiujemy pliki go.mod i go.sum do kontenera, aby zainstalować zależności
COPY go.mod go.sum ./

# Instalujemy zależności
RUN go mod tidy

# Kopiujemy całą aplikację do kontenera
COPY . .

# Kompilujemy aplikację
RUN GOOS=linux GOARCH=amd64 go build -o app .

# Używamy lekkiego obrazu jako bazowego do uruchomienia aplikacji
FROM alpine:latest

# Instalujemy zależności wymagane do uruchomienia aplikacji (np. biblioteki PostgreSQL)
RUN apk --no-cache add ca-certificates

# Ustawiamy zmienną środowiskową na porcie 80
ENV PORT 80

# Tworzymy katalog do przechowywania aplikacji w kontenerze
WORKDIR /app

# Kopiujemy aplikację z obrazu "builder" do nowego obrazu
COPY --from=builder /app/app .

# Ustawiamy punkt wejścia do uruchamiania aplikacji
ENTRYPOINT ["./app"]

# Otwieramy port 80
EXPOSE 80

