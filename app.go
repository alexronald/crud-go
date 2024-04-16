package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("templates/*"))

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usiario := "root"
	Contrasenia := ""
	NombreDB := "sistema"

	conexion, err := sql.Open(Driver, Usiario+":"+Contrasenia+"@tcp(127.0.0.1)/"+NombreDB)

	if err != nil {
		panic(err.Error())
	}
	return conexion
}

func index(w http.ResponseWriter, r *http.Request) {

	conexionEstablesida := conexionDB()

	registros, err := conexionEstablesida.Query("SELECT *  FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arrengloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arrengloEmpleado = append(arrengloEmpleado, empleado)
	}
	fmt.Println(arrengloEmpleado)
	plantillas.ExecuteTemplate(w, "inicio", arrengloEmpleado)
}

func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablesida := conexionDB()
		insertar, err := conexionEstablesida.Prepare("INSERT INTO empleados(Nombres,Correos)VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertar.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

func eliminar(w http.ResponseWriter, r *http.Request) {
	idEmplead := r.URL.Query().Get("id")

	fmt.Println(idEmplead)
	conexionEstablesida := conexionDB()
	BorrarRegistro, err := conexionEstablesida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	BorrarRegistro.Exec(idEmplead)
	http.Redirect(w, r, "/", 301)
}

func editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablesida := conexionDB()
	registro, err := conexionEstablesida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err := registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	plantillas.ExecuteTemplate(w, "editar", empleado)
	fmt.Println(empleado)
}
func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var id = r.FormValue("id")
		var nombre = r.FormValue("nombre")
		var correo = r.FormValue("correo")

		conexionEstablesida := conexionDB()
		actualizarRegistro, err := conexionEstablesida.Prepare("UPDATE empleados SET Nombres=?, Correos=? WHERE id=?")
		if err != nil {
			panic(err.Error)
		}

		actualizarRegistro.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301)
	}
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/eliminar", eliminar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Puerto predeterminado si no se especifica
	}
	fmt.Println("Servidor en ejecución: http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("Run Server: http://localhost:3000")
	//http.ListenAndServe("localhost:3000", nil)

}
