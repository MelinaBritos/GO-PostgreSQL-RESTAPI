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
	baseDedatos.DB.AutoMigrate(modelos.VentaUnitaria{})
	baseDedatos.DB.AutoMigrate(modelos.Carrito{})
	baseDedatos.DB.AutoMigrate(modelos.Compra{})
	baseDedatos.DB.AutoMigrate(modelos.Catalogo{})

	r := mux.NewRouter()
	r.HandleFunc("/", rutas.HomeHandler)

	// Rutas de productos

	r.HandleFunc("/productos", rutas.GetProductosHandler).Methods("GET")
	r.HandleFunc("/productosFiltro", rutas.GetProductosFiltroHandler).Methods("GET")
	r.HandleFunc("/productos/{id}", rutas.GetProductoHandler).Methods("GET")
	r.HandleFunc("/productos", rutas.PostProductosHandler).Methods("POST")
	r.HandleFunc("/productos/{id}", rutas.DeleteProductoHandler).Methods("DELETE")

	// Rutas de carrito

	r.HandleFunc("/carrito", rutas.GetCarritosHandler).Methods("GET")
	r.HandleFunc("/carrito/{id}", rutas.GetCarritoHandler).Methods("GET")
	r.HandleFunc("/carrito", rutas.PostCarritoHandler).Methods("POST")
	r.HandleFunc("/carrito/{id}", rutas.DeleteCarritoHandler).Methods("DELETE")
	r.HandleFunc("/carrito/{id}", rutas.PutCarritoHandler).Methods("PUT")

	// Rutas de catalogo

	r.HandleFunc("/catalogo", rutas.GetCatalogosHandler).Methods("GET")
	r.HandleFunc("/catalogo/{id}", rutas.GetCatalogoHandler).Methods("GET")
	r.HandleFunc("/catalogo", rutas.PostCatalogosHandler).Methods("POST")
	r.HandleFunc("/catalogo/{id}", rutas.DeleteCatalogoHandler).Methods("DELETE")

	// Rutas de compras

	r.HandleFunc("/compras", rutas.GetComprasHandler).Methods("GET")
	r.HandleFunc("/compras/{id}", rutas.GetCompraHandler).Methods("GET")
	r.HandleFunc("/comprasFiltro", rutas.GetComprasFiltroHandler).Methods("GET")
	r.HandleFunc("/compras", rutas.PostCompraHandler).Methods("POST")
	r.HandleFunc("/compras/{id}", rutas.DeleteCompraHandler).Methods("DELETE")

	// Rutas de configuracion

	r.HandleFunc("/configuracion/{id}", rutas.PostConfiguracionHandler).Methods("POST")

	// Rutas de ventas unitarias

	r.HandleFunc("/ventaUnitaria", rutas.GetVentasUnitariasHandler).Methods("GET")
	r.HandleFunc("/ventaUnitaria/{id}", rutas.GetVentaUnitariaHandler).Methods("GET")

	// Configura CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Inicia el servidor con CORS habilitado
	http.ListenAndServe(":8080", corsHandler(r))

}
