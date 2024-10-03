package wfiber

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/N4r35h/gos2tsi"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const indent = "    "

type ParsedFiberConstraint struct {
	Name   string
	Params []string
}

type RouteParam struct {
	Name        string
	TSType      string
	IsOptional  bool
	Constraints []ParsedFiberConstraint
}

type Route struct {
	Path      string
	Method    string
	RawIPType interface{}
	RawOPType interface{}
	IPStruct  gos2tsi.ParsedStruct
	OPStruct  gos2tsi.ParsedStruct
	Name      string
	Params    []RouteParam
}

type WFiberAppConfig struct {
	APIPrefix              string
	FrontendFolder         string
	SwaggerDocFolder       string
	GeneratedAPIClientPath string
	GenerateClient         bool
}

type App struct {
	FApp    *fiber.App
	Routes  []Route
	Codegen gos2tsi.Converter
	Config  WFiberAppConfig
}

func New(wfiberConfig WFiberAppConfig, config ...fiber.Config) *App {
	a := &App{
		FApp:    fiber.New(config...),
		Codegen: *gos2tsi.New(),
		Config:  wfiberConfig,
	}
	a.Codegen.Indent = "	"
	return a
}

type Router struct {
	Prefix  string
	App     *App
	FRouter fiber.Router
}

// Just passthrough
func (a *App) Use(args ...interface{}) fiber.Router {
	return a.FApp.Use(args...)
}
func (a *App) Static(prefix string, root string, config ...fiber.Static) fiber.Router {
	return a.FApp.Static(prefix, root, config...)
}
func (a *Router) All(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return a.FRouter.All(path, handlers...)
}

func (a *App) Listen(addr string) error {
	go a.CodeGen()
	return a.FApp.Listen(addr)
}

type CodeGenOutput struct {
	TSAPIClientData TSAPIClientData
}

func (a *App) CodeGen() CodeGenOutput {
	var cgo CodeGenOutput
	if a.Config.GenerateClient {
		start0 := time.Now()
		fmt.Println("CODE GEN START")
		for i, r := range a.Routes {
			if r.RawIPType != nil {
				a.Routes[i].IPStruct = a.Codegen.ParseStruct(r.RawIPType)
			}
			if r.RawOPType != nil {
				a.Routes[i].OPStruct = a.Codegen.ParseStruct(r.RawOPType)
			}
		}
		fmt.Println("CODE GEN Routes -", time.Since(start0))
		start1 := time.Now()
		cgo.TSAPIClientData = a.GenerateTSFile()
		fmt.Println("CODE GEN GenerateTSFile -", time.Since(start1))
		start2 := time.Now()
		a.GenerateAPIDocJSON()
		fmt.Println("CODE GEN GenerateAPIDocJSON -", time.Since(start2))
		fmt.Println("CODE GEN END -", time.Since(start0))
	}
	return cgo
}

type TSAPIClientData struct {
	RawFileContent string
}

