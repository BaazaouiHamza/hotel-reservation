package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/baazaouihamza/hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setUp(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "admin", "admin", true)
	hotel := fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 4.5, hotel.ID)

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
	fmt.Println(booking)
}
