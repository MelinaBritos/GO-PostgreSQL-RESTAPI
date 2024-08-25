package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetProductosCompraHandler(w http.ResponseWriter, r *http.Request) {
	var productosCompra []modelos.ProductoCompra

	baseDedatos.DB.Find(&productosCompra)
	json.NewEncoder(w).Encode(&productosCompra)

}

func GetProductoCompraHandler(w http.ResponseWriter, r *http.Request) {
	var productoCompra modelos.ProductoCompra
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&productoCompra, parametros["id"])

	if productoCompra.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("ProductoCompra no encontrado"))
		return
	}

	json.NewEncoder(w).Encode(&productoCompra)

}

func PostProductosCompraHandler(w http.ResponseWriter, r *http.Request) {
	var productosCompra []modelos.ProductoCompra

	json.NewDecoder(r.Body).Decode(&productosCompra)

	tx := baseDedatos.DB.Begin()
	for _, producto := range productosCompra {
		productoCreado := tx.Create(&producto)

		err := productoCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear el productoCompra: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func DeleteProductoCompraHandler(w http.ResponseWriter, r *http.Request) {
	var productoCompra modelos.ProductoCompra
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&productoCompra, parametros["id"])

	if productoCompra.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ProductoCompra no encontrado"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&productoCompra)
	w.WriteHeader(http.StatusOK)

}
