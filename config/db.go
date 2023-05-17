package config

import (
	log "DIA-NFT-Sales-Bot/debug"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() {
	var err error
	host := "dia.ep-fragrant-limit-541358.us-east-2.aws.neon.tech"
	port := "5432"
	dbname := "neondb"
	user := "Brymes"
	password := "va2pVxqIodY5"
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=require options=project=ep-fragrant-limit-541358",
		host,
		user,
		password,
		dbname,
		port,
	)

	DBClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Log.Println("Error Connecting to Database. Kindly set accurate Database environment variables")
		log.Log.Fatal(err)
	}
}

//docker run --name local-psql -v local_psql_data:/var/lib/postgresql/data -p 54320:5432 -e POSTGRES_PASSWORD=my_password -d postgres
