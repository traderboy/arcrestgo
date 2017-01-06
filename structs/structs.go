package structs

import (
	"encoding/json"
)

//JSONConfig stores the metadata about a service
type JSONConfig struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	//Services
	Services map[string]map[string]map[string]map[string]interface{} `json:"services"`
	FGDB string `json:fgdb`
	MXD string `json:mxd`
	//Services map[string]map[string]Service
	//map[string]Service
}

type FieldsStr struct {
	Fields json.RawMessage `json:"fields"`
	//Fields []Field `json:"fields"`
}

type TableField struct {
	//Fields json.RawMessage `json:"fields"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Domain       *Domain     `json:"domain"`
	Name         string      `json:"name"`
	Nullable     bool        `json:"nullable"`
	DefaultValue interface{} `json:"defaultValue"`
	Editable     bool        `json:"editable"`
	Alias        string      `json:"alias"`
	SqlType      string      `json:"sqlType"`
	Type         string      `json:"type"`
	Length       int         `json:"length,omitempty"`
}

type Domain struct {
	CodedValues []struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"codedValues,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Record []struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type Geometry struct {
	Rings [][][]float64 `json:"rings,omitempty"`
	Y     float64       `json:"y,omitempty"`
	X     float64       `json:"x,omitempty"`
}
type Feature struct {
	Geometry   *Geometry              `json:"geometry,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type FeatureTable struct {
	GlobalIDField    string `json:"globalIdField,omitempty"`
	SpatialReference *struct {
		Wkid       *int `json:"wkid,omitempty"`
		LatestWkid *int `json:"latestWkid,omitempty"`
	} `json:"spatialReference"`
	GeometryType      string    `json:"geometryType,omitempty"`
	ObjectIDField     string    `json:"objectIdField,omitempty"`
	ObjectIDFieldName string    `json:"objectIdFieldName,omitempty"`
	DisplayFieldName  string    `json:"displayFieldName,omitempty"`
	Fields            []Field   `json:"fields,omitempty"`
	Features          []Feature `json:"features,omitempty"`
}

/*
type JSON struct{
Features [] feature `json:"features"`
"displayFieldName": "",
    "spatialReference": {
        "wkid": 102100,
        "latestWkid": 3857
    },
    "geometryType": "esriGeometryPoint",
    "objectIdField": "OBJECTID",
    "objectIdFieldName": "OBJECTID"
}
*/

/*
type Configuration struct {
	Services struct {
		Service struct {
			Layers struct {
				Layer struct {
					ItemID        string `json:"itemId"`
					Data          string `json:"data"`
					Name          string `json:"name"`
					Oidname       string `json:"oidname"`
					Globaloidname string `json:"globaloidname"`
				} `json:"0"`
			} `json:"layers"`
			Relationships struct {
				Relationship struct {
					OID      int    `json:"oId"`
					DID      int    `json:"dId"`
					OTable   string `json:"oTable"`
					OJoinKey string `json:"oJoinKey"`
					DJoinKey string `json:"dJoinKey"`
					DTable   string `json:"dTable"`
				} `json:"0"`
			} `json:"relationships"`
		} `json:"accommodationagreementrentals"`
	} `json:"services"`
	Username string `json:"username"`
	Hostname string `json:"hostname"`
}
*/

/*
type Service struct {
	//Names map[string]interface{}
	Names map[string]Name
}
type Name struct {
	//Layers map[string]interface{}
	Layers map[string]Layer
}
type Layer struct {
	Items         map[string]Item
	Relationships map[string]Relationship
}
type Item struct {
	ItemID        string `json:"itemId"`
	Data          string `json:"data"`
	Name          string `json:"name"`
	Oidname       string `json:"oidname"`
	Globaloidname string `json:"globaloidname"`
}
type Relationship struct {
	Oid    int    `json:"oId"`
	DId    int    `json:"dId"`
	OTable string `json:"oTable"`

	OJoinKey string `json:"oJoinKey"`
	DJoinKey string `json:"dJoinKey"`
	DTable   string `json:"dTable"`
}
*/
