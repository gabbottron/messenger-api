package datatypes

// This is the user object attached to the JWT
type User struct {
	UserID       int
	SessionToken string
}

// Basic JSON http response object
type HttpResponseJSON struct {
	Status  string `json:"status" binding:"required"`
	Message string `json:"message" binding:"required"`
}
