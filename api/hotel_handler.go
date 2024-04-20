package api

import (
	"fmt"

	"github.com/baazaouihamza/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}

	fmt.Printf("query params %+v \n", qparams)

	hotels, err := h.hotelStore.GetHotels(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
