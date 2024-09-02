package modelos

import "gorm.io/gorm"

type Producto struct {
	gorm.Model

	CodigoUnico     string  `gorm:"unique;not null"`
	Tipo            string  `gorm:"not null"`
	Nombre          string  `gorm:"not null"`
	Marca           string  `gorm:"not null"`
	StockDisponible int     `gorm:"not null"`
	StockMinimo     int     `gorm:"not null"`
	Precio          float32 `gorm:"not null"` //para calcular ventas
	CantAComprar    int     `gorm:"not null"`
	PrecioDeseado   float32 `gorm:"not null"`
	Descripcion     string  `gorm:"not null"`
}
