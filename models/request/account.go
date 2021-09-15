package request

type Login struct {
	Account  string `form:"account" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
