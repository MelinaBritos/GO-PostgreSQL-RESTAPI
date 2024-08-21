package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetVentasHandler(w http.ResponseWriter, r *http.Request) {
	var ventas []modelos.Venta

	baseDedatos.DB.Find(&ventas)

	json.NewEncoder(w).Encode(&ventas)
}

func GetVentaHandler(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	var venta modelos.Venta

	baseDedatos.DB.First(&venta, parametros["id"])

	if venta.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Venta no encontrada"))
		return
	}

	json.NewEncoder(w).Encode(&venta)
}

func PostVentaHandler(w http.ResponseWriter, r *http.Request) {
	var venta modelos.Venta

	json.NewDecoder(r.Body).Decode(&venta)

	ventaCreada := baseDedatos.DB.Create(&venta)

	err := ventaCreada.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&venta)
}

func DeleteVentaHandler(w http.ResponseWriter, r *http.Request) {
	var venta modelos.Venta
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&venta, parametros["id"])

	if venta.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Venta no encontrada"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&venta)
	w.WriteHeader(http.StatusOK)
}
