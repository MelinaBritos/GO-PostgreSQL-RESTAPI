package modelos

import "gorm.io/gorm"

type FiltroCompra struct {
	gorm.Model

	Estado string
}