func (a *App) GenerateTSFile() TSAPIClientData {
	generatedFilePath := a.Config.FrontendFolder + a.Config.GeneratedAPIClientPath
	os.MkdirAll(path.Dir(generatedFilePath), 0755)
	file, err := os.Create(generatedFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	fileGeneratedContent := "// this file was generated DO NOT EDIT\n"
	// Generate structs
	sortedStructs := make([]gos2tsi.ParsedStruct, 0, len(a.Codegen.Structs))
	for _, v := range a.Codegen.Structs {
		sortedStructs = append(sortedStructs, v)
	}
	sort.Slice(sortedStructs, func(i, j int) bool {
		return sortedStructs[i].Name < sortedStructs[j].Name
	})
	for _, v := range sortedStructs {
		if !v.Required {
			continue
		}
		fileGeneratedContent += a.Codegen.GetStructAsInterfaceString(v) + "\n"
		if !strings.Contains(gos2tsi.GetFormattedInterfaceName(v.Name), "<") {
			fileGeneratedContent += "export interface " + gos2tsi.GetFormattedInterfaceName(v.Name) + "SearchCritCondition {\n"
			fileGeneratedContent += a.Codegen.Indent + "field: "
			for i, f := range v.Fields {
				segment := ""
				if i != 0 {
					segment += " | "
				}
				if f.Var.Embedded() {
					for _, f := range a.Codegen.Structs[f.Var.Pkg().Path()+"."+f.TSName].Fields {
						if segment != "" {
							segment += " | "
						}
						segment += "'" + f.TSName + "'"
					}
				} else {
					segment += "'" + f.TSName + "'"
				}
				fileGeneratedContent += segment
			}
			fileGeneratedContent += "\n" + a.Codegen.Indent + "condition: 'equals' | 'not equals' | 'contains' | 'not contains' | 'greater than' | 'lesser than'"
			fileGeneratedContent += "\n" + a.Codegen.Indent + "value: any\n"
			fileGeneratedContent += "}\n"
		}
	}
	fileGeneratedContent += "import http from './http'\n"
	// Generate api helpers
	for i := range a.Routes {
		fileGeneratedContent += a.GetAPIHelperFuncString(&a.Routes[i], a.Config.APIPrefix)
	}
	file.WriteString(fileGeneratedContent)
	return TSAPIClientData{
		RawFileContent: fileGeneratedContent,
	}
}

func (a *App) GetAPIHelperFuncString(r *Route, apiPrefix string) string {
	toRet := "export const " + GetFunctionNameOfRoute(r, apiPrefix) + " = ("
	for _, p := range r.Params {
		toRet += p.Name
		if p.IsOptional {
			toRet += "?"
		}
		toRet += ": " + p.TSType + ", "
	}
	if r.IPStruct.Name != "" {
		toRet += "_ip: " + r.IPStruct.Name + ", "
	}
	toRet += "query?: string"
	if r.Method == "GET" && len(r.OPStruct.GenericPopulations) > 0 {
		toRet += ", filter?: " + r.OPStruct.GenericPopulations[0].TSType + "SearchCritCondition[]"
	}
	toRet += ")"
	if r.OPStruct.Name != "" {
		toRet += ": Promise<" + getStructNameAfterGenericsPopulation(r.OPStruct)
		for i := 0; i < r.OPStruct.IsSlice; i++ {
			toRet += "[]"
		}
		toRet += ">"
	} else {
		toRet += ": Promise<number>"
	}
	toRet += " => {\n"
	toRet += indent + "return new Promise((resolve, reject) => {\n"
	toRet += indent + indent + "let q: string = query || ''\n"
	if r.Method == "GET" && len(r.OPStruct.GenericPopulations) > 0 {
		toRet += indent + indent + "if (filter && filter.length > 0) {\n"
		toRet += indent + indent + indent + "q = q + (q != '' && '&' || '?') + 'filter=' + encodeURI(JSON.stringify(filter))\n"
		toRet += indent + indent + "}\n"
	}
	toRet += indent + indent + "http." + strings.ToLower(r.Method) + "(" + GetParamsInsertedPath(r) + " + q"
	if r.IPStruct.Name != "" {
		toRet += ", _ip"
	}
	toRet += ")\n"
	toRet += indent + indent + indent + ".then(response => {\n"
	toRet += indent + indent + indent + indent + "return resolve("
	if r.OPStruct.Name != "" {
		toRet += "response.data"
	} else {
		toRet += "response.status"
	}
	toRet += ")\n"
	toRet += indent + indent + indent + "})\n"
	toRet += indent + indent + indent + ".catch(reject)\n"
	toRet += indent + "})\n"
	toRet += "}\n"
	return toRet
}

func getStructNameAfterGenericsPopulation(ps gos2tsi.ParsedStruct) string {
	var toRet string
	toRet = ps.Name
	if len(ps.GenericPopulations) > 0 {
		nameWithOutGenerics := strings.Split(toRet, "[")
		toRet = nameWithOutGenerics[0] + "<"
		for i, v := range ps.GenericPopulations {
			if i != 0 {
				toRet += " "
			}
			toRet += v.TSType
			for i := 0; i < v.IsSlice; i++ {
				toRet += "[]"
			}
		}
		toRet += ">"
	}
	return toRet
}

func (a *App) Group(prefix string, handlers ...fiber.Handler) *Router {
	return &Router{App: a, FRouter: a.FApp.Group(prefix, handlers...), Prefix: prefix}
}

func (a *Router) Group(prefix string, handlers ...fiber.Handler) *Router {
	return &Router{App: a.App, FRouter: a.FRouter.Group(prefix, handlers...), Prefix: a.Prefix + prefix}
}

func (a *Router) Get(prefix string, ip interface{}, op interface{}, handlers ...fiber.Handler) *Route {
	Route := Route{
		Path:      a.Prefix + prefix,
		Method:    fiber.MethodGet,
		RawIPType: ip,
		RawOPType: op,
	}
	a.App.Routes = append(a.App.Routes, Route)
	a.FRouter.Get(prefix, handlers...)
	return &Route
}

func (a *Router) Post(prefix string, ip interface{}, op interface{}, handlers ...fiber.Handler) *Route {
	Route := Route{
		Path:      a.Prefix + prefix,
		Method:    fiber.MethodPost,
		RawIPType: ip,
		RawOPType: op,
	}
	a.App.Routes = append(a.App.Routes, Route)
	a.FRouter.Post(prefix, handlers...)
	return &Route
}

func (a *Router) Put(prefix string, ip interface{}, op interface{}, handlers ...fiber.Handler) *Route {
	Route := Route{
		Path:      a.Prefix + prefix,
		Method:    fiber.MethodPut,
		RawIPType: ip,
		RawOPType: op,
	}
	a.App.Routes = append(a.App.Routes, Route)
	a.FRouter.Put(prefix, handlers...)
	return &Route
}

func (a *Router) Delete(prefix string, ip interface{}, op interface{}, handlers ...fiber.Handler) *Route {
	Route := Route{
		Path:      a.Prefix + prefix,
		Method:    fiber.MethodDelete,
		RawIPType: ip,
		RawOPType: op,
	}
	a.App.Routes = append(a.App.Routes, Route)
	a.FRouter.Delete(prefix, handlers...)
	return &Route
}

func GetFunctionNameOfRoute(r *Route, apiPrefix string) string {
	toRet := strings.Title(strings.ToLower(r.Method))
	r.Path = strings.ReplaceAll(r.Path, apiPrefix, "")
	//nameSeg := strings.ReplaceAll(r.Path, ":", "By")
	for _, v := range strings.Split(r.Path, "/") {
		OpenAngularBrackets := strings.Index(v, "<")
		ClosingAngularBrackets := strings.Index(v, ">")
		if strings.HasPrefix(v, ":") {
			r.Params = append(r.Params, GetFiberParsedRouteParam(v))
		}
		if OpenAngularBrackets != -1 && ClosingAngularBrackets != -1 {
			v = v[:OpenAngularBrackets] + v[ClosingAngularBrackets:]
		}
		v = strings.Title(v)
		v = strings.ReplaceAll(v, ":", "By")
		v = strings.ReplaceAll(v, "?", "")
		toRet += strings.ReplaceAll(v, ">", "")
	}
	return toRet
}

var FIBER_CONSTRAINT_NAME_TO_TS_TYPE = map[string]string{
	"bool":  "boolean",
	"int":   "number",
	"float": "number",
	"min":   "number",
	"max":   "number",
	"range": "number",
}

func GetFiberParsedRouteParam(segment string) RouteParam {
	var rp RouteParam
	if strings.Contains(segment, "?") {
		rp.IsOptional = true
	}
	rawParamName, rawConstraints := SplitURLParamNameAndConstraints(segment)
	rp.Name = rawParamName
	rp.Constraints = ParseAndGetConstraints(rawConstraints)
	rp.TSType = "string"
	for _, v := range rp.Constraints {
		if tsType, exists := FIBER_CONSTRAINT_NAME_TO_TS_TYPE[v.Name]; exists {
			rp.TSType = tsType
			break
		}
	}
	return rp
}

func SplitURLParamNameAndConstraints(ip string) (name string, constraints string) {
	URLParamSegments := strings.Split(ip, "<")
	var constraint string
	if len(URLParamSegments) > 1 {
		constraint = strings.Replace(URLParamSegments[1], ">", "", 1)
	}
	var paramName = strings.Replace(URLParamSegments[0], ":", "", 1)
	paramName = strings.ReplaceAll(paramName, "?", "")
	return paramName, constraint
}

func ParseAndGetConstraints(ip string) []ParsedFiberConstraint {
	var constraints []ParsedFiberConstraint
	for _, v := range strings.Split(ip, ";") {
		constraintSegments := strings.Split(v, "(")
		var parsedConstraint ParsedFiberConstraint
		parsedConstraint.Name = constraintSegments[0]
		if len(constraintSegments) > 1 {
			paramters := strings.Replace(constraintSegments[1], ")", "", 1)
			parsedConstraint.Params = append(parsedConstraint.Params, strings.Split(paramters, ",")...)
		}
		constraints = append(constraints, parsedConstraint)
	}
	return constraints
}

func GetParamsInsertedPath(r *Route) string {
	toRet := ""
	hasParams := false
	for _, seg := range strings.Split(r.Path, "/") {
		if seg == "" {
			continue
		}
		toRet += "/"
		if strings.HasPrefix(seg, ":") {
			hasParams = true
			toRet += "${" + GetFiberParsedRouteParam(seg).Name + "}"
		} else {
			toRet += seg
		}
	}
	if hasParams {
		return "`" + toRet + "`"
	}
	return "'" + toRet + "'"
}
