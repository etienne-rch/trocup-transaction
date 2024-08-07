FROM golang:1.18-alpine

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go.mod et go.sum et installer les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste de l'application
COPY . .

# Construire l'application
RUN go build -o app

# Exposer le port sur lequel l'application écoute
EXPOSE 5003

# Définir la commande par défaut pour exécuter l'application
CMD ["./app"]
