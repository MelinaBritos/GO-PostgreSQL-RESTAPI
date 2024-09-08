package rutas

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

func PutComprasFiltroHandler(w http.ResponseWriter, r *http.Request) {
	var compras []modelos.Compra
	var filtros modelos.FiltroCompra

	if err := json.NewDecoder(r.Body).Decode(&filtros); err != nil {
		http.Error(w, "Error al decodificar los filtros: "+err.Error(), http.StatusBadRequest)
		return
	}

	db := baseDedatos.DB

	if filtros.Estado != "" {
		db = db.Where("estado ILIKE ?", filtros.Estado)
	}

	db.Find(&compras)

	json.NewEncoder(w).Encode(&compras)

}

func PostCompraHandler(w http.ResponseWriter, r *http.Request) {
	var compra modelos.Compra

	if err := json.NewDecoder(r.Body).Decode(&compra); err != nil {
		http.Error(w, "Error al decodificar la compra: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarCompra(compra); err != nil {
		http.Error(w, "Compra inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := baseDedatos.DB.Begin()
	var Producto modelos.Producto
	baseDedatos.DB.Where("codigo_unico = ?", compra.CodigoProducto).First(&Producto)
	var productoCatalogo modelos.Catalogo

	err := tx.Where("codigo_producto = ?", Producto.CodigoUnico).First(&productoCatalogo).Error
	if err != nil {
		tx.Rollback()
		http.Error(w, "el producto no existe en el catalogo: "+err.Error(), http.StatusInternalServerError)
	}

	compra.Estado = "Completado"
	compra.Monto = productoCatalogo.PrecioActual * float32(compra.CantComprada)
	compra.Fecha = time.Now().Format("02-01-2006")
	Producto.StockDisponible += compra.CantComprada
	tx.Save(Producto)

	compraCreada := tx.Create(&compra)
	err1 := compraCreada.Error
	if err1 != nil {
		tx.Rollback()
		http.Error(w, "Error al crear la compra: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
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

func validarCompra(compra modelos.Compra) error {
	var producto modelos.Producto

	err := baseDedatos.DB.Where("codigo_unico = ?", compra.CodigoProducto).First(&producto).Error
	if err != nil {
		return errors.New("el producto no existe: " + compra.CodigoProducto)
	}
	if compra.CantComprada <= 0 {
		return errors.New("cantidad a comprar invalida")
	}
	return nil
}
