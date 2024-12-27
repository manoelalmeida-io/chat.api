package model

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	GoogleSub string `json:"googleSub"`
}

type UserContact struct {
	Id        *int64  `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
	UserId    *int64  `json:"userId"`
}
