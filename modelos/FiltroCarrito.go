package modelos

import "gorm.io/gorm"

type FiltroCarrito struct {
	gorm.Model

	Fecha string
	Tipo  string
}
