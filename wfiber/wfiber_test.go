package wfiber

import (
	"strings"
	"testing"

	"github.com/N4r35h/wfiber/wfiber/models_test"
	"github.com/N4r35h/wfiber/wfiber/structs_test"
	"github.com/gofiber/fiber/v2"
)

func TestWFiber(t *testing.T) {
	app := New(WFiberAppConfig{
		APIPrefix:              "/api",
		FrontendFolder:         "frontend_test",
		SwaggerDocFolder:       "swagger_test",
		GeneratedAPIClientPath: "/src/api/api.gen.ts",
		GenerateClient:         true,
	})
	api := app.Group("/api")
	// tests entity
	tests := api.Group("/tests")
	tests.Get("/", nil, structs_test.SimpleAPIResponseWithData[[]models_test.Tests]{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponseWithData[[]models_test.Tests]{})
	})
	tests.Post("/", models_test.Tests{}, structs_test.SimpleAPIResponseWithData[models_test.Tests]{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponseWithData[models_test.Tests]{})
	})
	testswid := tests.Group("/:id")
	testswid.Get("/", nil, structs_test.SimpleAPIResponseWithData[models_test.Tests]{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponseWithData[models_test.Tests]{})
	})
	testswid.Put("/", models_test.Tests{}, structs_test.SimpleAPIResponseWithData[models_test.Tests]{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponseWithData[models_test.Tests]{})
	})
	testswid.Delete("/", nil, structs_test.SimpleAPIResponse{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponse{})
	})
	testswid.Get("/_special", nil, structs_test.SimpleAPIResponse{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponse{})
	})
	// Experimental
	api.Get("/tests1", nil, structs_test.SimpleAPIResponseWithData[models_test.Tests]{}, func(c *fiber.Ctx) error {
		return c.JSON(structs_test.SimpleAPIResponseWithData[models_test.Tests]{})
	})
	api.Get("/tests3", nil, []models_test.Tests{}, func(c *fiber.Ctx) error {
		return c.JSON([]models_test.Tests{})
	})
	cgo := app.CodeGen()

	generatedContent := cgo.TSAPIClientData.RawFileContent
	// Check if "// this file was generated DO NOT EDIT" is set at the top of the file
	if !strings.HasPrefix(generatedContent, "// this file was generated DO NOT EDIT") {
		t.Errorf("Doesnt have generated file dont edit comment!")
	}

	// Check if http.ts module is imported
	if !strings.Contains(generatedContent, "import http from './http'") {
		t.Errorf("Doesnt have http client imported!")
	}

	// Check if generated api client file has all required structs
	var structsToBePresent []interface{} = []interface{}{
		structs_test.SimpleAPIResponseWithData[any]{},
		models_test.Tests{},
	}
	for _, Struct := range structsToBePresent {
		parsedStructMeta := app.Codegen.ParseStruct(Struct)
		if !strings.Contains(generatedContent, app.Codegen.GetStructAsInterfaceString(parsedStructMeta)) {
			t.Errorf(parsedStructMeta.Name + " interface not present!")
		}
	}
}
