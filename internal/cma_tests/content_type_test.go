package cma_tests

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/flaconi/contentful-go/internal/testutil"
	"github.com/flaconi/contentful-go/pkgs/common"
	"github.com/flaconi/contentful-go/pkgs/model"
	"github.com/flaconi/contentful-go/pkgs/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentTypesService_List(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/content_type/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types", r.URL.Path)
	})

	defer ts.Close()

	contentTypes, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().List(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(contentTypes.Items, 4)
	assertions.Equal("1t9IbcfdCk6m04uISSsaIK", contentTypes.Items[0].Sys.ID)
	assertions.Equal("City", contentTypes.Items[0].Name)
}

func TestContentTypesService_ListActivated(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/content_type/list.json"}, nil, func(r *http.Request) {
		assertions.Equal("GET", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/public/content_types", r.URL.Path)
	})

	defer ts.Close()

	contentTypes, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().ListActivated(context.Background()).Next()
	assertions.Nil(err)
	assertions.Len(contentTypes.Items, 4)
	assertions.Equal("1t9IbcfdCk6m04uISSsaIK", contentTypes.Items[0].Sys.ID)
	assertions.Equal("City", contentTypes.Items[0].Name)
}

func TestContentTypesService_Get(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		resultValidation func(assertions *assert.Assertions, ct *model.ContentType, err error)
		path             string
		name             string
		statusCode       int
	}{
		{
			resultValidation: func(assertions *assert.Assertions, contentType *model.ContentType, err error) {
				assertions.Nil(err)
				assertions.Equal("ct-name", contentType.Name)
				assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", contentType.Sys.ID)
			},
			path:       "/content_type/get.json",
			name:       "found",
			statusCode: 200,
		},
		{
			resultValidation: func(assertions *assert.Assertions, ct *model.ContentType, err error) {
				assertions.NotNil(err)
				var contentfulError common.NotFoundError
				assertions.True(errors.As(err, &contentfulError))
			},
			path:       "/content_type/not_found.json",
			statusCode: 404,
			name:       "not found",
		},
		{
			path:       "/content_type/validations.json",
			statusCode: 200,
			name:       "validations",
			resultValidation: func(assertions *assert.Assertions, ct *model.ContentType, err error) {
				assertions.Nil(err)

				var uniqueValidations []model.FieldValidation
				var linkValidations []model.FieldValidation
				var sizeValidations []model.FieldValidation
				var regexValidations []model.FieldValidation
				var preDefinedValidations []model.FieldValidation
				var rangeValidations []model.FieldValidation
				var dateValidations []model.FieldValidation
				var mimeTypeValidations []model.FieldValidation
				var dimensionValidations []model.FieldValidation
				var fileSizeValidations []model.FieldValidation

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
					_, ok := validation.(model.FieldValidationUnique)
					assertions.True(ok)
				}

				for _, validation := range linkValidations {
					_, ok := validation.(model.FieldValidationLink)
					assertions.True(ok)
				}

				for _, validation := range sizeValidations {
					_, ok := validation.(model.FieldValidationSize)
					assertions.True(ok)
				}

				for _, validation := range regexValidations {
					_, ok := validation.(model.FieldValidationRegex)
					assertions.True(ok)
				}

				for _, validation := range preDefinedValidations {
					_, ok := validation.(model.FieldValidationPredefinedValues)
					assertions.True(ok)
				}

				for _, validation := range rangeValidations {
					_, ok := validation.(model.FieldValidationRange)
					assertions.True(ok)
				}

				for _, validation := range dateValidations {
					_, ok := validation.(model.FieldValidationDate)
					assertions.True(ok)
				}

				for _, validation := range mimeTypeValidations {
					_, ok := validation.(model.FieldValidationMimeType)
					assertions.True(ok)
				}

				for _, validation := range dimensionValidations {
					_, ok := validation.(model.FieldValidationDimension)
					assertions.True(ok)
				}

				for _, validation := range fileSizeValidations {
					_, ok := validation.(model.FieldValidationFileSize)
					assertions.True(ok)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: tt.statusCode, Path: tt.path}, nil, func(r *http.Request) {
				assertions.Equal("GET", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/63Vgs0BFK0USe4i2mQUGK6", r.URL.Path)
			})

			defer ts.Close()

			contentType, err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Get(context.Background(), "63Vgs0BFK0USe4i2mQUGK6")
			tt.resultValidation(assertions, contentType, err)
		})
	}
}

