package main

import "fmt"

/*
Реализовать паттерн «состояние».
Реализуется конечный автомат, объект меняет поведение в зависимости от состояния
*/

// контекст, включает в себя состояния
type WaterCooler struct {
	hasWater State
	noWater  State

	currentState State

	waterAmount int
	cap         int
}

func newWaterCooler(waterAmount, capacity int) *WaterCooler {
	w := &WaterCooler{
		waterAmount: waterAmount,
		cap:         capacity,
	}
	hasWaterState := &HasWaterState{
		wc: w,
	}
	noWaterState := &NoWaterState{
		wc: w,
	}
	w.setState(hasWaterState)
	w.hasWater = hasWaterState
	w.noWater = noWaterState
	return w
}

func (w *WaterCooler) requestWater(amount int) {
	w.currentState.requestWater(amount)
}

func (w *WaterCooler) addWater(amount int) {
	w.currentState.addWater(amount)
}

func (w *WaterCooler) IncrWater(amount int) {
	fmt.Println("adding water")
	if w.waterAmount+amount > w.cap {
		w.waterAmount = w.cap
		return
	}
	w.waterAmount += amount
}

func (w *WaterCooler) setState(st State) {
	w.currentState = st
}

// интерфейс для действий у каждого состояния
type State interface {
	addWater(int)
	requestWater(int)
}

// состояние когда нет воды
type NoWaterState struct {
	wc *WaterCooler
}

func (s *NoWaterState) addWater(amount int) {
	s.wc.IncrWater(amount)
	//переводим в другое состояние
	s.wc.setState(s.wc.hasWater)
}

func (s *NoWaterState) requestWater(amount int) {
	fmt.Println("Error, no water for this amount: ", amount)
}

// состояние когда есть вода
type HasWaterState struct {
	wc *WaterCooler
}

func (s *HasWaterState) addWater(amount int) {
	s.wc.IncrWater(amount)
}

func (s *HasWaterState) requestWater(amount int) {
	if amount > s.wc.waterAmount {
		fmt.Println("Not enough water")
		return
	}
	if amount == s.wc.waterAmount {
		fmt.Println("Pouring water amount ", amount)
		s.wc.waterAmount = 0
		//переводим в другое состояние если вода кончилась
		s.wc.setState(s.wc.noWater)
		return
	}
	fmt.Println("Pouring water amount ", amount)
	s.wc.waterAmount -= amount
}

func main() {
	wc := newWaterCooler(15, 25)
	wc.addWater(15)
	fmt.Println("-Water amount after adding: ", wc.waterAmount)
	fmt.Println("-Requesting 26 water while having 25: ")
	wc.requestWater(26)
	fmt.Println("-Requesting 25 water while having 25: ")
	wc.requestWater(25)
	fmt.Println("-Adding 10 water: ")
	wc.addWater(10)
	fmt.Println("-Requesting 8 water while having 10: ")
	wc.requestWater(8)
}
