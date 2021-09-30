package mongodb

import (
	"context"
	"github.com/opensourceways/app-robot-server/global"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/app-robot-server/dbmodels"
)

func (c *client) AddInstance(ins dbmodels.Instance) dbmodels.IDBError {
	doc, err := structToMap(ins)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) dbmodels.IDBError {
		_, err := c.newDocIfNotExist(ctx, c.instanceCollection, bson.M{"p_name": ins.PName}, doc)
		return err
	}

	return withContext1(f)
}

func (c *client) GetInstance(insID string) (dbmodels.Instance, dbmodels.IDBError) {
	docFilter := bson.M{"id": insID}
	var result dbmodels.Instance
	dbErr := withContext1(func(ctx context.Context) dbmodels.IDBError {
		return c.getDoc(ctx, c.instanceCollection, docFilter, nil, &result)
	})
	return result, dbErr
}

func (c *client) UpdateInstanceStatus(insID string, status global.InstanceStatus) dbmodels.IDBError {
	docFilter := bson.M{"id": insID}
	update := bson.M{"status":status}
	return withContext1(func(ctx context.Context) dbmodels.IDBError {
		return c.updateDoc(ctx,c.instanceCollection,docFilter,update)
	})
}
