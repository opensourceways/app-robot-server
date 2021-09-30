package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/opensourceways/app-robot-server/config"
	"github.com/opensourceways/app-robot-server/dbmodels"
)

var _ dbmodels.IDB = (*client)(nil)

type client struct {
	*mongo.Client
	db *mongo.Database

	usersCollection    string
	pluginsCollection  string
	instanceCollection string
}

func (c *client) Close() error {
	return withContext(c.Disconnect)
}

func (c *client) collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (c *client) getUsersCollection() *mongo.Collection {
	return c.collection(c.usersCollection)
}

func (c *client) doTransaction(f func(mongo.SessionContext) error) error {

	callback := func(sc mongo.SessionContext) (interface{}, error) {
		return nil, f(sc)
	}

	s, err := c.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongodb session: %s", err.Error())
	}

	ctx := context.Background()
	defer s.EndSession(ctx)

	_, err = s.WithTransaction(ctx, callback)
	return err
}

func Initialize(cfg *config.MongoDBConfig) (dbmodels.IDB, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(cfg.ConnURI))
	if err != nil {
		return nil, err
	}

	if err = withContext(c.Connect); err != nil {
		return nil, err
	}

	// verify if database connection is created successfully
	err = withContext(func(ctx context.Context) error {
		return c.Ping(ctx, nil)
	})
	if err != nil {
		return nil, err
	}

	cli := &client{
		Client:             c,
		db:                 c.Database(cfg.DBName),
		usersCollection:    cfg.UsersCollection,
		pluginsCollection:  cfg.PluginsCollection,
		instanceCollection: cfg.InstanceCollection,
	}
	return cli, nil
}

func toUID(oid interface{}) (string, error) {
	v, ok := oid.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("retrieve id failed")
	}
	return v.Hex(), nil
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return f(ctx)
}
