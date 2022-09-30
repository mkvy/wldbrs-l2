package main

import "fmt"

/*
	Реализовать паттерн «стратегия».
	Паттерн позволяет реализовать разные вариации одного алгоритма внутри объекта
*/

// интерфейс стратегии
type CompareI interface {
	compareAtoB(name1 Name, name2 Name) bool
}

// стратегия сравнения по имени
type CompareFirstName struct{}

func (c *CompareFirstName) compareAtoB(name1 Name, name2 Name) bool {
	return name1.firstName > name2.firstName
}

// стратегия сравнения по фамилии
type CompareLastName struct{}

func (c *CompareLastName) compareAtoB(name1 Name, name2 Name) bool {
	return name1.lastName > name2.lastName
}

// контекст, задаем стратегию сравнения
type Name struct {
	firstName string
	lastName  string
	cAlg      CompareI
}

func InitName(fName string, lName string, alg CompareI) *Name {
	return &Name{firstName: fName, lastName: lName, cAlg: alg}
}

// можем менять стратегию
func (n *Name) SetCompareAlg(alg CompareI) {
	n.cAlg = alg
}

// функция, через которую вызываем стратегию
func (n *Name) Compare(name Name) bool {
	return n.cAlg.compareAtoB(*n, name)
}

func main() {
	cfn := &CompareFirstName{}
	person := InitName("Boris", "Atest", cfn)
	comparePerson := &Name{"Andre", "Checkfamilia", nil}
	isMore := person.Compare(*comparePerson)
	fmt.Printf("Сравниваем имя в лексиграфическом порядке: %s > %s ? %t\n", person.firstName, comparePerson.firstName, isMore)

	person.SetCompareAlg(&CompareLastName{})

	isMore = person.Compare(*comparePerson)
	fmt.Printf("Сравниваем фамилию в лексиграфическом порядке: %s > %s ? %t\n", person.lastName, comparePerson.lastName, isMore)

}
