package baseDedatos

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=go-postgresql-apirest-melina-67b1.h.aivencloud.com user=avnadmin password=AVNS_Klv6DpbXmApLvr1axHR dbname=defaultdb port=22433"
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
