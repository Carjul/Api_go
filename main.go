package main

import (
	"log"
	"net/http"

	"github.com/Carjul/GOLAN_API/db"
	"github.com/Carjul/GOLAN_API/routes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Configurar la conexión a la base de datos PostgreSQL

	db.ConexionDB()

	// Configurar el enrutador Mux
	r := mux.NewRouter()

	// Rutas de API
	r.HandleFunc("/", routes.IndexRoute)
	r.HandleFunc("/usuarios", routes.ObtenerUsuarios).Methods("GET")
	r.HandleFunc("/usuarios", routes.CrearUsuario).Methods("POST")
	r.HandleFunc("/usuarios", routes.EliminarUsuario).Methods("DELETE")
	r.HandleFunc("/usuarios", routes.ActualizarUser).Methods("PUT")
	r.HandleFunc("/usuarios/{id}", routes.ObtenerUsuario).Methods("GET")
	r.HandleFunc("/roles", routes.CrearRol).Methods("POST")
	r.HandleFunc("/upload", routes.UploadHandler).Methods("POST")
	// Iniciar el servidor HTTP
	log.Println("Servidor escuchando en el puerto :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
