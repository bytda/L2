/*
	Паттерн «строитель».

	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Строитель:

	Паттерн строитель используется для упрощения создания сложного объекта. Построение
сложного объекта основывается на использовании простых объектов и разделения процесса
на последовательные этапы. Паттерн использует интерфейс builder,
с описанием этапов строительства, но без конкретной реализации, а также
классы concreteBuilder, которые уже реализуют конкретный способ создания сложного объекта.
Паттерн строитель применяют, когда необходимо обеспечить разные варианты схожих объектов,
процесс создания объектов не зависит от того, из каких частей объект состоит и как они взаимосвязаны,
при этом создаваемые объекты сложные, состоящие из множества компонентов и их создание нельзя
осуществить в одно действие. Например, система для пользователей с разным доступом(покупатель, поставщик, администратор),
в общем объекты однотипны, но у них есть разные права доступа, покупатель не должен менять
цены в магазине и тд

Плюсы:

Облегчает код и делает его более читабельным, возможность пользоваться готовыми методами
создания объектов, вместо того чтобы тратить время на изучение в устройстве структур.
Так же мы видим детали создания, паттерн предоставляет интерфейс с поэтапным описанием как создается
объект, если необходимо мы можем изменить какой-то этап под нужную нам задачу, определив
новый concreteBuilder

Минусы:

Жесткая привязка к concreteBuilder с разновидностью объекта, под каждую вариацию объекта или
же изменений в описании объекта необходимо писать новый concreteBuilder

Таким образом, паттерн "строитель" обычно применяют для следующих задач:
	1) Пошаговая реализация сложной структуры,
	2) Избавление от конструктора с большим кол-вом параметров,
	3) Компоновка методов инициализации структуры в одном месте,
	4) Создание разных представлений одной структуры
*/

type home struct { //Home собираемый объект
	/*...*/
}

type homeBuilder struct { //строитель структуры Home

	result House
}

func (h *homeBuilder) Reset() { //инициализирует новый экземпляр структуры Home

	h.result = House{}
}

func (h *homeBuilder) SetWalls(num int) { //добавляет указанное кол-во стен

	/*...*/
}

func (h *homeBuilder) SetRoof() { //добавляет крышу

	/*...*/
}

func (h *homeBuilder) SetDoors(num int) { //добавляет указанное кол-во дверей

	/*...*/
}

func (h *homeBuilder) SetWindows(num int) { //добавляет указанное кол-во окон

	/*...*/
}

func (h *homeBuilder) GetResult() Home { //возвращает собранную структуру Home

	return h.result
}