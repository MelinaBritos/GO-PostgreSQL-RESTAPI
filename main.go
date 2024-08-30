package main

import (
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/rutas"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	baseDedatos.Conexiondb()

	baseDedatos.DB.AutoMigrate(modelos.Producto{})
	baseDedatos.DB.AutoMigrate(modelos.Venta{})
	baseDedatos.DB.AutoMigrate(modelos.Compra{})
	baseDedatos.DB.AutoMigrate(modelos.ProductoCompra{})

	r := mux.NewRouter()
	r.HandleFunc("/", rutas.HomeHandler)

	// Rutas de productos

	r.HandleFunc("/productos", rutas.GetProductosHandler).Methods("GET")
	r.HandleFunc("/productos/{id}", rutas.GetProductoHandler).Methods("GET")
	r.HandleFunc("/productos", rutas.PostProductosHandler).Methods("POST")
	r.HandleFunc("/productos/{id}", rutas.DeleteProductoHandler).Methods("DELETE")

	// Rutas de ventas

	r.HandleFunc("/ventas", rutas.GetVentasHandler).Methods("GET")
	r.HandleFunc("/ventas/{id}", rutas.GetVentaHandler).Methods("GET")
	r.HandleFunc("/ventas", rutas.PostVentasHandler).Methods("POST")
	r.HandleFunc("/ventas/{id}", rutas.DeleteVentaHandler).Methods("DELETE")
	r.HandleFunc("/ventas/{id}", rutas.PutVentaHandler).Methods("PUT")

	// Rutas de productos compra

	r.HandleFunc("/productocompra", rutas.GetProductosCompraHandler).Methods("GET")
	r.HandleFunc("/productocompra/{id}", rutas.GetProductoCompraHandler).Methods("GET")
	r.HandleFunc("/productocompra", rutas.PostProductosCompraHandler).Methods("POST")
	r.HandleFunc("/productocompra/{id}", rutas.DeleteProductoCompraHandler).Methods("DELETE")

	// Rutas de compras

	r.HandleFunc("/compras", rutas.GetComprasHandler).Methods("GET")
	r.HandleFunc("/compras/{id}", rutas.GetCompraHandler).Methods("GET")
	r.HandleFunc("/compras/{id}", rutas.DeleteCompraHandler).Methods("DELETE")

	// Rutas de configuracion

	r.HandleFunc("/configuracion/{id}", rutas.PutConfiguracionHandler).Methods("PUT")

	// Configura CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Inicia el servidor con CORS habilitado
	http.ListenAndServe(":8080", corsHandler(r))
}
