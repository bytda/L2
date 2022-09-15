/*
	Паттерн «посетитель»

	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Паттерн Посетитель

Паттерн Посетитель используется для добавления нового функционала к уже существующему
объекту не меняя его структуру или же когда постоянно меняется функционал объекта.
Например, когда у нас есть объект выполняющих какую-то функцию и нам нужно разово добавить
какой-то дополнительный функционал, мы можем использовать паттерн Посетитель

Плюсы:

Упрощает процесс добавления нового функционала, не нагромождая каждый раз объект множеством
новых методов, облегчает читабельность кода, так как основной объект не загроможден
дополнительным кодом.

Минусы:

Добавление новых объектов, которым тоже необходим visitor затрудняет использование паттерна,
необходимо будет прописать много кода, если иерархия объектов поменяется, поэтому паттерн лучше
применять когда часто меняется функционал объектов, а не сами объекты

Таким образом, паттерн "Посетитель" обычно применяют для следующих задач:
	1) Выполнить операцию над всеми элементами сложной объектной структуры
	2) Добавить новое поведение для некоторых классов в иерархии (оставив поведение пустым для остальных)
*/

type Visitor struct { //содержит методы для определенных структур
}

func (v *Visitor) VisitA(a *A) { //метод для структуры A
	/*...*/
}

func (v *Visitor) VisitB(b *B) { //метод для структуры B
	/*...*/
}

type Element interface { //содержит метод вызова Visit в обходимых структурах
	Accept(v *Visitor)
}

type A struct { //структура содержащая метод вызова визитера
}

type B struct { //структура содержащая метод вызова визитера
}

func (a *A) Accept(v *Visitor) { //вызывает метод визитера для этой структуры
	v.VisitA(a)
}

func (b *B) Accept(v *Visitor) { //вызывает метод визитера для этой структуры

	v.VisitB(b)
}
