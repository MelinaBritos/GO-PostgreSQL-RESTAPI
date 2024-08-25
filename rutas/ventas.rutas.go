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

func PostVentasHandler(w http.ResponseWriter, r *http.Request) {
	var ventas []modelos.Venta

	json.NewDecoder(r.Body).Decode(&ventas)

	tx := baseDedatos.DB.Begin()

	for _, venta := range ventas {
		var Producto modelos.Producto
		err := tx.Where("codigo_unico = ?", venta.CodigoUnicoProducto).First(&Producto).Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Producto no encontrado: "+err.Error(), http.StatusNotFound)
			return
		}

		//Restar el stock de la venta si es posible
		if Producto.StockDisponible < venta.Cantidad {
			tx.Rollback()
			http.Error(w, "Stock no disponible para realizar la venta", http.StatusBadRequest)
			return
		}
		Producto.StockDisponible -= venta.Cantidad

		// Reabastecer si el stock es bajo
		if Producto.StockDisponible < Producto.StockMinimo {
			var compra modelos.Compra
			var productoCompra modelos.ProductoCompra
			err := tx.Where("codigo = ?", Producto.CodigoUnico).First(&productoCompra).Error
			if err != nil {
				tx.Rollback()
				http.Error(w, "ProductoCompra no encontrado: "+err.Error(), http.StatusNotFound)
				return
			}
			compra.CodigoProductoCompra = productoCompra.Codigo
			if productoCompra.PrecioActual > productoCompra.PrecioDeseado {
				compra.Estado = "No completado"
			} else {
				compra.Estado = "Completado"
				compra.Precio = productoCompra.PrecioActual
				Producto.StockDisponible += 30
			}
			tx.Create(&compra)
		}
		tx.Save(&Producto)
	}

	//Crear las ventas
	for _, venta := range ventas {
		ventaCreada := tx.Create(&venta)
		err := ventaCreada.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear las ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	w.WriteHeader(http.StatusOK)
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
