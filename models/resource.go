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

func (r *Resource) Get(criteria map[string][]database.Criteria) ([]Resource, error) {

	var items []Resource

	sql := fmt.Sprintf("SELECT * FROM %s ", strings.ToLower(r.ResourceType))
	where := "WHERE 1=1 "

	if criteria_where, ok := criteria["where"]; ok == true{
		for _, v := range criteria_where {

			where += fmt.Sprintf(" %s %s %s %s ", v.Relation, v.Key, v.Compare, v.Value)

		}
    	sql += fmt.Sprintf(" %s ", where)
	}

	if order, ok := criteria["order"]; ok == true{
		sql += fmt.Sprintf(" %s ", order)
	}else{
		sql += fmt.Sprintf(" %s ", "ORDER BY id DESC")
	}

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
	defer rows.Close()

	return items, nil
}

func (r *Resource) DataUnmarshal() error {
	b, ok := r.Data.([]byte)

	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &r.Data)
}

func (r *Resource) BuildCriteria(params map[string][]string) map[string][]database.Criteria {
	query := make(map[string][]database.Criteria)
	for k, v := range params{

		if k == "patient" {
			if patient_id := v[0]; patient_id != ""{
				query["where"] = append(query["where"], r.GetPatientCriteria(patient_id, "AND", "="))
			}

		}
	}
	return query
}

func (r *Resource) GetPatientCriteria(patient_id, relation, compare string) database.Criteria {
	key := " data #>> '{subject,reference}'"
	value := fmt.Sprintf("'Patient/%s'", patient_id)

	if r.ResourceType == "Appointment" {
		key = " data #>> '{participant,0,actor,reference}' "
	}

	return database.Criteria{
		Relation:relation,
		Key:key,
		Value:value,
		Compare:compare,
	}

}


