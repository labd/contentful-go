package model_tests

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/labd/contentful-go/pkgs/model"
	"github.com/labd/contentful-go/pkgs/util"

	"github.com/stretchr/testify/assert"
)

func TestFieldValidationLink(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &model.FieldValidationLink{
		LinkContentType: []string{"test", "test2"},
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"linkContentType\":[\"test\",\"test2\"]}", string(data))
}

func TestFieldValidationUnique(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &model.FieldValidationUnique{
		Unique: false,
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"unique\":false}", string(data))
}

func TestFieldValidationPredefinedValues(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &model.FieldValidationPredefinedValues{
		In:           []interface{}{5, 10, "string", 6.4},
		ErrorMessage: util.ToPointer("error message"),
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"in\":[5,10,\"string\",6.4],\"message\":\"error message\"}", string(data))
}

func TestFieldValidationRange(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// between
	validation := &model.FieldValidationRange{
		Range: &model.MinMax{
			Min: util.ToPointer(float64(60)),
			Max: util.ToPointer(float64(100)),
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"min\":60,\"max\":100},\"message\":\"error message\"}", string(data))

	var validationCheck model.FieldValidationRange
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(60), *validationCheck.Range.Min)
	assertions.Equal(float64(100), *validationCheck.Range.Max)
	assertions.Equal("error message", validationCheck.ErrorMessage)

	// greater than equal to
	validation = &model.FieldValidationRange{
		Range: &model.MinMax{
			Min: util.ToPointer(float64(10)),
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"min\":10},\"message\":\"error message\"}", string(data))
	validationCheck = model.FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(10), *validationCheck.Range.Min)
	assertions.Nil(validationCheck.Range.Max)
	assertions.Equal("error message", validationCheck.ErrorMessage)

	// less than equal to
	validation = &model.FieldValidationRange{
		Range: &model.MinMax{
			Max: util.ToPointer(float64(90)),
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"max\":90},\"message\":\"error message\"}", string(data))
	validationCheck = model.FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(90), *validationCheck.Range.Max)
	assertions.Nil(validationCheck.Range.Min)
	assertions.Equal("error message", validationCheck.ErrorMessage)
}

func TestFieldValidationSize(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// between
	validation := &model.FieldValidationSize{
		Size: &model.MinMax{
			Min: util.ToPointer(float64(4)),
			Max: util.ToPointer(float64(6)),
		},
		ErrorMessage: util.ToPointer("error message"),
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"size\":{\"min\":4,\"max\":6},\"message\":\"error message\"}", string(data))

	var validationCheck model.FieldValidationSize
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(4), *validationCheck.Size.Min)
	assertions.Equal(float64(6), *validationCheck.Size.Max)
	assertions.Equal("error message", *validationCheck.ErrorMessage)
}

func TestFieldValidationDate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	layout := "2006-01-02T03:04:05"
	minTime := time.Now()
	maxTime := time.Now()

	minStr := minTime.Format(layout)
	maxStr := maxTime.Format(layout)

	validation := &model.FieldValidationDate{
		Range: &model.DateMinMax{
			Min: minTime,
			Max: maxTime,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"dateRange\":{\"min\":\""+minStr+"\",\"max\":\""+maxStr+"\"},\"message\":\"error message\"}", string(data))

	var validationCheck model.FieldValidationDate
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(minStr, validationCheck.Range.Min.Format(layout))
	assertions.Equal(maxStr, validationCheck.Range.Max.Format(layout))
	assertions.Equal("error message", validationCheck.ErrorMessage)
}
