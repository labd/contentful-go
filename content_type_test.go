package contentful

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleContentTypesService_Get() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(env, "content-type-id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(contentType.Name)
}

func ExampleContentTypesService_List() {
	cma := NewCMA("cma-token")

	collection, err := cma.ContentTypes.List(env).Next()
	if err != nil {
		log.Fatal(err)
	}

	contentTypes := collection.ToContentType()

	for _, contentType := range contentTypes {
		fmt.Println(contentType.Sys.ID, contentType.Sys.PublishedAt)
	}
}

func ExampleContentTypesService_Upsert_create() {
	cma := NewCMA("cma-token")

	contentType := &ContentType{
		Name:         "test content type",
		DisplayField: "field1_id",
		Description:  "content type description",
		Fields: []*Field{
			{
				ID:       "field1_id",
				Name:     "field1",
				Type:     "Symbol",
				Required: false,
				Disabled: false,
			},
			{
				ID:       "field2_id",
				Name:     "field2",
				Type:     "Symbol",
				Required: false,
				Disabled: true,
			},
		},
	}

	err := cma.ContentTypes.Upsert(env, contentType)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleContentTypesService_Upsert_Update() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(env, "content-type-id")
	if err != nil {
		log.Fatal(err)
	}

	contentType.Name = "modified content type name"

	err = cma.ContentTypes.Upsert(env, contentType)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleContentTypesService_Activate() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(env, "content-type-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.ContentTypes.Activate(env, contentType)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleContentTypesService_Deactivate() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(env, "content-type-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.ContentTypes.Deactivate(env, contentType)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleContentTypesService_Delete() {
	cma := NewCMA("cma-token")

	contentType, err := cma.ContentTypes.Get(env, "content-type-id")
	if err != nil {
		log.Fatal(err)
	}

	err = cma.ContentTypes.Delete(env, contentType)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleContentTypesService_Delete_allDrafts() {
	cma := NewCMA("cma-token")

	collection, err := cma.ContentTypes.List(env).Next()
	if err != nil {
		log.Fatal(err)
	}

	contentTypes := collection.ToContentType()

	for _, contentType := range contentTypes {
		if contentType.Sys.PublishedAt == "" {
			err := cma.ContentTypes.Delete(env, contentType)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func TestContentTypesServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_types.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.ContentTypes.List(env).Next()
	assertions.Nil(err)
	contentType := collection.ToContentType()
	assertions.Equal(4, len(contentType))
	assertions.Equal("City", contentType[0].Name)
}

func TestContentTypesServiceListActivated(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/public/content_types")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_types.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.ContentTypes.ListActivated(env).Next()
	assertions.Nil(err)
}

func TestContentTypesService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	contentType, err := cma.ContentTypes.Get(env, "63Vgs0BFK0USe4i2mQUGK6")
	assertions.Nil(err)
	assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", contentType.Sys.ID)
}

func TestContentTypesService_Get_2(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")

		checkHeaders(r, assertions)

		w.WriteHeader(400)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.ContentTypes.Get(env, "63Vgs0BFK0USe4i2mQUGK6")
	assertions.NotNil(err)
}

func TestContentTypesServiceActivate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assertions.Nil(err)

	err = cma.ContentTypes.Activate(env, ct)
	assertions.Nil(err)
}

func TestContentTypesServiceDeactivate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6/published")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assertions.Nil(err)

	err = cma.ContentTypes.Deactivate(env, ct)
	assertions.Nil(err)
}

func TestContentTypeSaveForCreate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/ct-name")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("ct-name", payload["name"])
		assertions.Equal("ct-description", payload["description"])

		fields := payload["fields"].([]interface{})
		assertions.Equal(2, len(fields))

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})

		assertions.Equal("field1", field1["id"].(string))
		assertions.Equal("field1-name", field1["name"].(string))
		assertions.Equal("Symbol", field1["type"].(string))

		assertions.Equal("field2", field2["id"].(string))
		assertions.Equal("field2-name", field2["name"].(string))
		assertions.Equal("Symbol", field2["type"].(string))
		assertions.Equal(true, field2["disabled"].(bool))

		assertions.Equal(field1["id"].(string), payload["displayField"])

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:       "field1",
		Name:     "field1-name",
		Type:     "Symbol",
		Required: true,
	}

	field2 := &Field{
		ID:       "field2",
		Name:     "field2-name",
		Type:     "Symbol",
		Disabled: true,
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1, field2},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
	assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assertions.Equal("ct-name", ct.Name)
	assertions.Equal("ct-description", ct.Description)
}

