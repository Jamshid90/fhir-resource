package database

import (
	"database/sql"
	"fmt"
	"github.com/Jamshid90/fhir-schema"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var db *sql.DB

func Connect(){
	var err error
	var dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s ",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_SSLMODE"))
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return db
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec( query, args...)
}

func CreateTable(){
	fhir_resourc := schema.GetFhirResourceMap()
	for k, _ := range fhir_resourc{
		query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
										id serial primary key,
										data jsonb,
										created_at timestamp(0) without time zone )`,
			strings.ToLower(k))
		Exec(query)
	}
}

func Cloce(){
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

type Criteria struct {
	Relation string
	Key string
	Value string
	Compare string
}
