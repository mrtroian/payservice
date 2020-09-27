package api

import "testing"

func TestReadID(t *testing.T) {
	api := &API{pattern: "host/"}

	if got, err := api.readID("/"); err == nil {
		t.Errorf(`api.readID("/") = %d; want err`, got)
	}

	if got, err := api.readID(""); err == nil {
		t.Errorf(`api.readID("") = %d; want err`, got)
		t.Fail()
	}

	if got, err := api.readID("42"); err == nil && got != 42 {
		t.Errorf(`api.readID("42") = %d; want 42`, got)
		t.Fail()
	}

	if got, err := api.readID("/42"); err == nil {
		t.Errorf(`api.readID("/42") = %d; want err`, got)
		t.Fail()
	}

	api.pattern = "/"

	if got, err := api.readID("/"); err == nil {
		t.Errorf(`api.readID("/") = %d; want err`, got)
	}

	if got, err := api.readID(""); err == nil {
		t.Errorf(`api.readID("") = %d; want err`, got)
		t.Fail()
	}

	if got, err := api.readID("42"); err == nil && got != 42 {
		t.Errorf(`api.readID("42") = %d; want 42`, got)
		t.Fail()
	}

	if got, err := api.readID("/42"); err != nil {
		t.Errorf(`api.readID("/42") = %d; want 42`, got)
		t.Fail()
	}
}
