package job

import (
	"digi-bot/model"
	"testing"
)

func Test_compare(t *testing.T) {
	newProduct := model.ProductDto{
		Price: 100,
	}
	oldProduct := model.ProductDto{
		Price: 104,
	}

	expected := false
	_, isChanged := compare(newProduct, oldProduct)

	got := isChanged
	if expected != got {
		t.Fatalf("fail. expect: %t, got: %t", expected, got)
	}

}
func Test_compare_greater(t *testing.T) {
	newProduct := model.ProductDto{
		Price: 100,
	}
	oldProduct := model.ProductDto{
		Price: 200,
	}

	expected := true
	_, isChanged := compare(newProduct, oldProduct)

	got := isChanged
	if expected != got {
		t.Fatalf("fail. expect: %t, got: %t", expected, got)
	}

}
