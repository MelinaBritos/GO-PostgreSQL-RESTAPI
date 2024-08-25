package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetProductosHandler(w http.ResponseWriter, r *http.Request) {
	var productos []modelos.Producto

	baseDedatos.DB.Find(&productos)
	json.NewEncoder(w).Encode(&productos)

}

func GetProductoHandler(w http.ResponseWriter, r *http.Request) {
	var producto modelos.Producto
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&producto, parametros["id"])

	if producto.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Producto no encontrado"))
		return
	}

	json.NewEncoder(w).Encode(&producto)

}

func PostProductosHandler(w http.ResponseWriter, r *http.Request) {
	var productos []modelos.Producto

	json.NewDecoder(r.Body).Decode(&productos)

	tx := baseDedatos.DB.Begin()
	for _, producto := range productos {
		productoCreado := tx.Create(&producto)

		err := productoCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear las ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func DeleteProductoHandler(w http.ResponseWriter, r *http.Request) {
	var producto modelos.Producto
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&producto, parametros["id"])

	if producto.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Producto no encontrado"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&producto)
	w.WriteHeader(http.StatusOK)

}
