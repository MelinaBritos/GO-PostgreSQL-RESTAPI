package main

import (
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/rutas"
	"github.com/gorilla/mux"
)

func main() {
	baseDedatos.Conexiondb()

	baseDedatos.DB.AutoMigrate(modelos.Producto{})
	baseDedatos.DB.AutoMigrate(modelos.Venta{})

	r := mux.NewRouter()

	r.HandleFunc("/", rutas.HomeHandler)
	r.HandleFunc("/productos", rutas.GetProductosHandler).Methods("GET")
	r.HandleFunc("/productos/{id}", rutas.GetProductoHandler).Methods("GET")
	r.HandleFunc("/productos", rutas.PostProductoHandler).Methods("POST")
	r.HandleFunc("/productos", rutas.DeleteProductoHandler).Methods("DELETE")

	http.ListenAndServe(":3000", r)
}
