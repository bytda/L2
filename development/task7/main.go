package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===
Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.
Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}
Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}
start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)
fmt.Printf(“fone after %v”, time.Since(start))
*/

func or(channels ...<-chan interface{}) <-chan interface{} { //функция or-канал, объединяет каналы из слайса channels на входе,
	//в один канал, он дожидается, когда закрывается каждый из этих каналов и потом закрывается сам
	if len(channels) == 0 { //проверяем наличие элементов в слайсе
		return nil //возращаем nil
	}
	if len(channels) == 1 { //если на входе один канал
		return channels[0] //то просто его и возращаем
	}
	resultChannel := make(chan interface{}) //объявляем результирующий канал or
	go func() {                             //объявляем горутину, каналы могут закрываться долго,
		//поэтому мы сразу вернем or-канал,
		defer close(resultChannel) //а закроем канал по завершении  горутины
		<-channels[0]              //ждем завершения первого канала в слайсе
		<-or(channels[1:]...)      //вызываем снова функцию or, передавая в качестве аргумента слайс без первого элемента
		//получается рекурсия так как функция вызывает саму себя

	}()
	return resultChannel //возращаем итоговый канал
}
func main() {
	sig := func(after time.Duration) <-chan interface{} { //реализуем функцию из условия задания
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Second), //время уменьшили
		sig(5*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
		sig(1*time.Second),
	)
	fmt.Printf("Завершил работу спустя %v\n", time.Since(start)) //должен вывести время самого долгого выполнения,
	//то есть 5 сек, так как ждем выполнения всех горутин
}
