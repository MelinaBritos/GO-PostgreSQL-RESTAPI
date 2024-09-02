package rutas

import (
	"encoding/json"
	"errors"
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

func PutProductosFiltroHandler(w http.ResponseWriter, r *http.Request) {
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

func PostProductosHandler(w http.ResponseWriter, r *http.Request) {
	var productos []modelos.Producto

	if err := json.NewDecoder(r.Body).Decode(&productos); err != nil {
		http.Error(w, "Error al decodificar los productos: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, producto := range productos {
		if err := validarProducto(producto); err != nil {
			http.Error(w, "Producto inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	tx := baseDedatos.DB.Begin()
	for _, producto := range productos {
		productoCreado := tx.Create(&producto)

		err := productoCreado.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear los productos: "+err.Error(), http.StatusInternalServerError)
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

func validarProducto(producto modelos.Producto) error {
	if producto.CodigoUnico == "" {
		return errors.New("codigo unico no puede estar vacio")
	}
	if producto.Nombre == "" {
		return errors.New("nombre no puede estar vacio")
	}
	if producto.Tipo == "" {
		return errors.New("tipo no puede estar vacio")
	}
	if producto.Marca == "" {
		return errors.New("marca no puede estar vacia")
	}
	if producto.StockDisponible <= 0 || producto.StockMinimo <= 0 {
		return errors.New("stock minimo y stock disponible no pueden ser igual o menor a cero")
	}
	if producto.Precio <= 0 {
		return errors.New("precio no puede ser cero")
	}
	if producto.CantAComprar <= 0 {
		return errors.New("cantidad a comprar no puede ser igual o menor a cero")
	}
	if producto.PrecioDeseado <= 0 {
		return errors.New("precio deseado no puede ser igual o menor a cero")
	}
	return nil
}
