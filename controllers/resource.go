package controllers

import (
	"../models"
	"encoding/json"
	"github.com/Jamshid90/fhir-schema"
	"io/ioutil"
	"net/http"
	"strings"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	path_split := strings.Split(r.URL.Path, "/")
	resourceType := path_split[2]
	resource, err := schema.GetFhirResourceInstance(resourceType)

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"message" : err.Error(),
		})
		return
	}

	rc := ResourceController{
		ResourceType : resourceType,
		Resource : resource,
	}

	switch r.Method {
		case "GET":
			rc.GetResource(w, r)
		case "POST":
			rc.CreateResource(w, r)
	}
}

type ResourceController struct {
	ResourceType string
	Resource interface{}
}

func (rc *ResourceController) GetResource(w http.ResponseWriter, r *http.Request)  {

	request_query := r.URL.Query()
	model_resource := models.Resource{ResourceType:rc.ResourceType}
	criteria := model_resource.BuildCriteria(request_query)

	items, err := model_resource.Get(criteria)

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"message" : err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		"success" : true,
		"data" : items,
	})
	return
}

func (rc *ResourceController) CreateResource(w http.ResponseWriter, r *http.Request)  {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(b, &rc.Resource)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"message" : err.Error(),
		})
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"message" : err.Error(),
		})
		return
	}

	if data["resourceType"] != rc.ResourceType{
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"data" : "Resource type invalid",
		})
		return
	}

	model_resource := models.Resource{
		Data:rc.Resource,
		ResourceType:rc.ResourceType,
	}
	err = model_resource.Create()

	if err != nil {
		json.NewEncoder(w).Encode(Response{
			"success" : false,
			"message" : err.Error(),
		})
	}

	json.NewEncoder(w).Encode(Response{
		"success" : true,
		"data" : model_resource,
	})

}
