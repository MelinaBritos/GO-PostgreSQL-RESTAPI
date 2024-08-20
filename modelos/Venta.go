package modelos

import "gorm.io/gorm"

type Venta struct {
	gorm.Model

	IDproducto string `gorm:"not null"`
	Cantidad   uint   `gorm:"not null"`
	FechaVenta uint   `gorm:"not null"`
}
