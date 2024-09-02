package rutas

import (
	"encoding/json"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func GetCatalogosHandler(w http.ResponseWriter, r *http.Request) {
	var Catalogos []modelos.Catalogo

	baseDedatos.DB.Find(&Catalogos)
	json.NewEncoder(w).Encode(&Catalogos)

}

func GetCatalogoHandler(w http.ResponseWriter, r *http.Request) {
	var Catalogo modelos.Catalogo
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&Catalogo, parametros["id"])

	if Catalogo.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Catalogo no encontrado"))
		return
	}

	json.NewEncoder(w).Encode(&Catalogo)

}

func PostCatalogosHandler(w http.ResponseWriter, r *http.Request) {
	var Catalogos []modelos.Catalogo

	json.NewDecoder(r.Body).Decode(&Catalogos)

	tx := baseDedatos.DB.Begin()
	for _, catalogo := range Catalogos {

		catalogoCreado := tx.Create(&catalogo)

		err := catalogoCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear el Catalogo: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

func DeleteCatalogoHandler(w http.ResponseWriter, r *http.Request) {
	var Catalogo modelos.Catalogo
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&Catalogo, parametros["id"])

	if Catalogo.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Catalogo no encontrado"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&Catalogo)
	w.WriteHeader(http.StatusOK)

}
