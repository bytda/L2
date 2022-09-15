/*
	Паттерн «Стратегия»

	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
	Паттерн "стратегия" используется для определения группы алгоритмов,
сокрытия их реализации (принцип инкапсуляция), обеспечения взаимозаменяемости -
можно выбирать нужный алгоритм в зависимости от потребности Паттерн применяется
в тех случаях, когда системе нужно использовать разные вариации какого-то алгоритма
в своей работе, когда есть множество схожих объектов, которые отличаются деталями
в поведении. Пример использования: система, где есть разные типы пользователей с
разными правами доступа (покупатели, сотрудники, администратор), нам необходимо
реализовать функцию просмотра ассортимента товаров: тогда для покупателей должна
выводиться информация, интересующая именно их, сотрудникам выводится техническая
информация о товаре, администраторам сервиса другая техническая информация.
Тогда алгоритм просмотра перечня товаров будет реализован в нескольких вариациях,
и тут мы можем воспользоваться паттерном "стратегия"

Плюсы:

	Можно динамически определять, какой алгоритм будет запущен,
соблюдается инкапсуляция - код алгоритмов отделен и скрыт от остального кода,
алгоритмы вызываются единообразно: без if-ов и других подобных конструкций

Минусы:

	Приходится создавать дополнительные структуры, что немного нагромождает код,
при проектировании системы нужно четко понимать, какие алгоритмы и когда применять,
в чем их отличия, чтобы не допустить логических ошибок

Таким образом, паттерн "Стратегия" обычно применяют для следующих задач:
	1) Нужно использовать разные виды одного алгоритма
	2) Приведение похожих структур в единую структуру
	3) Необходимо скрыть детали реализации алгоритмов
	4) Когда есть большое дерево условных операторов, где каждая ветвь представляет собой вариацию алгоритма
*/

type Strategy interface { //содержит общую для всех стратегий функцию выполнения

	Execute(data Data)
}

type StrategyA struct {
}

func (s StrategyA) Execute(data Data) { //выполняет алгоритм соответствующий стратегии А

	/*...*/
}

type Context struct { //позволяет сохранять и использовать определенную стратегию

	strategy Strategy
}

func (c *Context) SetStrategy(s Strategy) {
	c.strategy = s
}

func (c *Context) ExecStrategy(data Data) {
	c.strategy.Execute(data)
}

func StrategyClient() {
	context := new(Context)
	str := new(StrategyA)
	context.SetStrategy(str)
	context.ExecStrategy(Data{})
}