func TestContentTypeSaveForUpdate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("ct-name-updated", payload["name"])
		assertions.Equal("ct-description-updated", payload["description"])

		fields := payload["fields"].([]interface{})
		assertions.Equal(3, len(fields))

		field1 := fields[0].(map[string]interface{})
		field2 := fields[1].(map[string]interface{})
		field3 := fields[2].(map[string]interface{})

		assertions.Equal("field1", field1["id"].(string))
		assertions.Equal("field1-name-updated", field1["name"].(string))
		assertions.Equal("String", field1["type"].(string))

		assertions.Equal("field2", field2["id"].(string))
		assertions.Equal("field2-name-updated", field2["name"].(string))
		assertions.Equal("Integer", field2["type"].(string))
		assertions.Nil(field2["disabled"])

		assertions.Equal("field3", field3["id"].(string))
		assertions.Equal("field3-name", field3["name"].(string))
		assertions.Equal("Date", field3["type"].(string))

		assertions.Equal(field3["id"].(string), payload["displayField"])

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_type-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assertions.Nil(err)

	ct.Name = "ct-name-updated"
	ct.Description = "ct-description-updated"

	field1 := ct.Fields[0]
	field1.Name = "field1-name-updated"
	field1.Type = "String"
	field1.Required = false

	field2 := ct.Fields[1]
	field2.Name = "field2-name-updated"
	field2.Type = "Integer"
	field2.Disabled = false

	field3 := &Field{
		ID:   "field3",
		Name: "field3-name",
		Type: "Date",
	}

	ct.Fields = append(ct.Fields, field3)
	ct.DisplayField = ct.Fields[2].ID

	_ = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
	assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assertions.Equal("ct-name-updated", ct.Name)
	assertions.Equal("ct-description-updated", ct.Description)
	assertions.Equal(2, ct.Sys.Version)
}

func TestContentTypeCreateWithID(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/mycontenttype")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("content_type-updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct := &ContentType{
		Sys: &Sys{
			ID: "mycontenttype",
		},
		Name: "MyContentType",
	}

	_ = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
}

func TestContentTypeDelete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/63Vgs0BFK0USe4i2mQUGK6")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	ct, err := contentTypeFromTestData("content_type.json")
	assertions.Nil(err)

	// delete content type
	err = cma.ContentTypes.Delete(env, ct)
	assertions.Nil(err)
}

