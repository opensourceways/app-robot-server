package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/global"
)

func (c *client) AddPlugin(p dbmodels.Plugin) dbmodels.IDBError {
	doc, err := structToMap(p)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) dbmodels.IDBError {
		_, err := c.newDocIfNotExist(ctx, c.pluginsCollection, bson.M{"name": p.Name}, doc)
		return err
	}

	return withContext1(f)
}

func (c *client) AuditedPlugin(pID, user string) dbmodels.IDBError {
	return nil
}

func (c *client) AddPluginVersion(pName string, version dbmodels.PluginVersion) dbmodels.IDBError {
	docFilter := genPluginFilter(pName, "")
	elemFilter := bson.M{dbmodels.FieldPVNumber: version.VersionNumber}
	arrayFilterByElemMatch(dbmodels.FieldPVersions, false, elemFilter, docFilter)
	v, idbError := structToMap(version)
	if idbError != nil {
		return idbError
	}

	f := func(ctx context.Context) dbmodels.IDBError {
		return c.pushArrayElem(ctx, c.pluginsCollection, dbmodels.FieldPVersions, docFilter, v)
	}
	return withContext1(f)
}

func (c *client) GetUserPlugins(userName string) ([]dbmodels.Plugin, dbmodels.IDBError) {
	projection := bson.M{
		"versions":            0,
		dbmodels.FieldMongoID: 0,
		"audit_by":            0,
	}
	docFilter := bson.M{
		dbmodels.FieldPluginAuthor: userName,
	}
	var result []dbmodels.Plugin

	f := func(ctx context.Context) dbmodels.IDBError {
		if err := c.getDocs(ctx, c.pluginsCollection, docFilter, projection, &result); err != nil {
			if isErrNoDocuments(err) {
				return errNoDBRecord
			}
			return newSystemError(err)
		}
		return nil
	}
	err := withContext1(f)
	return result, err
}

func (c *client) GetPluginDetail(pName, uName string) (dbmodels.Plugin, dbmodels.IDBError) {
	projection := bson.M{
		dbmodels.FieldMongoID: 0,
	}

	docFilter := genPluginFilter(pName, uName)

	var result dbmodels.Plugin
	f := func(ctx context.Context) dbmodels.IDBError {
		return c.getDoc(ctx, c.pluginsCollection, docFilter, projection, &result)
	}
	idbError := withContext1(f)
	return result, idbError
}

func (c *client) UpdatePluginLastVersion(pName, lv string, publish bool) dbmodels.IDBError {
	docFilter := genPluginFilter(pName, "")
	updateFilter := bson.M{"last_version": lv}
	if publish {
		updateFilter["status"] = global.PluginStatusPublished
	}
	f := func(ctx context.Context) dbmodels.IDBError {
		return c.updateDoc(ctx, c.pluginsCollection, docFilter, updateFilter)
	}
	return withContext1(f)
}

func genPluginFilter(pName, userName string) bson.M {
	filter := bson.M{
		dbmodels.FieldPluginName: pName,
	}
	if userName != "" {
		filter[dbmodels.FieldPluginAuthor] = userName
	}
	return filter
}
