package main

import "fmt"

/*
	Реализовать паттерн «посетитель».
	Паттерн позволяет добавлять логику к объекту, не добавляя реализацию в сам класс объекта
	Добавляем лишь метод, который вызывает сторонний интерфейс
*/

type Visitor interface {
	VisitForETH(w ETHWallet)
	VisitForBTC(w BTCWallet)
}

// элемент
type Wallet interface {
	accept(v Visitor)
}

// конкретные элементы
type ETHWallet struct {
	balance float32
}

// метод accept, который вызывает необходимый метод посещения
func (w *ETHWallet) accept(v Visitor) {
	v.VisitForETH(*w)
}

type BTCWallet struct {
	balance float32
}

func (w *BTCWallet) accept(v Visitor) {
	v.VisitForBTC(*w)
}

// конкретные посетители
type USDConvertVisitor struct {
}

func (cv *USDConvertVisitor) VisitForETH(w ETHWallet) {
	fmt.Printf("ETH balance in USD: %f\n", w.balance*1338)
}

func (cv *USDConvertVisitor) VisitForBTC(w BTCWallet) {
	fmt.Printf("BTC balance in USD: %f\n", w.balance*19476)
}

type RUBConvertVisitor struct {
}

func (cv *RUBConvertVisitor) VisitForETH(w ETHWallet) {
	fmt.Printf("ETH balance in RUB: %f\n", w.balance*76882)
}

func (cv *RUBConvertVisitor) VisitForBTC(w BTCWallet) {
	fmt.Printf("BTC balance in RUB: %f\n", w.balance*1119870)
}

func main() {
	ethWallet := &ETHWallet{
		balance: 0.31,
	}
	btcWallet := &BTCWallet{balance: 0.09}

	wallets := []Wallet{ethWallet, btcWallet}
	fmt.Printf("Баланс изначальный, ETH:%f BTC:%f\n", ethWallet.balance, btcWallet.balance)
	fmt.Println("Конвертация в USD")
	ucv := &USDConvertVisitor{}
	//вызываем метод accept который у каждой структуры вызовет метод для своего типа
	for _, w := range wallets {
		w.accept(ucv)
	}
	fmt.Println("Конвертация в RUB")
	rcv := &RUBConvertVisitor{}
	for _, w := range wallets {
		w.accept(rcv)
	}
}
