package modelos

import "gorm.io/gorm"

type ProductoVenta struct {
	gorm.Model

	Nombre string `gorm:"not null"`
	Tipo   string `gorm:"not null"`
	Marca  string `gorm:"not null"`
	Precio uint   `gorm:"not null"`
}
