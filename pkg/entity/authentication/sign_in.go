package entity

type LoginRequestDto struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type LoginDto struct {
	Token   string `json:"token"`
	Expired string `json:"expiredAt"`
	Status  string `json:"status"`
}
