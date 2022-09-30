package main

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
	Паттерн определяет интерфейс создания объекта и делегирует операцию создания экземпляра субклассам
*/

// фабрика
type Factory interface {
	createVehicle(model string) VehicleProduct
}

// конкретные фабрики для каждого продукта
type suvFactory struct {
}

func (f *suvFactory) createVehicle(model string) VehicleProduct {
	return &SuvVehicle{model: model}
}

type sportFactory struct {
}

func (f *sportFactory) createVehicle(model string) VehicleProduct {
	return &SportVehicle{model: model}
}

// определяем общий интерфейс для продуктов
type VehicleProduct interface {
	Run()
}

// конкретные продукты
type SuvVehicle struct {
	model string
}

func (v *SuvVehicle) Run() {
	fmt.Printf("Suv vehicle %s runs\n", v.model)
}

type SportVehicle struct {
	model string
}

func (v *SportVehicle) Run() {
	fmt.Printf("Sport vehicle %s runs\n", v.model)
}

// сам фабричный метод
func getVehicleFactory(vehicleType string) Factory {
	switch vehicleType {
	case "suv":
		return &suvFactory{}
	case "sport":
		return &sportFactory{}
	}
	return nil
}

func main() {
	suvF := getVehicleFactory("suv")
	if suvF == nil {
		panic("No such factory")
	}
	suv := suvF.createVehicle("Test")
	suv.Run()
	sportF := getVehicleFactory("sport")
	if sportF == nil {
		panic("No such factory")
	}
	sport := sportF.createVehicle("Test 2")
	sport.Run()

}
