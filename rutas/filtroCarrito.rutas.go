package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
)

func PostCarritosFiltroHandler(w http.ResponseWriter, r *http.Request) {
	var carritos []modelos.Carrito
	var filtros modelos.FiltroCarrito

	if err := json.NewDecoder(r.Body).Decode(&filtros); err != nil {
		http.Error(w, "Error al decodificar los filtros: "+err.Error(), http.StatusBadRequest)
		return
	}

	db := baseDedatos.DB

	if filtros.Tipo == "mensual" {
		db = db.Where("fechaVenta ILIKE ?", "%-"+filtros.Fecha)
	}
	if filtros.Tipo == "diario" {
		db = db.Where("fechaVenta ILIKE ?", filtros.Fecha)
	}

	db.Find(&carritos)

	json.NewEncoder(w).Encode(&carritos)

}
