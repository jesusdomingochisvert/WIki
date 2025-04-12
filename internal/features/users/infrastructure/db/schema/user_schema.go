package schema

type UserSchema struct {
	ID       string `bson:"_id:omitempty"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}
