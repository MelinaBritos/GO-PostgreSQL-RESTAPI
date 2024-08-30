package rutas

import (
	"encoding/json"
	"errors"
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

	for _, venta := range ventas {
		if err := validarVenta(venta); err != nil {
			http.Error(w, "Venta inválida: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	tx := baseDedatos.DB.Begin()

	for _, venta := range ventas {

		// Validar que exista el producto de la venta
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
			var compra modelos.Compra //creo una compra automatica
			var productoCompra modelos.ProductoCompra
			err := tx.Where("codigo = ?", Producto.CodigoUnico).First(&productoCompra).Error //busco el producto compra que coincida con el producto vendido
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
				compra.Precio = productoCompra.PrecioActual * venta.Cantidad
				Producto.StockDisponible += productoCompra.CantAComprar
			}
			tx.Create(&compra)
		}
		tx.Save(&Producto)

	}

	//Crear las ventas y sumar monto total
	var montoTotal uint
	for _, venta := range ventas {
		venta.Estado = "Pendiente"
		montoTotal += venta.Monto
		tx.Save(venta)
		ventaCreada := tx.Create(&venta)
		err := ventaCreada.Error
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al crear las ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	json.NewEncoder(w).Encode(&montoTotal)
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

func PutVentaHandler(w http.ResponseWriter, r *http.Request) {
	var venta modelos.Venta
	parametros := mux.Vars(r)
	var nuevoEstado string

	baseDedatos.DB.First(&venta, parametros["id"])

	if venta.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Venta no encontrada"))
		return
	}

	json.NewDecoder(r.Body).Decode(&nuevoEstado)
	if err := validarEstado(nuevoEstado); err != nil {
		http.Error(w, "Estado inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if nuevoEstado == "confirmado" {
		venta.Estado = "Confirmado"
		baseDedatos.DB.Save(venta)
	}
	if nuevoEstado == "cancelado" {
		baseDedatos.DB.Unscoped().Delete(&venta)
	}
	w.WriteHeader(http.StatusOK)
}

func validarVenta(venta modelos.Venta) error {
	if venta.CodigoUnicoProducto == "" {
		return errors.New("codigo de producto no puede estar vacio")
	}
	if venta.Cantidad == 0 {
		return errors.New("cantidad no puede ser cero")
	}
	if venta.FechaVenta == "" {
		return errors.New("fecha venta no puede estar vacia")
	}
	if venta.Monto == 0 {
		return errors.New("monto no puede ser cero")
	}
	return nil
}

func validarEstado(nuevoEstado string) error {
	if nuevoEstado != "confirmado" && nuevoEstado != "cancelado" {
		return errors.New("el nuevo estado debe ser confirmado o cancelado")
	}
	return nil
}
