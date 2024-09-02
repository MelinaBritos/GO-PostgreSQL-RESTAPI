package rutas

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func PostConfiguracionHandler(w http.ResponseWriter, r *http.Request) {
	var configuracion modelos.Configuracion
	var producto modelos.Producto
	parametros := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&configuracion); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarConfiguracion(configuracion); err != nil {
		http.Error(w, "Configuracion inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := baseDedatos.DB.Begin()

	tx.First(&producto, parametros["id"])
	if producto.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Producto no encontrado"))
		return
	}

	if configuracion.StockMinimo != 0 {
		producto.StockMinimo = configuracion.StockMinimo
	}
	if configuracion.CantAComprar != 0 {
		producto.CantAComprar = configuracion.CantAComprar
	}
	if configuracion.PrecioDeseado != 0 {
		producto.PrecioDeseado = configuracion.PrecioDeseado
	}

	tx.Save(&producto)
	tx.Commit()
	w.WriteHeader(http.StatusOK)

}

func validarConfiguracion(configuracion modelos.Configuracion) error {
	if configuracion.CantAComprar == 0 && configuracion.PrecioDeseado == 0 && configuracion.StockMinimo == 0 {
		return errors.New("configuracion invalida")
	}
	if configuracion.CantAComprar < 0 || configuracion.PrecioDeseado < 0 || configuracion.StockMinimo < 0 {
		return errors.New("configuracion invalida, los valores no pueden ser negativos")
	}
	return nil
}
