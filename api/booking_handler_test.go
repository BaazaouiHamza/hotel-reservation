package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/baazaouihamza/hotel-reservation/db/fixtures"
	"github.com/baazaouihamza/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setUp(t)
	defer db.teardown(t)

	var (
		user          = fixtures.AddUser(db.Store, "hamza", "baazaoui", false)
		hotel         = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room          = fixtures.AddRoom(db.Store, "small", true, 4.5, hotel.ID)
		from          = time.Now()
		till          = from.AddDate(0, 0, 5)
		booking       = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app           = fiber.New()
		route         = app.Group("/", JWTAuthentication(db.User))
		bookingHadler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHadler.HandlerGetBooking)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.RoomID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

}

func TestAdminGetBookings(t *testing.T) {
	db := setUp(t)
	defer db.teardown(t)

	var (
		adminUser     = fixtures.AddUser(db.Store, "admin", "admin", true)
		user          = fixtures.AddUser(db.Store, "hamza", "baazaoui", false)
		hotel         = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room          = fixtures.AddRoom(db.Store, "small", true, 4.5, hotel.ID)
		from          = time.Now()
		till          = from.AddDate(0, 0, 5)
		booking       = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app           = fiber.New()
		admin         = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHadler = NewBookingHandler(db.Store)
	)

	admin.Get("/", bookingHadler.HandlerGetBookings)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if len(bookings) != 1 {
		t.Fatalf("expected one booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got response got %d", resp.StatusCode)
	}
}
