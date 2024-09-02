package modelos

import (
	"time"

	"gorm.io/gorm"
)

type Compra struct {
	gorm.Model

	CodigoProducto string    `gorm:"not null"`
	CantComprada   int       `gorm:"not null"`
	Monto          float32   `gorm:"not null"` //Monto es precio de producto(catalogo) x cantComprada
	Fecha          time.Time `gorm:"type:date;default:CURRENT_DATE"`
	Estado         string    `gorm:"not null"`
}