func TestContentTypesService_Upsert_Create(t *testing.T) {
	assertions := assert.New(t)

	var tests = []struct {
		payloadValidation func(payload map[string]interface{}, assertions *assert.Assertions)
		contentType       *model.ContentType
		name              string
	}{
		{
			name: "base",
			payloadValidation: func(payload map[string]interface{}, assertions *assert.Assertions) {
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
			},
			contentType: &model.ContentType{
				Name:        "ct-name",
				Description: util.ToPointer("ct-description"),
				Fields: []*model.Field{{
					ID:       "field1",
					Name:     "field1-name",
					Type:     "Symbol",
					Required: true,
				}, {
					ID:       "field2",
					Name:     "field2-name",
					Type:     "Symbol",
					Disabled: true,
				}},
				DisplayField: "field1",
			},
		},
		{
			name: "validation link",
			payloadValidation: func(payload map[string]interface{}, assertions *assert.Assertions) {
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
			},
			contentType: &model.ContentType{
				Name:        "ct-name",
				Description: util.ToPointer("ct-description"),
				Fields: []*model.Field{{
					ID:   "field1",
					Name: "field1-name",
					Type: model.FieldTypeLink,
					Validations: []model.FieldValidation{
						model.FieldValidationLink{
							LinkContentType: []string{"63Vgs0BFK0USe4i2mQUGK6"},
						},
					},
				}},
				DisplayField: "field1",
			},
		},
		{
			name: "validation array",
			payloadValidation: func(payload map[string]interface{}, assertions *assert.Assertions) {
				fields := payload["fields"].([]interface{})
				assertions.Equal(1, len(fields))

				field1 := fields[0].(map[string]interface{})
				assertions.Equal("Array", field1["type"].(string))

				arrayItemSchema := field1["items"].(map[string]interface{})
				assertions.Equal("Text", arrayItemSchema["type"].(string))

				arrayItemSchemaValidations := arrayItemSchema["validations"].([]interface{})
				validation1 := arrayItemSchemaValidations[0].(map[string]interface{})
				assertions.Equal(true, validation1["unique"].(bool))
			},
			contentType: &model.ContentType{
				Name:        "ct-name",
				Description: util.ToPointer("ct-description"),
				Fields: []*model.Field{{
					ID:   "field1",
					Name: "field1-name",
					Type: model.FieldTypeArray,
					Items: &model.FieldTypeArrayItem{
						Type: model.FieldTypeText,
						Validations: []model.FieldValidation{
							&model.FieldValidationUnique{
								Unique: true,
							},
						},
					},
				}},
				DisplayField: "field1",
			},
		},
		{
			name: "validation range unique predefinedValues",
			payloadValidation: func(payload map[string]interface{}, assertions *assert.Assertions) {
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
			},
			contentType: &model.ContentType{
				Name:        "ct-name",
				Description: util.ToPointer("ct-description"),
				Fields: []*model.Field{{
					ID:   "field1",
					Name: "field1-name",
					Type: model.FieldTypeInteger,
					Validations: []model.FieldValidation{
						&model.FieldValidationUnique{
							Unique: false,
						},
						&model.FieldValidationRange{
							Range: &model.MinMax{
								Min: util.ToPointer(float64(20)),
								Max: util.ToPointer(float64(30)),
							},
							ErrorMessage: "error message",
						},
						&model.FieldValidationPredefinedValues{
							In:           []interface{}{20, 21, 22},
							ErrorMessage: util.ToPointer("error message 2"),
						},
					},
				}},
				DisplayField: "field1",
			},
		},
		{
			name: "validation media",
			payloadValidation: func(payload map[string]interface{}, assertions *assert.Assertions) {
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
					model.MimeTypeAttachment,
					model.MimeTypePlainText,
					model.MimeTypeImage,
					model.MimeTypeAudio,
					model.MimeTypeVideo,
					model.MimeTypeRichText,
					model.MimeTypePresentation,
					model.MimeTypeSpreadSheet,
					model.MimeTypePDF,
					model.MimeTypeArchive,
					model.MimeTypeCode,
					model.MimeTypeMarkup,
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
				assertions.Equal(30, int(sizeData["min"].(float64)))
				assertions.Equal(400, int(sizeData["max"].(float64)))
			},
			contentType: &model.ContentType{
				Name:        "ct-name",
				Description: util.ToPointer("ct-description"),
				Fields: []*model.Field{{
					ID:       "field1",
					Name:     "field1-name",
					Type:     model.FieldTypeLink,
					LinkType: "Asset",
					Validations: []model.FieldValidation{
						&model.FieldValidationMimeType{
							MimeTypes: []string{
								model.MimeTypeAttachment,
								model.MimeTypePlainText,
								model.MimeTypeImage,
								model.MimeTypeAudio,
								model.MimeTypeVideo,
								model.MimeTypeRichText,
								model.MimeTypePresentation,
								model.MimeTypeSpreadSheet,
								model.MimeTypePDF,
								model.MimeTypeArchive,
								model.MimeTypeCode,
								model.MimeTypeMarkup,
							},
						},
						&model.FieldValidationDimension{
							Width: &model.MinMax{
								Min: util.ToPointer(float64(100)),
							},
							Height: &model.MinMax{
								Max: util.ToPointer(float64(300)),
							},
							ErrorMessage: "custom error message",
						},
						&model.FieldValidationFileSize{
							Size: &model.MinMax{
								Min: util.ToPointer(float64(30)),
								Max: util.ToPointer(float64(400)),
							},
						},
					},
				}},
				DisplayField: "field1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 201, Path: "/content_type/get.json"}, nil, func(r *http.Request) {
				assertions.Equal("PUT", r.Method)
				assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/ct-name", r.URL.Path)

				var payload map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&payload)
				assertions.Nil(err)
				tt.payloadValidation(payload, assertions)
			})

			defer ts.Close()

			ct := tt.contentType

			err := cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Upsert(context.Background(), ct)
			assertions.Nil(err)
			assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
			assertions.Equal("ct-name", ct.Name)
			assertions.Equal("ct-description", *ct.Description)
		})
	}
}

