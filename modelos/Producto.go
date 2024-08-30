package modelos

import "gorm.io/gorm"

type Producto struct {
	gorm.Model

	CodigoUnico     string  `gorm:"unique;not null"`
	Nombre          string  `gorm:"not null"`
	Tipo            string  `gorm:"not null"`
	Marca           string  `gorm:"not null"`
	StockDisponible float64 `gorm:"not null"`
	StockMinimo     float64 `gorm:"not null"`
}
