package Models

type User struct {
	Id      uint   `json:"id"`
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"email"`
	Password   string `json:"password"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (b *User) TableName() string {
	return "users"
}
