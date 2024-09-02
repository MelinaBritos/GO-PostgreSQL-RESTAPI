package modelos

import (
	"time"

	"gorm.io/gorm"
)

type Carrito struct {
	gorm.Model

	CodigoCarrito   string          `gorm:"unique;not null"`
	VentasUnitarias []VentaUnitaria `gorm:"foreignKey:CodigoCarritoFK;references:CodigoCarrito;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MontoTotal      float32
	FechaVenta      time.Time `gorm:"type:date;default:CURRENT_DATE"`
	Estado          string
}
