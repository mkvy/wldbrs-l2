package main

import "fmt"

/*
Реализовать паттерн «строитель».

Может использоваться когда у создаваемого объекта много параметров и не все опциональные параметры нужны при инициализации
Используя строитель можно собирать объект пошагово
Также может использоваться для создания разных представлений какого-либо объекта
*/

// создаваемый объект
type Vehicle struct {
	Doors           int
	EngineDispl     float32
	gearbox         string
	class           string
	hasHeater       bool
	hasAC           bool
	hasPowerWindows bool
}

// интерфейс билдера
type VehicleBuilderI interface {
	Doors(val int) VehicleBuilderI
	EngineDispl(val float32) VehicleBuilderI
	Gearbox(val string) VehicleBuilderI
	Class(val string) VehicleBuilderI
	Heater(val bool) VehicleBuilderI
	AC(val bool) VehicleBuilderI
	PowerWindows(val bool) VehicleBuilderI
	Build() Vehicle
}

// структура билдера
type vehicleBuilder struct {
	doors           int
	engineDispl     float32
	gearbox         string
	class           string
	hasHeater       bool
	hasAC           bool
	hasPowerWindows bool
}

func NewVehicleBuilder() VehicleBuilderI {
	return &vehicleBuilder{}
}

// задаем значения структуры, возвращаем ссылку на самого себя
func (vb *vehicleBuilder) Doors(val int) VehicleBuilderI {
	vb.doors = val
	return vb
}

func (vb *vehicleBuilder) EngineDispl(val float32) VehicleBuilderI {
	vb.engineDispl = val
	return vb
}

func (vb *vehicleBuilder) Gearbox(val string) VehicleBuilderI {
	vb.gearbox = val
	return vb
}

func (vb *vehicleBuilder) Class(val string) VehicleBuilderI {
	vb.class = val
	return vb
}

func (vb *vehicleBuilder) Heater(val bool) VehicleBuilderI {
	vb.hasHeater = val
	return vb
}

func (vb *vehicleBuilder) AC(val bool) VehicleBuilderI {
	vb.hasAC = val
	return vb
}

func (vb *vehicleBuilder) PowerWindows(val bool) VehicleBuilderI {
	vb.hasPowerWindows = val
	return vb
}

// возвращаем объект с параметрами, которые были заданы через билдер
func (vb *vehicleBuilder) Build() Vehicle {
	return Vehicle{
		Doors:           vb.doors,
		EngineDispl:     vb.engineDispl,
		gearbox:         vb.gearbox,
		class:           vb.class,
		hasHeater:       vb.hasHeater,
		hasAC:           vb.hasAC,
		hasPowerWindows: vb.hasPowerWindows,
	}
}

// директор
type Maker struct {
	vb VehicleBuilderI
}

func NewMaker(b VehicleBuilderI) *Maker {
	return &Maker{
		vb: b,
	}
}

// директор, вызывает методы строителя
func (m *Maker) BuildSUV() Vehicle {
	return m.vb.Doors(5).Class("SUV").EngineDispl(2.5).AC(true).Heater(true).PowerWindows(true).Gearbox("auto").Build()
}

func (m *Maker) BuildCheapCar() Vehicle {
	//можем использовать не все параметры, т.к. есть zero-value
	return m.vb.Doors(2).Class("B class").EngineDispl(1.2).Heater(true).Gearbox("manual").Build()
}

func main() {
	//вызываем без директора
	vehicleBuilder := NewVehicleBuilder()
	vehicle := vehicleBuilder.Doors(4).AC(true).Heater(true).Gearbox("manual").EngineDispl(1.3).Build()
	fmt.Println("Car:")
	fmt.Println(vehicle)
	dir := NewMaker(vehicleBuilder)
	dir.BuildSUV()
	fmt.Println("SUV:")
	fmt.Println(dir.BuildSUV())
	fmt.Println("Cheap car:")
	fmt.Println(dir.BuildCheapCar())

}
