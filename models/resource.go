package models

import (
	"../database"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Resource struct {
	ID uint64 `json: "id,omitempty" `
	ResourceType string `json:"resourceType,omitempty"`
	Data interface{} `json:"data,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (r *Resource) CreateTable()  {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( 
									id serial primary key, 
									data jsonb, 
									created_at timestamp(0) without time zone )`,
									strings.ToLower(r.ResourceType))
	database.Exec(query)
}

func (r *Resource) Create() error  {
	sql := fmt.Sprintf("INSERT INTO %s ( data, created_at ) VALUES( $1, $2 )", strings.ToLower(r.ResourceType))

	data, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	_, err = database.Exec(sql, string(data), time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (r *Resource) Get() ([]Resource, error) {

	var items []Resource

	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC", strings.ToLower(r.ResourceType))

	rows, err := database.Query(sql)
	if err != nil{
		return items, err
	}

	for rows.Next(){
		item := Resource{}

		rows.Scan(&item.ID, &item.Data, &item.CreatedAt )

		err = item.DataUnmarshal()
		if err == nil{
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *Resource) DataUnmarshal() error {
	b, ok := r.Data.([]byte)

	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &r.Data)
}
