package modelos

import (
	"gorm.io/gorm"
)

type Compra struct {
	gorm.Model

	CodigoProductoCompra string `gorm:"not null"`
	Precio               uint   `gorm:"not null"`
	Fecha                string `gorm:"not null"`
	Estado               string `gorm:"not null"`
}
