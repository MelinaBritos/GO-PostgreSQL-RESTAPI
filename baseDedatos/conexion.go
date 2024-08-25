package baseDedatos

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=bddproyectop-ayelenvalentinaa-ab78.h.aivencloud.com user=avnadmin password=AVNS_7qAktm_m5i9zUOz9xjf dbname=defaultdb port=21647"
var DB *gorm.DB

func Conexiondb() {
	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Base de datos conectada")
	}
}
