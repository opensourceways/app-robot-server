package dbmodels

const (
	FieldUsers     = "users"
	FieldEmail     = "email"
	FieldLoginName = "login_name"
	FieldPassword  = "password"
	FieldUserID = "id"
)

type CUsers struct {
	Users []User `json:"users"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	LoginName string `json:"login_name" bson:"login_name"`
	Password  string `json:"password"`
}
