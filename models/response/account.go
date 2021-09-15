package response

type LoginResult struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
