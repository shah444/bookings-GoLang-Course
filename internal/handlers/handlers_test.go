package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/shah444/bookings-GoLang-Course/internal/driver"
	"github.com/shah444/bookings-GoLang-Course/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "24-05-2022"},
	// 	{key: "end", value: "25-05-2022"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "24-05-2022"},
	// 	{key: "end", value: "25-05-2022"},
	// }, http.StatusOK},
	// {"make-reservation-post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Vidit"},
	// 	{key: "last_name", value: "Shah"},
	// 	{key: "email", value: "test@example.com"},
	// 	{key: "phone", value: "111-111-1111"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// testing the case when there is no reservation in context
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing the case when the room requested is invalid/non-existent
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "Vidit")
	postedData.Add("last_name", "Shah")
	postedData.Add("email", "vid@it.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// testing for missing reservation in session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for missing body in the post request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for invalid form
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "V") // First name should be atleast 3 chars long
	postedData.Add("last_name", "Shah")
	postedData.Add("email", "vid@it.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// testing for error during insertion of reservation into database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "Vidit")
	postedData.Add("last_name", "Shah")
	postedData.Add("email", "vid@it.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	reservation = models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for error during insertion of room restriction into the database
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "Vidit")
	postedData.Add("last_name", "Shah")
	postedData.Add("email", "vid@it.com")
	postedData.Add("phone", "123456789")
	postedData.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	reservation = models.Reservation{
		RoomID: 1000,
		Room: models.Room{
			ID:       1000,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("got availability when none should be there")
	}

	// testing for invalid form
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("got OK = true but expected false")
	}

	// testing for when availability is there
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if !j.OK {
		t.Error("got no availability when it should be there")
	}

	// testing for when start date is invalid
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("got valid date when it should be invalid")
	}

	// testing for when roomID is invalid
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("got valid room when it should be invalid")
	}

	// testing for database error
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("got valid room when it should be invalid")
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("error choosing room, got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// testing when roomID can't be parsed to int
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	req.RequestURI = "/choose-room/invalid"

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("can't parse roomID to int, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing when reservation not there in context
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("error choosing room, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_ReservationSummary(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got incorrect status code, expected %d, got %d", rr.Code, http.StatusOK)
	}

	// testing for the case when reservation is not in context
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, expected %d, got %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// Empty post body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Invalid dates
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// database query error
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// rooms available
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	req, _ := http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// testing for invalid dates
	req, _ = http.NewRequest("GET", "/book-room?s=invalid&e=2040-01-02&id=1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//testing for database query error
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=1000", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("got incorrect status code, got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
