package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
)

func PostProductosFiltroHandler(w http.ResponseWriter, r *http.Request) {
	var productos []modelos.Producto
	var filtros modelos.FiltroProducto

	if err := json.NewDecoder(r.Body).Decode(&filtros); err != nil {
		http.Error(w, "Error al decodificar los filtros: "+err.Error(), http.StatusBadRequest)
		return
	}

	db := baseDedatos.DB

	if filtros.Nombre != "" {
		db = db.Where("nombre ILIKE ?", "%"+filtros.Nombre+"%")
	}
	if filtros.Marca != "" {
		db = db.Where("marca ILIKE ?", "%"+filtros.Marca+"%")
	}
	if filtros.Tipo != "" {
		db = db.Where("tipo ILIKE ?", "%"+filtros.Tipo+"%")
	}

	db.Find(&productos)

	json.NewEncoder(w).Encode(&productos)

}
