package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/utils"
)

func (c *client) InitCUsers() error {
	var count int64
	dcf := func(ctx context.Context) error {
		docCount, err := c.getUsersCollection().CountDocuments(ctx, bson.D{})
		if err != nil {
			return err
		}
		count = docCount
		return nil
	}
	if err := withContext(dcf); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	cUsers := dbmodels.CUsers{Users: []dbmodels.User{}}
	//doc, _ := structToMap(cUsers)
	f := func(ctx context.Context) error {
		_, err := c.getUsersCollection().InsertOne(ctx, cUsers)
		return err
	}
	return withContext(f)
}

func (c *client) EmailExist(email string) (bool, error) {
	isNotExist := false
	err := withContext(genIsUserNotExistFunc(bson.M{dbmodels.FieldEmail: email}, &isNotExist, c))
	return !isNotExist, err
}

func (c *client) LoginNameExist(loginName string) (bool, error) {
	isNotExist := false
	err := withContext(genIsUserNotExistFunc(bson.M{dbmodels.FieldLoginName: loginName}, &isNotExist, c))
	return !isNotExist, err
}

func (c *client) AddUser(user dbmodels.User) dbmodels.IDBError {
	value, idbError := structToMap(user)
	if idbError != nil {
		return idbError
	}
	addFilter := bson.M{
		"$nor": bson.A{
			bson.M{"users.login_name": user.LoginName},
			bson.M{"users.email": user.Email},
		},
	}
	return withContext1(func(ctx context.Context) dbmodels.IDBError {
		return c.pushArrayElem(ctx, c.usersCollection, dbmodels.FieldUsers, addFilter, value)
	})
}

func (c *client) GetUser(account, password string) (dbmodels.User, dbmodels.IDBError) {
	var elemFilter bson.M
	if utils.IsEmail(account) {
		elemFilter = bson.M{dbmodels.FieldEmail: account, dbmodels.FieldPassword: password}
	} else {
		elemFilter = bson.M{dbmodels.FieldLoginName: account, dbmodels.FieldPassword: password}
	}
	docFilter := bson.M{}
	arrayFilterByElemMatch(dbmodels.FieldUsers, true, elemFilter, docFilter)
	var result []dbmodels.CUsers
	f := func(ctx context.Context) error {
		return c.getArrayElem(ctx, c.usersCollection, dbmodels.FieldUsers, docFilter, elemFilter, nil, &result)
	}
	var user dbmodels.User
	if err := withContext(f); err != nil {
		return user, newSystemError(err)
	}
	if len(result) != 1 || len(result[0].Users) != 1 {
		return user, errNoDBRecord
	}
	user = result[0].Users[0]
	return user, nil
}

func (c *client) UpdatePassword(userId, oldPwd, newPwd string) dbmodels.IDBError {
	updateCmd := bson.M{dbmodels.FieldPassword: newPwd}
	elemFilter := bson.M{dbmodels.FieldUserID: userId, dbmodels.FieldPassword: oldPwd}
	docFilter := bson.M{}
	arrayFilterByElemMatch(dbmodels.FieldUsers, true, elemFilter, docFilter)
	f := func(ctx context.Context) dbmodels.IDBError {
		return c.updateArrayElem(ctx, c.usersCollection, dbmodels.FieldUsers, docFilter, elemFilter, updateCmd)
	}
	return withContext1(f)
}

func genIsUserNotExistFunc(elemFilter bson.M, result *bool, c *client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		v, err := c.isArrayElemNotExists(ctx, c.usersCollection, dbmodels.FieldUsers, bson.M{}, elemFilter)
		if err != nil {
			return err
		}
		*result = v
		return nil
	}
}
