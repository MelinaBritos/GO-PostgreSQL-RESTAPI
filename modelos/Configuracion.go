package modelos

import "gorm.io/gorm"

type Configuracion struct {
	gorm.Model

	StockMinimo   uint `gorm:"not null"`
	PrecioDeseado uint `gorm:"not null"`
	CantAComprar  uint `gorm:"not null"`
}
