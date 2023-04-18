# Usamos la imagen oficial de Go como base
FROM golang:latest

# Establecemos el directorio de trabajo en el GOPATH
WORKDIR /go/src/app

# Copiamos los archivos del proyecto al contenedor
COPY . .

# Instalamos las dependencias del proyecto
RUN go get -d -v ./...

# Compilamos la aplicación Go
RUN go install -v ./...

# Exponemos el puerto en el que la aplicación escucha
EXPOSE 8080

# Ejecutamos la aplicación al iniciar el contenedor
CMD ["app"]