func TestContentTypeFieldRef(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/ct-name")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].([]interface{})
		assertions.Equal(1, len(fields))

		field1 := fields[0].(map[string]interface{})
		assertions.Equal("Link", field1["type"].(string))
		validations := field1["validations"].([]interface{})
		assertions.Equal(1, len(validations))
		validation := validations[0].(map[string]interface{})
		linkValidationValue := validation["linkContentType"].([]interface{})
		assertions.Equal(1, len(linkValidationValue))
		assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", linkValidationValue[0].(string))

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test content type
	linkCt, err := contentTypeFromTestData("content_type.json")
	assertions.Nil(err)

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: "Link",
		Validations: []FieldValidation{
			FieldValidationLink{
				LinkContentType: []string{linkCt.Sys.ID},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
}

func TestContentTypeFieldArray(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/ct-name")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].([]interface{})
		assertions.Equal(1, len(fields))

		field1 := fields[0].(map[string]interface{})
		assertions.Equal("Array", field1["type"].(string))

		arrayItemSchema := field1["items"].(map[string]interface{})
		assertions.Equal("Text", arrayItemSchema["type"].(string))

		arrayItemSchemaValidations := arrayItemSchema["validations"].([]interface{})
		validation1 := arrayItemSchemaValidations[0].(map[string]interface{})
		assertions.Equal(true, validation1["unique"].(bool))

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: FieldTypeArray,
		Items: &FieldTypeArrayItem{
			Type: FieldTypeText,
			Validations: []FieldValidation{
				&FieldValidationUnique{
					Unique: true,
				},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
}

func TestContentTypeFieldValidationRangeUniquePredefinedValues(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/ct-name")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].([]interface{})
		assertions.Equal(1, len(fields))

		field1 := fields[0].(map[string]interface{})
		assertions.Equal("Integer", field1["type"].(string))

		validations := field1["validations"].([]interface{})

		// unique validation
		validationUnique := validations[0].(map[string]interface{})
		assertions.Equal(false, validationUnique["unique"].(bool))

		// range validation
		validationRange := validations[1].(map[string]interface{})
		rangeValues := validationRange["range"].(map[string]interface{})
		errorMessage := validationRange["message"].(string)
		assertions.Equal("error message", errorMessage)
		assertions.Equal(float64(20), rangeValues["min"].(float64))
		assertions.Equal(float64(30), rangeValues["max"].(float64))

		// predefined validation
		validationPredefinedValues := validations[2].(map[string]interface{})
		predefinedValues := validationPredefinedValues["in"].([]interface{})
		assertions.Equal(3, len(predefinedValues))
		assertions.Equal("error message 2", validationPredefinedValues["message"].(string))
		assertions.Equal(float64(20), predefinedValues[0].(float64))
		assertions.Equal(float64(21), predefinedValues[1].(float64))
		assertions.Equal(float64(22), predefinedValues[2].(float64))

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:   "field1",
		Name: "field1-name",
		Type: FieldTypeInteger,
		Validations: []FieldValidation{
			&FieldValidationUnique{
				Unique: false,
			},
			&FieldValidationRange{
				Range: &MinMax{
					Min: 20,
					Max: 30,
				},
				ErrorMessage: "error message",
			},
			&FieldValidationPredefinedValues{
				In:           []interface{}{20, 21, 22},
				ErrorMessage: "error message 2",
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
}

func TestContentTypeFieldTypeMedia(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/environments/"+environmentID+"/content_types/ct-name")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)

		fields := payload["fields"].([]interface{})
		assertions.Equal(1, len(fields))

		field1 := fields[0].(map[string]interface{})
		assertions.Equal("Link", field1["type"].(string))
		assertions.Equal("Asset", field1["linkType"].(string))

		validations := field1["validations"].([]interface{})

		// mime type validation
		validationMimeType := validations[0].(map[string]interface{})
		linkMimetypeGroup := validationMimeType["linkMimetypeGroup"].([]interface{})
		assertions.Equal(12, len(linkMimetypeGroup))
		var mimetypes []string
		for _, mimetype := range linkMimetypeGroup {
			mimetypes = append(mimetypes, mimetype.(string))
		}
		assertions.Equal(mimetypes, []string{
			MimeTypeAttachment,
			MimeTypePlainText,
			MimeTypeImage,
			MimeTypeAudio,
			MimeTypeVideo,
			MimeTypeRichText,
			MimeTypePresentation,
			MimeTypeSpreadSheet,
			MimeTypePDF,
			MimeTypeArchive,
			MimeTypeCode,
			MimeTypeMarkup,
		})

		// dimension validation
		validationDimension := validations[1].(map[string]interface{})
		errorMessage := validationDimension["message"].(string)
		assetImageDimensions := validationDimension["assetImageDimensions"].(map[string]interface{})
		widthData := assetImageDimensions["width"].(map[string]interface{})
		heightData := assetImageDimensions["height"].(map[string]interface{})
		widthMin := int(widthData["min"].(float64))
		heightMax := int(heightData["max"].(float64))

		_, ok := widthData["max"].(float64)
		assertions.False(ok)

		_, ok = heightData["min"].(float64)
		assertions.False(ok)

		assertions.Equal("custom error message", errorMessage)
		assertions.Equal(100, widthMin)
		assertions.Equal(300, heightMax)

		// size validation
		validationSize := validations[2].(map[string]interface{})
		sizeData := validationSize["assetFileSize"].(map[string]interface{})
		min := int(sizeData["min"].(float64))
		max := int(sizeData["max"].(float64))
		assertions.Equal(30, min)
		assertions.Equal(400, max)

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	field1 := &Field{
		ID:       "field-id",
		Name:     "media-field",
		Type:     FieldTypeLink,
		LinkType: "Asset",
		Validations: []FieldValidation{
			&FieldValidationMimeType{
				MimeTypes: []string{
					MimeTypeAttachment,
					MimeTypePlainText,
					MimeTypeImage,
					MimeTypeAudio,
					MimeTypeVideo,
					MimeTypeRichText,
					MimeTypePresentation,
					MimeTypeSpreadSheet,
					MimeTypePDF,
					MimeTypeArchive,
					MimeTypeCode,
					MimeTypeMarkup,
				},
			},
			&FieldValidationDimension{
				Width: &MinMax{
					Min: 100,
				},
				Height: &MinMax{
					Max: 300,
				},
				ErrorMessage: "custom error message",
			},
			&FieldValidationFileSize{
				Size: &MinMax{
					Min: 30,
					Max: 400,
				},
			},
		},
	}

	ct := &ContentType{
		Name:         "ct-name",
		Description:  "ct-description",
		Fields:       []*Field{field1},
		DisplayField: field1.ID,
	}

	err = cma.ContentTypes.Upsert(env, ct)
	assertions.Nil(err)
}

func TestContentTypeFieldValidationsUnmarshal(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("content_type_with_validations.json"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	ct, err := cma.ContentTypes.Get(env, "validationsTest")
	assertions.Nil(err)

	var uniqueValidations []FieldValidation
	var linkValidations []FieldValidation
	var sizeValidations []FieldValidation
	var regexValidations []FieldValidation
	var preDefinedValidations []FieldValidation
	var rangeValidations []FieldValidation
	var dateValidations []FieldValidation
	var mimeTypeValidations []FieldValidation
	var dimensionValidations []FieldValidation
	var fileSizeValidations []FieldValidation

	for _, field := range ct.Fields {
		if field.Name == "text-short" {
			assertions.Equal(4, len(field.Validations))
			uniqueValidations = append(uniqueValidations, field.Validations[0])
			sizeValidations = append(sizeValidations, field.Validations[1])
			regexValidations = append(regexValidations, field.Validations[2])
			preDefinedValidations = append(preDefinedValidations, field.Validations[3])
		}

		if field.Name == "text-long" {
			assertions.Equal(3, len(field.Validations))
			sizeValidations = append(sizeValidations, field.Validations[0])
			regexValidations = append(regexValidations, field.Validations[1])
			preDefinedValidations = append(preDefinedValidations, field.Validations[2])
		}

		if field.Name == "number-integer" || field.Name == "number-decimal" {
			assertions.Equal(3, len(field.Validations))
			uniqueValidations = append(uniqueValidations, field.Validations[0])
			rangeValidations = append(rangeValidations, field.Validations[1])
			preDefinedValidations = append(preDefinedValidations, field.Validations[2])
		}

		if field.Name == "date" {
			assertions.Equal(1, len(field.Validations))
			dateValidations = append(dateValidations, field.Validations[0])
		}

		if field.Name == "location" || field.Name == "bool" {
			assertions.Equal(0, len(field.Validations))
		}

		if field.Name == "media-onefile" {
			assertions.Equal(3, len(field.Validations))
			mimeTypeValidations = append(mimeTypeValidations, field.Validations[0])
			dimensionValidations = append(dimensionValidations, field.Validations[1])
			fileSizeValidations = append(fileSizeValidations, field.Validations[2])
		}

		if field.Name == "media-manyfiles" {
			assertions.Equal(1, len(field.Validations))
			assertions.Equal(3, len(field.Items.Validations))
			sizeValidations = append(sizeValidations, field.Validations[0])
			mimeTypeValidations = append(mimeTypeValidations, field.Items.Validations[0])
			dimensionValidations = append(dimensionValidations, field.Items.Validations[1])
			fileSizeValidations = append(fileSizeValidations, field.Items.Validations[2])
		}

		if field.Name == "json" {
			assertions.Equal(1, len(field.Validations))
			sizeValidations = append(sizeValidations, field.Validations[0])
		}

		if field.Name == "ref-onref" {
			assertions.Equal(1, len(field.Validations))
			linkValidations = append(linkValidations, field.Validations[0])
		}

		if field.Name == "ref-manyRefs" {
			assertions.Equal(1, len(field.Validations))
			assertions.Equal(1, len(field.Items.Validations))
			linkValidations = append(linkValidations, field.Items.Validations[0])
			sizeValidations = append(sizeValidations, field.Validations[0])
		}
	}

	for _, validation := range uniqueValidations {
		_, ok := validation.(FieldValidationUnique)
		assertions.True(ok)
	}

	for _, validation := range linkValidations {
		_, ok := validation.(FieldValidationLink)
		assertions.True(ok)
	}

	for _, validation := range sizeValidations {
		_, ok := validation.(FieldValidationSize)
		assertions.True(ok)
	}

	for _, validation := range regexValidations {
		_, ok := validation.(FieldValidationRegex)
		assertions.True(ok)
	}

	for _, validation := range preDefinedValidations {
		_, ok := validation.(FieldValidationPredefinedValues)
		assertions.True(ok)
	}

	for _, validation := range rangeValidations {
		_, ok := validation.(FieldValidationRange)
		assertions.True(ok)
	}

	for _, validation := range dateValidations {
		_, ok := validation.(FieldValidationDate)
		assertions.True(ok)
	}

	for _, validation := range mimeTypeValidations {
		_, ok := validation.(FieldValidationMimeType)
		assertions.True(ok)
	}

	for _, validation := range dimensionValidations {
		_, ok := validation.(FieldValidationDimension)
		assertions.True(ok)
	}

	for _, validation := range fileSizeValidations {
		_, ok := validation.(FieldValidationFileSize)
		assertions.True(ok)
	}
}
