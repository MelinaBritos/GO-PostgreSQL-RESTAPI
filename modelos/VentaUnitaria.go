package modelos

import (
	"gorm.io/gorm"
)

type VentaUnitaria struct {
	gorm.Model

	CodigoCarritoFK     string  `gorm:"not null"`
	CodigoUnicoProducto string  `gorm:"not null"`
	Cantidad            int     `gorm:"not null"`
	Monto               float32 // precio de producto x cantidad

}
