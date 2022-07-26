package mongo

import (
	"context"
	"github.com/navisot/go-url-shortener/shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// mongoRepository struct keeps the  client, database, timeout
type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// newMongoClient returns a new mongo client
func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewMongoRepository returns a new repo
func NewMongoRepository(mongoURL, mongoDb string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{timeout: time.Duration(mongoTimeout) * time.Second, database: mongoDb}
	newClient, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = newClient
	return repo, nil
}

// Find returns a redirect from DB
func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find.NotFound")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

// Store stores a redirect to DB
func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
