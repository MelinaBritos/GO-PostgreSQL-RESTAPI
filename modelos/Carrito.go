package modelos

import (
	"gorm.io/gorm"
)

type Carrito struct {
	gorm.Model

	CodigoCarrito   string          `gorm:"unique;not null"`
	VentasUnitarias []VentaUnitaria `gorm:"foreignKey:CodigoCarritoFK;references:CodigoCarrito;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MontoTotal      float32
	FechaVenta      string
	Estado          string
}
