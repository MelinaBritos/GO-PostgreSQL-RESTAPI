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

func GetCarritosHandler(w http.ResponseWriter, r *http.Request) {
	var Carritos []modelos.Carrito

	baseDedatos.DB.Find(&Carritos)

	if err := baseDedatos.DB.Preload("VentasUnitarias").Find(&Carritos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&Carritos)
}

func GetCarritoHandler(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	var Carrito modelos.Carrito

	baseDedatos.DB.First(&Carrito, parametros["id"])

	if Carrito.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // error 404
		w.Write([]byte("Carrito no encontrado"))
		return
	}

	baseDedatos.DB.Model(&Carrito).Association("VentasUnitarias").Find(&Carrito.VentasUnitarias)
	json.NewEncoder(w).Encode(&Carrito)
}

func PostCarritoHandler(w http.ResponseWriter, r *http.Request) {
	var Carrito modelos.Carrito

	if err := json.NewDecoder(r.Body).Decode(&Carrito); err != nil {
		http.Error(w, "Error al decodificar el carrito: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := baseDedatos.DB.Begin()

	CarritoCreado := tx.Create(&Carrito)
	err := CarritoCreado.Error
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al crear el carrito: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarVentaUnitaria(Carrito); err != nil {
		tx.Rollback()
		http.Error(w, "Venta invalida: "+err.Error(), http.StatusBadRequest)
		return
	}

	Carrito.Estado = "Pendiente"
	Carrito.FechaVenta = time.Now().Format("02-01-2006")
	if err := tx.Save(&Carrito).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar el carrito: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx.Commit()

	tx1 := baseDedatos.DB.Begin()

	err = restarStock(Carrito)
	if err != nil {
		tx1.Rollback()
		http.Error(w, "Error al restar stock: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = calcularMontoVentaUnitaria(Carrito)
	if err != nil {
		tx1.Rollback()
		http.Error(w, "Error al calcular monto de venta: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = realizarComprasAutomaticas(Carrito)
	if err != nil {
		tx1.Rollback()
		http.Error(w, "Error al realizar compras automaticas: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx1.Commit()
}

func DeleteCarritoHandler(w http.ResponseWriter, r *http.Request) {
	var Carrito modelos.Carrito
	parametros := mux.Vars(r)

	baseDedatos.DB.First(&Carrito, parametros["id"])

	if Carrito.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Carrito no encontrado"))
		return
	}

	baseDedatos.DB.Unscoped().Delete(&Carrito)
	w.WriteHeader(http.StatusOK)
}

func PutCarritoHandler(w http.ResponseWriter, r *http.Request) {
	var Carrito modelos.Carrito
	parametros := mux.Vars(r)
	var nuevoEstado string

	baseDedatos.DB.First(&Carrito, parametros["id"])

	if Carrito.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Carrito no encontrado"))
		return
	}

	json.NewDecoder(r.Body).Decode(&nuevoEstado)

	if err := validarEstado(nuevoEstado); err != nil {
		http.Error(w, "Estado inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if Carrito.Estado == "Confirmado" {
		http.Error(w, "No se puede cancelar una venta ya confirmada", http.StatusBadRequest)
		return
	}
	if nuevoEstado == "confirmado" {
		Carrito.Estado = "Confirmado"
		baseDedatos.DB.Save(Carrito)
	}
	if nuevoEstado == "cancelado" {
		baseDedatos.DB.Unscoped().Delete(&Carrito)
	}
	w.WriteHeader(http.StatusOK)
}

func restarStock(carrito modelos.Carrito) error {

	for _, ventaUnitaria := range carrito.VentasUnitarias {
		var producto modelos.Producto
		err := baseDedatos.DB.Where("codigo_unico = ?", ventaUnitaria.CodigoUnicoProducto).First(&producto).Error
		if err != nil {
			return errors.New("error al buscar producto" + err.Error())
		}
		producto.StockDisponible -= ventaUnitaria.Cantidad
		baseDedatos.DB.Save(&producto)
	}
	return nil
}

func calcularMontoVentaUnitaria(carrito modelos.Carrito) error {
	var montoTotal float32
	for _, ventaUnitaria := range carrito.VentasUnitarias {
		var producto modelos.Producto

		err := baseDedatos.DB.Where("codigo_unico = ?", ventaUnitaria.CodigoUnicoProducto).First(&producto).Error
		if err != nil {
			return errors.New("error al buscar producto" + err.Error())
		}

		baseDedatos.DB.Model(&ventaUnitaria).Update("monto", producto.Precio*float32(ventaUnitaria.Cantidad))
		montoTotal += producto.Precio * float32(ventaUnitaria.Cantidad)
	}
	baseDedatos.DB.Model(&carrito).Update("monto_total", montoTotal)
	return nil
}

func realizarComprasAutomaticas(carrito modelos.Carrito) error {

	tx := baseDedatos.DB.Begin()

	for _, ventaUnitaria := range carrito.VentasUnitarias {
		var Producto modelos.Producto
		err := tx.Where("codigo_unico = ?", ventaUnitaria.CodigoUnicoProducto).First(&Producto).Error
		if err != nil {
			return errors.New("error al buscar producto" + err.Error())
		}

		// Reabastecer si el stock es bajo
		if Producto.StockDisponible <= Producto.StockMinimo {
			var compra modelos.Compra
			var productoCatalogo modelos.Catalogo

			compra.CodigoProducto = Producto.CodigoUnico

			err := tx.Where("codigo_producto = ?", Producto.CodigoUnico).First(&productoCatalogo).Error
			if err != nil {
				return errors.New("el producto no existe en el catalogo: " + err.Error())
			}
			if productoCatalogo.PrecioActual > Producto.PrecioDeseado {
				compra.Estado = "No completado"
			} else {
				compra.Estado = "Completado"
				compra.CantComprada = Producto.CantAComprar
				compra.Monto = productoCatalogo.PrecioActual * float32(Producto.CantAComprar)
				Producto.StockDisponible += Producto.CantAComprar
				compra.Fecha = time.Now().Format("02-01-2006")
			}

			compraCreada := tx.Create(&compra)
			err1 := compraCreada.Error
			if err1 != nil {
				tx.Rollback()
				return errors.New("error al crear compra: " + err1.Error())
			}
			tx.Save(&Producto)
		}
	}
	tx.Commit()
	return nil
}

func validarVentaUnitaria(carrito modelos.Carrito) error {

	for _, venta := range carrito.VentasUnitarias {
		var producto modelos.Producto

		err := baseDedatos.DB.Where("codigo_unico = ?", venta.CodigoUnicoProducto).First(&producto).Error
		if err != nil {
			return errors.New("el producto no existe: " + venta.CodigoUnicoProducto)
		}
		if venta.CodigoCarritoFK != carrito.CodigoCarrito {
			return errors.New("el codigo de venta es erroneo: " + venta.CodigoCarritoFK)
		}
		if producto.StockDisponible < venta.Cantidad {
			return errors.New("no hay suficiente stock: " + venta.CodigoUnicoProducto)
		}
		if venta.Cantidad <= 0 {
			return errors.New("la cantidad es invalida")
		}
	}
	return nil
}

func validarEstado(nuevoEstado string) error {
	if nuevoEstado != "confirmado" && nuevoEstado != "cancelado" {
		return errors.New("el nuevo estado debe ser confirmado o cancelado")
	}
	return nil
}
