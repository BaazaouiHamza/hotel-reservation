package db

import (
	"context"

	"github.com/baazaouihamza/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	hotelColl = "hotels"
)

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

var _ HotelStore = (*MongoHotelStore)(nil)

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil

}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter Map) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil

}

func (s *MongoHotelStore) Update(ctx context.Context, filter Map, update Map) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
