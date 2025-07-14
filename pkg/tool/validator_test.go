package tool

import (
	"testing"
)

type testStruct struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=0,lte=130"`
}

func TestStructValidator_Validate(t *testing.T) {
	v := NewStruckValidator()

	// valid data
	data := testStruct{Name: "Tom", Age: 18}
	err := v.Validate(data)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// missing required field
	data = testStruct{Name: "", Age: 18}
	err = v.Validate(data)
	if err == nil {
		t.Error("expected error for missing name, got nil")
	}

	// age out of range
	data = testStruct{Name: "Tom", Age: 200}
	err = v.Validate(data)
	if err == nil {
		t.Error("expected error for age out of range, got nil")
	}
}

func TestStructValidator_Engine(t *testing.T) {
	v := NewStruckValidator()
	if v.Engine() == nil {
		t.Error("expected non-nil engine")
	}
}

func TestValidatorSingleton(t *testing.T) {
	if Validator() == nil {
		t.Error("expected non-nil default validator")
	}
}
