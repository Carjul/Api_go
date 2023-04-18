package db

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

// PGPASSWORD=wZ8qpAwocCkl4twZSEAx psql -h containers-us-west-32.railway.app -U postgres -p 7429 -d railway
func ConexionDB() {
	var err error
	DB, err = sql.Open("postgres", "host=containers-us-west-32.railway.app port=7429 dbname=railway user=postgres password=wZ8qpAwocCkl4twZSEAx sslmode=disable")
	if err != nil {
		log.Fatal(err)
	} else {
		// Ejecutar la consulta SQL para crear la tabla de rols
		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS Rols (id SERIAL PRIMARY KEY, rol VARCHAR(50))")
		if err != nil {
			fmt.Println("Error al crear la tabla de Rols:", err)
			return
		}
		// Ejecutar la consulta SQL para crear la tabla de usuarios
		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS usuarios (id SERIAL PRIMARY KEY, nombre VARCHAR(50), rol INT REFERENCES Rols(id))")
		if err != nil {
			fmt.Println("Error al crear la tabla de usuarios:", err)
			return
		}

		fmt.Println("db Conectada")
	}

}
