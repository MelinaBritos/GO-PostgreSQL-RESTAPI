package modelos

import "gorm.io/gorm"

type Venta struct {
	gorm.Model

	CodigoUnicoProducto string `gorm:"not null"`
	Cantidad            uint   `gorm:"not null"`
	FechaVenta          string `gorm:"not null"`
	Monto               uint   `gorm:"not null"`
	Estado              string `gorm:"not null"`
}
