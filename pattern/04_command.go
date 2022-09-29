package main

import "fmt"

/*
	Реализовать паттерн «команда».
	Паттерн позволяет инкапсулировать методы
*/

// команда
type Command interface {
	execute()
}

//отправитель

type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}

//конкретные команды

type OnCommand struct {
	device Device
}

type OffCommand struct {
	device Device
}

func (c *OnCommand) execute() {
	c.device.on()
}

func (c *OffCommand) execute() {
	c.device.off()
}

// получатель
type Device interface {
	on()
	off()
}

//конкретный получатель

type Radio struct {
	stateRun bool
}

func (r *Radio) on() {
	r.stateRun = true
	fmt.Println("Turning radio on")
}

func (r *Radio) off() {
	r.stateRun = false
	fmt.Println("Turning radio off")
}

type Fan struct {
	stateRun bool
}

func (f *Fan) on() {
	f.stateRun = true
	fmt.Println("Turning fan on")
}

func (f *Fan) off() {
	f.stateRun = false
	fmt.Println("Turning fan off")
}

func main() {
	fan := &Fan{}
	radio := &Radio{}

	onRadio := OnCommand{device: radio}
	offRadio := OffCommand{device: radio}
	onFan := OnCommand{device: fan}
	offFan := OffCommand{device: fan}

	buttonsRadio := []Button{{command: &onRadio}, {command: &offRadio}}
	buttonsFan := []Button{{command: &onFan}, {command: &offFan}}

	for _, btn := range buttonsRadio {
		btn.press()
	}

	for _, btn := range buttonsFan {
		btn.press()
	}
}
