package rutas

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/baseDedatos"
	"github.com/MelinaBritos/GO-PostgreSQL-RESTAPI/modelos"
	"github.com/gorilla/mux"
)

func PutConfiguracionHandler(w http.ResponseWriter, r *http.Request) {
	var configuracion modelos.Configuracion
	var producto modelos.Producto
	var productocompra modelos.ProductoCompra
	parametros := mux.Vars(r)

	json.NewDecoder(r.Body).Decode(&configuracion)
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
		productocompra.CantAComprar = configuracion.CantAComprar
	}
	if configuracion.PrecioDeseado != 0 {
		productocompra.PrecioDeseado = configuracion.PrecioDeseado
	}

	tx.Save(&producto)
	tx.Save(&productocompra)
	tx.Commit()
	w.WriteHeader(http.StatusOK)

}

func validarConfiguracion(configuracion modelos.Configuracion) error {
	if configuracion.CantAComprar == 0 && configuracion.PrecioDeseado == 0 && configuracion.StockMinimo == 0 {
		return errors.New("configuracion invalida")
	}
	return nil
}
