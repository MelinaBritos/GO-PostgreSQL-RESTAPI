package modelos

import "gorm.io/gorm"

type Venta struct {
	gorm.Model

	ProductoID int      `gorm:"not null"`
	Producto   Producto `gorm:"foreignKey:ProductoID`
	Cantidad   uint     `gorm:"not null"`
	FechaVenta string   `gorm:"not null"`
}
