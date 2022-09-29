package main

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
	Паттерн позволяет передавать запрос последовательно по цепочке и в случае необходимость прерывать дальнейшее выполнение
*/

//пример - система регистрации/аутентификации/проверки роли

// интерфейс обработчика
type Handler interface {
	handle(u User)
	setNext(h Handler)
}

//обработчики регистрации, логина, роли

type UserRegHandler struct {
	users map[string]User
	next  Handler
}

func InitUserRegHandler(u ...User) *UserRegHandler {
	m := make(map[string]User)
	for _, us := range u {
		m[us.login] = us
	}
	return &UserRegHandler{users: m}
}

func (urh *UserRegHandler) handle(u User) {
	if _, ok := urh.users[u.login]; !ok {
		fmt.Println("user not found. creating new user")
		urh.users[u.login] = u
		urh.next.handle(u)
		return
	}
	fmt.Println("user already exists!")
	urh.next.handle(u)
}

func (urh *UserRegHandler) setNext(h Handler) {
	urh.next = h
}

type UserLogHandler struct {
	users map[string]User
	next  Handler
}

func InitUserLogHandler(m map[string]User) *UserLogHandler {
	return &UserLogHandler{users: m}
}

func (h *UserLogHandler) handle(u User) {
	if _, ok := h.users[u.login]; !ok {
		fmt.Printf(`User "%s" login failed!`, u.login)
		fmt.Println()
		return
	}
	if u == h.users[u.login] {
		fmt.Printf("User %s is logged in\n", u.login)
		h.next.handle(u)
	}
	return
}

func (h *UserLogHandler) setNext(h1 Handler) {
	h.next = h1
}

type UserRoleHandler struct {
	next Handler
}

func InitUserRoleHandler() *UserRoleHandler {
	return &UserRoleHandler{}
}

func (h *UserRoleHandler) handle(u User) {
	if u.login == "admin" {
		fmt.Println("Hello admin")
	} else if u.login != "" {
		fmt.Println("Hello user ", u.login)
	}
	return
}

func (h *UserRoleHandler) setNext(h1 Handler) {
	h.next = h1
}

type User struct {
	login    string
	password string
}

func main() {
	u1 := &User{
		login:    "admin",
		password: "admin",
	}
	u2 := &User{
		login:    "user1",
		password: "root",
	}

	//клиентские запросы
	regHandler := InitUserRegHandler(*u1, *u2)
	logHandler := InitUserLogHandler(regHandler.users)
	regHandler.setNext(logHandler)
	roleHandler := InitUserRoleHandler()
	logHandler.setNext(roleHandler)

	regHandler.handle(User{login: "admin", password: "admin"})

	regHandler.handle(User{login: "user1", password: "root"})

	regHandler.handle(User{login: "user2", password: "root"})
	//при необходимость можно отправлять запрос не первому элементу цепи
	logHandler.handle(User{login: "nobody", password: "1"})

}
