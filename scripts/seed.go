package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baazaouihamza/hotel-reservation/api"
	"github.com/baazaouihamza/hotel-reservation/db"
	"github.com/baazaouihamza/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := bookingStore.InsertBooking(context.Background(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking:", insertedBooking.ID)
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedhotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedhotel

}

func main() {
	hamza := seedUser(false, "hamza", "baazaoui", "hamza@baazaoui.com", "1234567")
	seedUser(true, "admin", "admin", "admin@admin.com", "admin")
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "The Nederlands", 4)
	hotel := seedHotel("Dont die in your sleep", "London", 1)
	seedRoom("small", true, 98.99, hotel.ID)
	seedRoom("medium", true, 198.99, hotel.ID)
	room := seedRoom("large", false, 298.99, hotel.ID)
	seedBooking(hamza.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))

}

func init() {
	var err error
	// create connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting!", err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
