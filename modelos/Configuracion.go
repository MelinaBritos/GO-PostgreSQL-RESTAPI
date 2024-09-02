package modelos

import "gorm.io/gorm"

type Configuracion struct {
	gorm.Model

	StockMinimo   int     `gorm:"not null"`
	PrecioDeseado float32 `gorm:"not null"`
	CantAComprar  int     `gorm:"not null"`
}
