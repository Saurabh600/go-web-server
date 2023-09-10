package data

type User struct {
	Id           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"Password_hash"`
	Age          int    `json:"age"`
	Gender       string `json:"gender"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserFormData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type EmailAndPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
