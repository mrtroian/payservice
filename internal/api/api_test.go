package api

import "testing"

func TestReadID(t *testing.T) {
	if got, err := ReadID("/"); err == nil {
		t.Errorf(`ReadID("/") = %d; want err`, got)
	}

	if got, err := ReadID(""); err == nil {
		t.Errorf(`ReadID("") = %d; want err`, got)
		t.Fail()
	}

	if got, err := ReadID("42"); err == nil && got != 42 {
		t.Errorf(`ReadID("42") = %d; want 42`, got)
		t.Fail()
	}

	if got, err := ReadID("/42"); err == nil {
		t.Errorf(`ReadID("/42") = %d; want err`, got)
		t.Fail()
	}

	apiEndpoint = "/"

	if got, err := ReadID("/"); err == nil {
		t.Errorf(`ReadID("/") = %d; want err`, got)
	}

	if got, err := ReadID(""); err == nil {
		t.Errorf(`ReadID("") = %d; want err`, got)
		t.Fail()
	}

	if got, err := ReadID("42"); err == nil && got != 42 {
		t.Errorf(`ReadID("42") = %d; want 42`, got)
		t.Fail()
	}

	if got, err := ReadID("/42"); err != nil {
		t.Errorf(`ReadID("/42") = %d; want 42`, got)
		t.Fail()
	}
}
