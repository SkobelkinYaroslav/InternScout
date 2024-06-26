package user

type User struct {
	id         int      `json:"id"`
	categories []string `json:"categories"`
}
