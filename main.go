package main

import (
	"context"
	"flag"
	"log"

	"github.com/baazaouihamza/hotel-reservation/api"
	"github.com/baazaouihamza/hotel-reservation/api/middleware"
	"github.com/baazaouihamza/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new fiber instance with custom config
var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	// create connection options
	clientOptions := options.Client().ApplyURI(db.DBURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting!", err)
	}

	// handlers initialization
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Room:  roomStore,
			Hotel: hotelStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		roomHandler  = api.NewRoomHandler(store)
		app          = fiber.New(config)
		auth         = app.Group("/api")
		apiv1        = app.Group("api/v1", middleware.JWTAuthentication(userStore))
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	app.Listen(*listenAddr)
}
