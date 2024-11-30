package movie

import (
	"context"
	"fmt"

	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/config"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Movie struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Movie  string             `json:"movie" bson:"movie"`
	Actors []string           `json:"actors" bson:"actors"`
}

type Store struct {
	mongoClient *mongo.Client
	config      *config.Config
}

func NewStore(mongoClient *mongo.Client, config *config.Config) *Store {
	return &Store{mongoClient: mongoClient, config: config}
}

func (s *Store) CreateMovie(ctx context.Context, movie Movie) error {
	ulog := logger.FromContext(ctx)
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return err
	}
	ulog.Info(fmt.Sprintf("Insert a record with ID: %s", result.InsertedID))
	return nil
}

func (s *Store) CreateMovies(ctx context.Context, movies []Movie) error {
	ulog := logger.FromContext(ctx)
	// convert movies to any slice
	newMovies := make([]any, len(movies))
	for i, movie := range movies {
		newMovies[i] = movie
	}
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	result, err := collection.InsertMany(ctx, newMovies)
	if err != nil {
		return err
	}
	ulog.Info(fmt.Sprintf("Insert records with : %v", result))
	return nil
}

func (s *Store) UpdateMovie(ctx context.Context, movieID string, movie Movie) error {
	ulog := logger.FromContext(ctx)
	// convert id from hex to primitive.ObjectID
	ID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": ID,
	}
	update := bson.M{
		"$set": bson.M{
			"movie":  movie.Movie,
			"actors": movie.Actors,
		},
	}
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	ulog.Info(fmt.Sprintf("new record %v", result))
	return nil
}

func (s *Store) DeleteMovie(ctx context.Context, movieID string) error {
	ulog := logger.FromContext(ctx)
	// convert id from hex to primitive.ObjectID
	ID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": ID,
	}
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	ulog.Info(fmt.Sprintf("deleted result %v", result))
	return nil
}

func (s Store) Find(ctx context.Context, movieName string) (*Movie, error) {
	ulog := logger.FromContext(ctx)
	var result Movie

	filter := bson.D{{Key: "movie", Value: movieName}}
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	return &result, nil
}

func (s *Store) FindAll(ctx context.Context, movieName string) ([]Movie, error) {
	ulog := logger.FromContext(ctx)
	var results []Movie

	filter := bson.D{{Key: "movie", Value: movieName}}
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	return results, nil
}

func (s *Store) ListAll(ctx context.Context) ([]Movie, error) {
	ulog := logger.FromContext(ctx)
	var results []Movie
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	cursor, err := collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	return results, nil
}

func (s *Store) DeleteAll(ctx context.Context) error {
	ulog := logger.FromContext(ctx)
	collection := s.mongoClient.Database(s.config.DBName).
		Collection(s.config.CollectionName)
	deleted, err := collection.DeleteMany(ctx, bson.D{{}}, nil)
	if err != nil {
		ulog.Error(fmt.Sprintf("%v", err))
		return err
	}
	ulog.Info(fmt.Sprintf("records deleted: %v", deleted.DeletedCount))
	return nil
}
