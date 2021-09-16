package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/huaweicloud/golangsdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/opensourceways/app-robot-server/dbmodels"
)

func withContext1(f func(context.Context) dbmodels.IDBError) dbmodels.IDBError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return f(ctx)
}

func structToMap(info interface{}) (bson.M, dbmodels.IDBError) {
	body, err := golangsdk.BuildRequestBody(info, "")
	if err != nil {
		return nil, newDBError(dbmodels.ErrMarshalDataFaield, err)
	}
	return body, nil
}

func arrayFilterByElemMatch(array string, exists bool, cond, filter bson.M) {
	match := bson.M{"$elemMatch": cond}
	if exists {
		filter[array] = match
	} else {
		filter[array] = bson.M{"$not": match}
	}
}

func (c *client) pushArrayElem(ctx context.Context, collection, array string, filterOfDoc, value bson.M) dbmodels.IDBError {
	update := bson.M{"$push": bson.M{array: value}}

	col := c.collection(collection)
	r, err := col.UpdateOne(ctx, filterOfDoc, update)
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}

func (c *client) pushArrayElems(ctx context.Context, collection, array string, filterOfDoc bson.M, value bson.A) dbmodels.IDBError {
	update := bson.M{"$push": bson.M{array: bson.M{"$each": value}}}

	col := c.collection(collection)
	r, err := col.UpdateOne(ctx, filterOfDoc, update)
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}

func (c *client) replaceDoc(ctx context.Context, collection string, filterOfDoc, docInfo bson.M) (string, dbmodels.IDBError) {
	upsert := true

	col := c.collection(collection)
	r, err := col.ReplaceOne(
		ctx, filterOfDoc, docInfo,
		&options.ReplaceOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", newSystemError(err)
	}

	if r.UpsertedID == nil {
		return "", nil
	}

	v, _ := toUID(r.UpsertedID)
	return v, nil
}

func (c *client) deleteDoc(ctx context.Context, collection string, filterOfDoc bson.M) dbmodels.IDBError {
	col := c.collection(collection)
	if _, err := col.DeleteOne(ctx, filterOfDoc); err != nil {
		return newSystemError(err)
	}

	return nil
}

func (c *client) pullArrayElem(ctx context.Context, collection, array string, filterOfDoc, filterOfArray bson.M) dbmodels.IDBError {
	update := bson.M{"$pull": bson.M{array: filterOfArray}}

	col := c.collection(collection)
	r, err := col.UpdateOne(ctx, filterOfDoc, update)
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}

func (c *client) pushNestedArrayElem(ctx context.Context, collection, array string, filterOfDoc, filterOfArray, updateCmd bson.M) dbmodels.IDBError {
	return c.updateArrayElemHelper(ctx, collection, array, filterOfDoc, filterOfArray, updateCmd, "$addToSet")
}

// r, _ := col.UpdateOne; r.ModifiedCount == 0 will happen in two case:
// 1. no matched array item;
// 2 update repeatedly with same update cmd.
func (c *client) updateArrayElem(ctx context.Context, collection, array string, filterOfDoc, filterOfArray, updateCmd bson.M) dbmodels.IDBError {
	return c.updateArrayElemHelper(ctx, collection, array, filterOfDoc, filterOfArray, updateCmd, "$set")
}

func (c *client) updateArrayElemHelper(ctx context.Context, collection, array string, filterOfDoc, filterOfArray, updateCmd bson.M, op string) dbmodels.IDBError {
	cmd := bson.M{}
	for k, v := range updateCmd {
		cmd[fmt.Sprintf("%s.$[i].%s", array, k)] = v
	}

	arrayFilter := bson.M{}
	for k, v := range filterOfArray {
		arrayFilter["i."+k] = v
	}

	col := c.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc,
		bson.M{op: cmd},
		&options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: bson.A{
					arrayFilter,
				},
			},
		},
	)
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}

func (c *client) pullAndReturnArrayElem(ctx context.Context, collection, array string, filterOfDoc, filterOfArray bson.M, result interface{}) dbmodels.IDBError {
	col := c.collection(collection)
	sr := col.FindOneAndUpdate(
		ctx, filterOfDoc,
		bson.M{"$pull": bson.M{array: filterOfArray}},
		&options.FindOneAndUpdateOptions{
			Projection: bson.M{array: bson.M{"$elemMatch": filterOfArray}},
		})

	if err := sr.Decode(result); err != nil {
		if isErrNoDocuments(err) {
			return errNoDBRecord
		}
		return newSystemError(err)
	}
	return nil
}

func (c *client) moveArrayElem(ctx context.Context, collection, from, to string, filterOfDoc, filterOfArray, value bson.M) dbmodels.IDBError {
	col := c.collection(collection)

	r, err := col.UpdateOne(
		ctx, filterOfDoc,
		bson.M{
			"$pull": bson.M{from: filterOfArray},
			"$push": bson.M{to: value},
		},
	)
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}

func (c *client) getDoc(ctx context.Context, collection string, filterOfDoc, project bson.M, result interface{}) dbmodels.IDBError {
	col := c.collection(collection)

	var sr *mongo.SingleResult
	if len(project) > 0 {
		sr = col.FindOne(ctx, filterOfDoc, &options.FindOneOptions{
			Projection: project,
		})
	} else {
		sr = col.FindOne(ctx, filterOfDoc)
	}

	if err := sr.Decode(result); err != nil {
		if isErrNoDocuments(err) {
			return errNoDBRecord
		}
		return newSystemError(err)
	}
	return nil
}

func (c *client) newDocIfNotExist(ctx context.Context, collection string, filterOfDoc, docInfo bson.M) (string, dbmodels.IDBError) {
	upsert := true

	col := c.collection(collection)
	r, err := col.UpdateOne(
		ctx, filterOfDoc, bson.M{"$setOnInsert": docInfo},
		&options.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		return "", newSystemError(err)
	}

	if r.UpsertedID == nil {
		return "", newDBError(dbmodels.ErrRecordExists, fmt.Errorf("the doc exists"))
	}

	v, _ := toUID(r.UpsertedID)
	return v, nil
}

func (c *client) updateDoc(ctx context.Context, collection string, filterOfDoc, update bson.M) dbmodels.IDBError {
	col := c.collection(collection)
	r, err := col.UpdateOne(ctx, filterOfDoc, bson.M{"$set": update})
	if err != nil {
		return newSystemError(err)
	}

	if r.MatchedCount == 0 {
		return errNoDBRecord
	}
	return nil
}
