package api

import (
	"net/http"

	"github.com/baazaouihamza/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{
		"canceled": true,
	}); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type: "msg",
		Msg:  "updated",
	})

}

// TODO: This needs to be admin authorized
func (h *BookingHandler) HandlerGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// TODO: this needs to be user authorized
func (h *BookingHandler) HandlerGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	return c.JSON(booking)
}
