package dbmodels

import "github.com/opensourceways/app-robot-server/global"

const FieldMongoID = "_id"

var db IDB

func RegisterDB(idb IDB) {
	db = idb
}

func GetDB() IDB {
	return db
}

type IDB interface {
	IUsers
	IPlugins
	IInstance
}

type IUsers interface {
	InitCUsers() error
	EmailExist(email string) (bool, error)
	LoginNameExist(loginName string) (bool, error)
	AddUser(user User) IDBError
	GetUser(account, password string) (User, IDBError)
	UpdatePassword(userId, oldPwd, newPwd string) IDBError
}

type IPlugins interface {
	AddPlugin(p Plugin) IDBError
	AuditedPlugin(pName, user string) IDBError
	AddPluginVersion(pName string, version PluginVersion) IDBError
	GetUserPlugins(userName string) ([]Plugin, IDBError)
	GetPluginDetail(pName, uName string) (Plugin, IDBError)
	UpdatePluginLastVersion(pName, lv string, publish bool) IDBError
}

type IInstance interface {
	AddInstance(i Instance) IDBError
	GetInstance(insID string) (Instance, IDBError)
	UpdateInstanceStatus(insID string, status global.InstanceStatus) IDBError
}
