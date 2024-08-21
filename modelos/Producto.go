package modelos

import "gorm.io/gorm"

type Producto struct {
	gorm.Model

	Nombre          string `gorm:"not null"`
	Tipo            string `gorm:"not null"`
	Marca           string `gorm:"not null"`
	StockDisponible uint   `gorm:"not null"`
	StockMinimo     uint   `gorm:"not null"`
}
