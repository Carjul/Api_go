package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Carjul/GOLAN_API/db"
	"github.com/gorilla/mux"
)

// Estructura para el modelo de Usuario
type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}
type Id struct {
	ID int `json:"id"`
}
type Rol struct {
	Rol string `json:"rol"`
}

func CrearRol(w http.ResponseWriter, r *http.Request) {
	// Obtener el rol del cuerpo de la solicitud
	var rol Rol
	err := json.NewDecoder(r.Body).Decode(&rol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insertar el rol en la base de datos
	row := db.DB.QueryRow("INSERT INTO Rols (rol) VALUES ($1) RETURNING id", rol.Rol)

	// Extraer el ID del rol insertado
	var id int
	err = row.Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")

	// Escribir el cuerpo de la respuesta HTTP con el ID del rol insertado
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bienvenido a la API de GOLANG"))
	return
}

func ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Consultar el usuario en la base de datos
	row := db.DB.QueryRow("SELECT id, nombre, rol FROM usuarios WHERE id = $1", id)

	// Crear una variable para almacenar los datos del usuario
	var usuario Usuario

	// Extraer los datos del usuario de la fila
	err := row.Scan(&usuario.ID, &usuario.Nombre, &usuario.Rol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir los datos del usuario a formato JSON
	usuarioJSON, err := json.Marshal(usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")

	// Escribir el cuerpo de la respuesta HTTP con los datos en formato JSON
	w.Write(usuarioJSON)
}
func ObtenerUsuarios(w http.ResponseWriter, r *http.Request) {
	// Consultar todos los usuarios en la base de datos
	rows, err := db.DB.Query("SELECT id, nombre, rol FROM usuarios")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Crear una lista de usuarios
	usuarios := []Usuario{}
	for rows.Next() {
		var usuario Usuario
		err := rows.Scan(&usuario.ID, &usuario.Nombre, &usuario.Rol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usuarios = append(usuarios, usuario)
	}

	// Convertir la lista de usuarios a JSON
	usuariosJSON, err := json.Marshal(usuarios)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")

	// Escribir el cuerpo de la respuesta HTTP con los datos en formato JSON
	w.Write(usuariosJSON)
}

func CrearUsuario(w http.ResponseWriter, r *http.Request) {
	// Leer los datos del usuario del cuerpo de la petición
	var usuario Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insertar el usuario en la base de datos
	_, err = db.DB.Exec("INSERT INTO Usuarios (id, nombre, rol) VALUES ($1, $2, $3)", usuario.ID, usuario.Nombre, usuario.Rol)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Escribir el cuerpo de la respuesta HTTP con la confirmación de creación del usuario
	fmt.Fprint(w, `{"mensaje": "Usuario creado correctamente"}`)
}

func EliminarUsuario(w http.ResponseWriter, r *http.Request) {

	var userid Id
	err := json.NewDecoder(r.Body).Decode(&userid)
	// Eliminar el usuario de la base de datos

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("DELETE FROM Usuarios WHERE id = $1", userid.ID)
	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Escribir el cuerpo de la respuesta HTTP con la confirmación de creación del usuario
	fmt.Fprint(w, `{"mensaje": "Usuario eliminado correctamente"}`)
}

func ActualizarUser(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insertar el usuario en la base de datos
	_, err = db.DB.Exec("UPDATE Usuarios SET nombre = $1, rol = $2 WHERE id = $3", usuario.Nombre, usuario.Rol, usuario.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Escribir el cuerpo de la respuesta HTTP con la confirmación de creación del usuario
	fmt.Fprint(w, `{"mensaje": "Usuario actualizado correctamente"}`)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el archivo del cuerpo de la solicitud
	file, handler, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crear un nuevo archivo en el sistema de archivos
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copiar el contenido del archivo recibido al archivo creado en el sistema de archivos
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con una confirmación de carga exitosa
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Archivo cargado exitosamente")
}
