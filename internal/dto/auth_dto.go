package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // NIP
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string       `json:"token"`
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

type UserResponse struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Jabatan  string `json:"jabatan"`
	Role     string `json:"role"`
}
