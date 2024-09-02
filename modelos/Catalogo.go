package modelos

import (
	"gorm.io/gorm"
)

type Catalogo struct {
	gorm.Model

	CodigoProducto string  `gorm:"unique;not null"`
	PrecioActual   float32 `gorm:"not null"`
	Proveedor      string
}
