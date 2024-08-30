package modelos

import (
	"time"

	"gorm.io/gorm"
)

type ProductoCompra struct {
	gorm.Model

	Codigo        string    `gorm:"unique;not null"` //fk
	PrecioActual  uint      `gorm:"not null"`
	PrecioDeseado uint      `gorm:"not null"`
	Fecha         time.Time `gorm:"type:date;default:CURRENT_DATE"`
	Proveedor     string
	CantAComprar  uint `gorm:"not null"`
}
