package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r = httptest.NewRequest("POST", "/", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid even when required fields are present")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.Has("a") {
		t.Error("form shows field present when it is not")
	}

	postedData.Add("a", "a")

	form = New(postedData)

	if !form.Has("a") {
		t.Error("form shows field is missing when it is not")
	}
}

func TestForm_MinLength(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "abc")

	form := New(postedData)

	res := form.MinLength("abc", 4)

	if res {
		t.Error("form shows valid min length for non-existent field")
	}

	if form.Errors.Get("abc") == "" {
		t.Error("did not get error when should have")
	}

	if form.MinLength("a", 4) {
		t.Error("form shows valid min length even when it is invalid")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)
	res = form.MinLength("a", 3)

	if !res {
		t.Error("form shows invalid min length even when it is valid")
	}

	if form.Errors.Get("a") != "" {
		t.Error("form shows error when it should not")
	}
}

func TestForm_IsEmail(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "abc")

	form := New(postedData)

	if form.IsEmail("abc") {
		t.Error("form shows field is email for non-existent field")
	}

	if form.IsEmail("a") {
		t.Error("form shows field is email even when it is not")
	}

	postedData.Add("b", "abc@test.com")
	form = New(postedData)

	if !form.IsEmail("b") {
		t.Error("form shows field is not email even when it is")
	}
}
