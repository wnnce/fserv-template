package sqlbuild

import (
	"reflect"
	"testing"
)

func TestSliceToAnySlice_Int(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []any{1, 2, 3}
	result := SliceToAnySlice(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSliceToAnySlice_String(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := []any{"a", "b", "c"}
	result := SliceToAnySlice(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSliceToAnySlice_Empty(t *testing.T) {
	input := []int{}
	expected := []any{}
	result := SliceToAnySlice(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSliceToAnySlice_Nil(t *testing.T) {
	var input []int = nil
	result := SliceToAnySlice(input)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
