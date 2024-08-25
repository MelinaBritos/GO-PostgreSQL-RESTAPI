package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetComprasHandler(w http.ResponseWriter, r *http.Request) {
	var compras []modelos.Compra

	baseDedatos.DB.Find(&compras)
	json.NewEncoder(w).Encode(&compras)

}

func GetCompraHandler(w http.ResponseWriter, r *http.Request) {
	var compra modelos.Compra
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&compra, parametros["id"])

	if compra.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Compra no encontrada"))
		return
	}

	json.NewEncoder(w).Encode(&compra)

}

func DeleteCompraHandler(w http.ResponseWriter, r *http.Request) {
	var compra modelos.Compra
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&compra, parametros["id"])

	if compra.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("compra no encontrada"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&compra)
	w.WriteHeader(http.StatusOK)

}
