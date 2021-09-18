package dbmodels

var db IDB

func RegisterDB(idb IDB) {
	db = idb
}

func GetDB() IDB {
	return db
}

type IDB interface {
	IUsers
}

type IUsers interface {
	InitCUsers() error
	EmailExist(email string) (bool, error)
	LoginNameExist(loginName string) (bool, error)
	AddUser(user User) IDBError
	GetUser(account, password string) (User, IDBError)
	UpdatePassword(userId, oldPwd, newPwd string) IDBError
}
