package contentful

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFieldValidationLink(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &FieldValidationLink{
		LinkContentType: []string{"test", "test2"},
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"linkContentType\":[\"test\",\"test2\"]}", string(data))
}

func TestFieldValidationUnique(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &FieldValidationUnique{
		Unique: false,
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"unique\":false}", string(data))
}

func TestFieldValidationPredefinedValues(t *testing.T) {
	var err error
	assertions := assert.New(t)

	validation := &FieldValidationPredefinedValues{
		In:           []interface{}{5, 10, "string", 6.4},
		ErrorMessage: "error message",
	}

	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"in\":[5,10,\"string\",6.4],\"message\":\"error message\"}", string(data))
}

func TestFieldValidationRange(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// between
	validation := &FieldValidationRange{
		Range: &MinMax{
			Min: 60,
			Max: 100,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"min\":60,\"max\":100},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationRange
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(60), validationCheck.Range.Min)
	assertions.Equal(float64(100), validationCheck.Range.Max)
	assertions.Equal("error message", validationCheck.ErrorMessage)

	// greater than equal to
	validation = &FieldValidationRange{
		Range: &MinMax{
			Min: 10,
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"min\":10},\"message\":\"error message\"}", string(data))
	validationCheck = FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(10), validationCheck.Range.Min)
	assertions.Equal(float64(0), validationCheck.Range.Max)
	assertions.Equal("error message", validationCheck.ErrorMessage)

	// less than equal to
	validation = &FieldValidationRange{
		Range: &MinMax{
			Max: 90,
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"range\":{\"max\":90},\"message\":\"error message\"}", string(data))
	validationCheck = FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(90), validationCheck.Range.Max)
	assertions.Equal(float64(0), validationCheck.Range.Min)
	assertions.Equal("error message", validationCheck.ErrorMessage)
}

func TestFieldValidationSize(t *testing.T) {
	var err error
	assertions := assert.New(t)

	// between
	validation := &FieldValidationSize{
		Size: &MinMax{
			Min: 4,
			Max: 6,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"size\":{\"min\":4,\"max\":6},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationSize
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(float64(4), validationCheck.Size.Min)
	assertions.Equal(float64(6), validationCheck.Size.Max)
	assertions.Equal("error message", validationCheck.ErrorMessage)
}

func TestFieldValidationDate(t *testing.T) {
	var err error
	assertions := assert.New(t)

	layout := "2006-01-02T03:04:05"
	min := time.Now()
	max := time.Now()

	minStr := min.Format(layout)
	maxStr := max.Format(layout)

	validation := &FieldValidationDate{
		Range: &DateMinMax{
			Min: min,
			Max: max,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	assertions.Nil(err)
	assertions.Equal("{\"dateRange\":{\"min\":\""+minStr+"\",\"max\":\""+maxStr+"\"},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationDate
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	assertions.Nil(err)
	assertions.Equal(minStr, validationCheck.Range.Min.Format(layout))
	assertions.Equal(maxStr, validationCheck.Range.Max.Format(layout))
	assertions.Equal("error message", validationCheck.ErrorMessage)
}
