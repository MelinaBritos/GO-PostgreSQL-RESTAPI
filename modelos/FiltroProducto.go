package modelos

import "gorm.io/gorm"

type FiltroProducto struct {
	gorm.Model

	Nombre string
	Tipo   string
	Marca  string
}
