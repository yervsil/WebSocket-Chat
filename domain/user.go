package domain

type User struct {
	Id 			  int  	 `db:"id"`
	Email         string `db:"email"`
	Name	      string `db:"username"`
	Password_hash string `db:"password_hash"`
}
