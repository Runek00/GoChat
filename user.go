package main

type User struct {
	login    string
	password string
}

var users = make(map[string]User, 0)

func CheckUser(login string, password string) bool {
	usr, ok := GetUser(login)
	return ok && usr.password == password
}

func GetUser(login string) (User, bool) {
	usr, exists := users[login]
	return usr, exists
}

func AddUser(login string, password string) {
	u := User{login, password}
	users[u.login] = u
}
