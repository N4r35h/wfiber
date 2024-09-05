package wfiber

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type SwaggerContact struct {
	Email string `json:"email"`
}

type SwaggerLicence struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SwaggerInfo struct {
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	TermsOfService string         `json:"termsOfService"`
	Contact        SwaggerContact `json:"contact"`
	License        SwaggerLicence `json:"license"`
	Version        string         `json:"version"`
}

type SwaggerAPI struct {
	Swagger      string                 `json:"swagger"`
	Info         SwaggerInfo            `json:"info"`
	Host         string                 `json:"host"`
	BasePath     string                 `json:"basePath"`
	Tags         []SwaggerTag           `json:"tags"`
	Paths        map[string]SwaggerPath `json:"paths"`
	ExternalDocs SwaggerExternalDocs    `json:"externalDocs"`
	Definitions  map[string]SwaggerDefn `json:"definitions"`
}

type SwaggerFieldProp struct {
	Type    string `json:"type,omitempty"`
	Format  string `json:"format,omitempty"`
	Example string `json:"example,omitempty"`
}

type SwaggerDefn struct {
	Type       string                      `json:"type,omitempty"`
	Properties map[string]SwaggerFieldProp `json:"properties,omitempty"`
	//Items      SwaggerFieldProp `json:"items,omitempty"`
}

type SwaggerRef struct {
	Ref string `json:"$ref,omitempty"`
}

type SwaggerOperationSchema struct {
	Ref   string     `json:"$ref,omitempty"`
	Type  string     `json:"type,omitempty"`
	Items SwaggerRef `json:"items,omitempty"`
}

type SwaggerOperationParameter struct {
	In          string                 `json:"in,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Required    bool                   `json:"required,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Schema      SwaggerOperationSchema `json:"schema,omitempty"`
}

type SwaggerOperation struct {
	Tags        []string                          `json:"tags,omitempty"`
	Summary     string                            `json:"summary,omitempty"`
	Description string                            `json:"description,omitempty"`
	OperationID string                            `json:"operationId,omitempty"`
	Parameters  []SwaggerOperationParameter       `json:"parameters"`
	Responses   map[int]SwaggerOperationParameter `json:"responses"`
}

type SwaggerPath struct {
	Get    *SwaggerOperation `json:"get,omitempty"`
	Post   *SwaggerOperation `json:"post,omitempty"`
	Put    *SwaggerOperation `json:"put,omitempty"`
	Delete *SwaggerOperation `json:"delete,omitempty"`
}

type SwaggerExternalDocs struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type SwaggerTag struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ExternalDocs SwaggerExternalDocs `json:"externalDocs"`
}

func (a *App) GenerateAPIDocJSON() {
	swaggerDocPath := a.Config.SwaggerDocFolder + "/doc.json"
	os.MkdirAll(path.Dir(swaggerDocPath), 0755)
	file, err := os.Create(swaggerDocPath)
	if err != nil {
		fmt.Println(err)
	}
	var SwaggerAPI SwaggerAPI
	SwaggerAPI.Swagger = "2.0"
	SwaggerAPI.Info.Title = "SwaggerAPI.Info.Title"
	SwaggerAPI.Info.Description = "SwaggerAPI.Info.Description"
	SwaggerAPI.Info.TermsOfService = "SwaggerAPI.Info.TermsOfService"
	SwaggerAPI.Info.Contact.Email = "SwaggerAPI.Info.Contact.Email"
	SwaggerAPI.Info.License.Name = "SwaggerAPI.Info.License.Name"
	SwaggerAPI.Info.License.URL = "SwaggerAPI.Info.License.URL"
	SwaggerAPI.Info.Version = "SwaggerAPI.Info.Version"
	SwaggerAPI.BasePath = a.Config.APIPrefix
	SwaggerAPI.Host = "boilerplate.example.com"
	SwaggerAPI.ExternalDocs.Description = "SwaggerAPI.ExternalDocs.Description"
	SwaggerAPI.ExternalDocs.URL = "SwaggerAPI.ExternalDocs.URL"
	SwaggerAPI.Tags = []SwaggerTag{}
	SwaggerAPI.Definitions = map[string]SwaggerDefn{}
	var SwaggerPaths map[string]SwaggerPath = map[string]SwaggerPath{}
	for _, v := range a.Routes {
		var SwaggerPath SwaggerPath
		v.Path = strings.ReplaceAll(v.Path, SwaggerAPI.BasePath, "")
		if PrevSwaggerPath, exists := SwaggerPaths[v.Path]; exists {
			SwaggerPath = PrevSwaggerPath
		}
		SwaggerOperation := SwaggerOperation{
			Tags: []string{strings.Title(strings.Split(v.Path, "/")[1])},
		}
		if v.IPStruct.Name != "" {
			SwaggerOperation.Parameters = append(SwaggerOperation.Parameters, SwaggerOperationParameter{
				In:          "body",
				Description: v.IPStruct.Name,
				Schema:      SwaggerOperationSchema{Ref: "#/definitions/" + v.IPStruct.Name},
			})
			var props map[string]SwaggerFieldProp = map[string]SwaggerFieldProp{}
			for _, v := range v.IPStruct.Fields {
				props[v.TSName] = SwaggerFieldProp{
					Type: v.TSType,
				}
			}
			SwaggerAPI.Definitions[v.IPStruct.Name] = SwaggerDefn{
				Type:       "object",
				Properties: props,
			}
		}
		if v.OPStruct.Name != "" {
			var Responses map[int]SwaggerOperationParameter = map[int]SwaggerOperationParameter{}
			Responses[200] = SwaggerOperationParameter{
				Description: v.OPStruct.Name,
				Schema:      SwaggerOperationSchema{Ref: "#/definitions/" + v.OPStruct.Name},
			}
			var props map[string]SwaggerFieldProp = map[string]SwaggerFieldProp{}
			for _, v := range v.OPStruct.Fields {
				props[v.TSName] = SwaggerFieldProp{
					Type: v.TSType,
				}
			}
			SwaggerAPI.Definitions[v.OPStruct.Name] = SwaggerDefn{
				Type:       "object",
				Properties: props,
			}
			SwaggerOperation.Responses = Responses
		}
		switch v.Method {
		case "GET":
			SwaggerPath.Get = &SwaggerOperation
		case "POST":
			SwaggerPath.Post = &SwaggerOperation
		case "PUT":
			SwaggerPath.Put = &SwaggerOperation
		case "DELETE":
			SwaggerPath.Delete = &SwaggerOperation
		}
		SwaggerPaths[v.Path] = SwaggerPath
	}
	SwaggerAPI.Paths = SwaggerPaths
	op, err := json.Marshal(SwaggerAPI)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(string(op))
}
