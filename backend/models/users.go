package models

type User struct {
	Username  string `bson:"username" json:"username"`
	PublicKey string `bson:"public_key" json:"public_key"`
}
