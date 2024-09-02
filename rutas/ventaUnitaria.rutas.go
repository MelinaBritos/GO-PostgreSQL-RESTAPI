package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetVentasUnitariasHandler(w http.ResponseWriter, r *http.Request) {
	var VentasUnitarias []modelos.VentaUnitaria

	baseDedatos.DB.Find(&VentasUnitarias)
	json.NewEncoder(w).Encode(&VentasUnitarias)

}

func GetVentaUnitariaHandler(w http.ResponseWriter, r *http.Request) {
	var ventaUnitaria modelos.VentaUnitaria
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&ventaUnitaria, parametros["id"])

	if ventaUnitaria.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Venta unitaria no encontrado"))
		return
	}

	json.NewEncoder(w).Encode(&ventaUnitaria)

}
