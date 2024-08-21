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

	// Rutas de productos

	r.HandleFunc("/productos", rutas.GetProductosHandler).Methods("GET")
	r.HandleFunc("/productos/{id}", rutas.GetProductoHandler).Methods("GET")
	r.HandleFunc("/productos", rutas.PostProductoHandler).Methods("POST")
	r.HandleFunc("/productos/{id}", rutas.DeleteProductoHandler).Methods("DELETE")

	// Rutas de ventas

	r.HandleFunc("/ventas", rutas.GetVentasHandler).Methods("GET")
	r.HandleFunc("/ventas/{id}", rutas.GetVentaHandler).Methods("GET")
	r.HandleFunc("/ventas", rutas.PostVentaHandler).Methods("POST")
	r.HandleFunc("/ventas/{id}", rutas.DeleteVentaHandler).Methods("DELETE")

	http.ListenAndServe(":3000", r)
}
