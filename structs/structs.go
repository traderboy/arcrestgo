package structs

import (
	"database/sql"
	"encoding/json"
)

type Catalog struct {
	Services        JSONConfig
	DataSource      int
	SqliteSource    string
	Postgresql      string
	RootPath        string
	Schema          string
	Db              *sql.DB
	DbQuery         *sql.DB
	DataPath        string
	ReplicaPath     string
	AttachmentsPath string
}

//JSONConfig stores the metadata about a service
type JSONConfig struct {
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	//Services
	Services   map[string]map[string]map[string]map[string]interface{} `json:"services"`
	FGDB       string                                                  `json:fgdb`
	MXD        string                                                  `json:mxd`
	PG         string                                                  `json:pg`
	SqliteDb   string                                                  `json:sqliteDb`
	DataSource string                                                  `json:dataSource`
	RootName   string                                                  `json:`
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
		Code interface{} `json:"code"`
		Name string      `json:"name"`
	} `json:"codedValues,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type RelatedRecords struct {
	Fields              []Field              `json:"fields,omitempty"`
	RelatedRecordGroups []RelatedRecordGroup `json:"relatedRecordGroups"`
}

type RelatedRecordGroup struct {
	ObjectId int `json:"objectId"`
	//RelatedRecord []map[string]interface{} `json:"relatedRecords"`
	RelatedRecords []RelatedRecord `json:"relatedRecords"`
}

type RelatedRecord struct {
	//Attributes []Attribute `json:"attributes"`
	Attributes map[string]interface{} `json:"attributes"`
}

/*
type Attribute struct {
	Attributes map[string]interface{}
	//`json:"attributes"`
}
*/
type Record []struct {
	Attributes map[string]interface{} `json:"attributes"`
	//RelatedRecordGroups RelatedRecordGroup     `json:"relatedRecordGroups"`
}

type Geometry struct {
	Rings [][][]float64 `json:"rings,omitempty"`
	Y     float64       `json:"y,omitempty"`
	X     float64       `json:"x,omitempty"`
}

type Feature struct {
	Geometry *Geometry `json:"geometry,omitempty"`
	//Attributes Attribute `json:"attributes,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	//Attributes map[string][]Attribute `json:"attributes,omitempty"`
}

type FeatureTable struct {
	GlobalIDField     string `json:"globalIdField,omitempty"`
	GlobalIDFieldName string `json:"globalIdFieldName,omitempty"`
	SpatialReference  *struct {
		Wkid       *int `json:"wkid,omitempty"`
		LatestWkid *int `json:"latestWkid,omitempty"`
	} `json:"spatialReference,omitempty"`
	GeometryType       string `json:"geometryType,omitempty"`
	ObjectIDField      string `json:"objectIdField,omitempty"`
	GeometryProperties *struct {
		ShapeLengthFieldName string `json:"shapeLengthFieldName,omitempty"`
		Units                string `json:"units,omitempty"`
	} `json:"geometryProperties,omitempty"`
	ObjectIDFieldName string    `json:"objectIdFieldName,omitempty"`
	DisplayFieldName  string    `json:"displayFieldName,omitempty"`
	Fields            []Field   `json:"fields,omitempty"`
	Features          []Feature `json:"features,omitempty"`
}

type Items struct {
	ObjectID        int    `json:"objectId"`
	UUID            string `json:"uuid"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	PhysicalName    string `json:"physicalName"`
	Path            string `json:"path"`
	Url             string `json:"url"`
	Properties      int32  `json:"properties"`
	Defaults        []byte `json:"defaults"`
	DatasetSubtype1 int32  `json:"datasetSubtype1"`
	DatasetSubtype2 int32  `json:"datasetSubtype2"`
	DatasetInfo1    string `json:"datasetInfo1"`
	DatasetInfo2    string `json:"datasetInfo"`
	Definition      string `json:"definition"`
	Documentation   string `json:"documentation"`
	ItemInfo        string `json:"itemInfo"`
	Shape           []byte `json:"shape"`
}

type ServiceItems struct {
	ObjectID            int    `json:"objectId"`
	DatasetName         string `json:"datasetName"`
	ItemType            int32  `json:"itemType"`
	ItemId              int32  `json:"itemId"`
	ItemInfo            string `json:"itemInfo"`
	AdvancedDrawingInfo string `json:"advancedDrawingInfo"`
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