func TestContentTypesService_Upsert_Update(t *testing.T) {
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/content_type/update.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/63Vgs0BFK0USe4i2mQUGK6", r.URL.Path)

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

	})

	defer ts.Close()

	var ct *model.ContentType
	err := testutil.ModelFromTestData("/content_type/get.json", &ct)
	assertions.Nil(err)

	ct.Name = "ct-name-updated"
	ct.Description = util.ToPointer("ct-description-updated")

	field1 := ct.Fields[0]
	field1.Name = "field1-name-updated"
	field1.Type = "String"
	field1.Required = false

	field2 := ct.Fields[1]
	field2.Name = "field2-name-updated"
	field2.Type = "Integer"
	field2.Disabled = false

	field3 := &model.Field{
		ID:   "field3",
		Name: "field3-name",
		Type: "Date",
	}

	ct.Fields = append(ct.Fields, field3)
	ct.DisplayField = ct.Fields[2].ID

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Upsert(context.Background(), ct)

	assertions.Nil(err)
	assertions.Equal("63Vgs0BFK0USe4i2mQUGK6", ct.Sys.ID)
	assertions.Equal("ct-name-updated", ct.Name)
	assertions.Equal("ct-description-updated", *ct.Description)
	assertions.Equal(2, ct.Sys.Version)
}

func TestContentTypesService_Delete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 204, Path: ""}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/63Vgs0BFK0USe4i2mQUGK6", r.URL.Path)
	})

	defer ts.Close()

	var key *model.ContentType
	err = testutil.ModelFromTestData("/content_type/get.json", &key)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Delete(context.Background(), key)
	assertions.Nil(err)
}

func TestContentTypesService_Activate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/content_type/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("PUT", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/63Vgs0BFK0USe4i2mQUGK6/published", r.URL.Path)
	})

	defer ts.Close()

	var key *model.ContentType
	err = testutil.ModelFromTestData("/content_type/get.json", &key)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Activate(context.Background(), key)
	assertions.Nil(err)
}

func TestContentTypesService_Deactivate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	cma, ts := testutil.MockCMAClient(t, assertions, testutil.ResponseData{StatusCode: 200, Path: "/content_type/get.json"}, nil, func(r *http.Request) {
		assertions.Equal("DELETE", r.Method)
		assertions.Equal("/spaces/"+testutil.SpaceID+"/environments/test/content_types/63Vgs0BFK0USe4i2mQUGK6/published", r.URL.Path)
	})

	defer ts.Close()

	var key *model.ContentType
	err = testutil.ModelFromTestData("/content_type/get.json", &key)
	assertions.Nil(err)

	err = cma.WithSpaceId(testutil.SpaceID).WithEnvironment("test").ContentTypes().Deactivate(context.Background(), key)
	assertions.Nil(err)
}
