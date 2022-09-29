package main

import (
	"errors"
	"fmt"
)

/*
	Реализовать паттерн «фасад».
	Фасад позволяет скрыть внутреннюю логику взаимодействия между интерфейсами,
	предоставляя лишь общий интерфейс для объединения под собой нескольких подсистем.
	Тем самым можно предоставить простой или урезанный интерфейс к сложной подсистеме.
	Из плюсов: имеем логическую изоляцию от внутренностей сложных подсистем
	Из минусов: сильная привязанность ко многим компонентам программы
*/

// Определяем сам фасад, включает в себя ссылки на используемые подсистемы
type OrderFacade struct {
	account  *Account
	delivery *Delivery
	items    map[string]Item
}

// Инициализируем фасад
func newOrderFacade(login string, pass string, address string) *OrderFacade {
	fmt.Println("Creating facade")
	items := make(map[string]Item)
	items["pizza"] = &Pizza{price: 50}
	items["sushi"] = &Sushi{price: 100}
	of := &OrderFacade{
		account:  newAccount(login, pass, address),
		delivery: &Delivery{},
		items:    items,
	}
	return of
}

// Определеяем методы фасада
func (of *OrderFacade) depositAccount(amount int) error {
	fmt.Println("depositing money to account")
	return of.account.addMoney(amount)
}

// Используем методы различных интерфейсов, предоставляя клиенту упрощенный интерфейс, скрываем внутреннюю логику
func (of *OrderFacade) placeOrder(item string, count int) error {
	var val Item
	var ok bool
	if val, ok = of.items[item]; !ok {
		return errors.New("item not found")
	}
	err := of.account.makePayment(val.GetPrice() * count)
	if err != nil {
		return err
	}
	val.Cook()
	of.delivery.Deliver(of.account.GetAddress())
	return nil
}

// подсистемы. инициализируем; определяем внутренние методы подсистем, которые используются в фасаде
type Account struct {
	login   string
	pass    string
	address string
	balance int
}

func newAccount(login string, pass string, address string) *Account {
	fmt.Println("Creating user account")
	acc := &Account{
		login:   login,
		pass:    pass,
		address: address,
		balance: 0,
	}
	return acc
}

func (a *Account) addMoney(amount int) error {
	if amount < 0 {
		return errors.New("Cannot add negative amount of money")
	}
	fmt.Println("Adding money to balance: ", amount)
	a.balance += amount
	return nil
}

func (a *Account) makePayment(charge int) error {
	if a.balance-charge < 0 {
		return errors.New("Not enough money to buy")
	}
	a.balance -= charge
	fmt.Printf("Made payment: %d, current balance: %d\n", charge, a.balance)
	return nil
}

func (a *Account) GetAddress() string {
	return a.address
}

// подсистемы
type Delivery struct{}

func (d *Delivery) Deliver(dest string) {
	fmt.Println("delivering to address = ", dest)
}

type Item interface {
	Cook()
	GetPrice() int
}

type Pizza struct {
	price int
}

type Sushi struct {
	price int
}

func (p *Pizza) Cook() {
	fmt.Println("Cooking pizza")
}
func (s *Sushi) Cook() {
	fmt.Println("Cooking sushi")
}
func (p *Pizza) GetPrice() int {
	return p.price
}
func (s *Sushi) GetPrice() int {
	return s.price
}

func main() {
	//инициализируем подсистему
	of := newOrderFacade("login", "123", "Order Street 18")
	//используем методы фасада, не вызывая методов других структур
	of.depositAccount(500)
	err := of.placeOrder("pizza", 2)
	if err != nil {
		fmt.Println(err)
	}
	err = of.placeOrder("sushi", 1)
	if err != nil {
		fmt.Println(err)
	}
	err = of.placeOrder("sushi", 10)
	if err != nil {
		fmt.Println(err)
	}
}
