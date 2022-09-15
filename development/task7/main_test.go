package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) { //для проведения теста: проверяем полученное от функции время
	//и сопоставляем  с тем, что выведет time.Now()
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	channels := []<-chan interface{}{sig(2 * time.Second), sig(5 * time.Second), sig(1 * time.Second), sig(1 * time.Second), sig(1 * time.Second)}
	or(
		channels...,
	)
	//получили канал, закрытый только при условии, что все каналы из chans уже закрылись
	//теперь проверяем, так ли это получилось или нет
	//сюда запоминаем информацию по каждому каналу
	isChannelClosed := []bool{false, false, false, false, false}
	_, ok := <-channels[0]
	if !ok {
		isChannelClosed[0] = true
	}
	_, ok = <-channels[1]
	if !ok {
		isChannelClosed[1] = true
	}
	_, ok = <-channels[2]
	if !ok {
		isChannelClosed[2] = true
		_, ok = <-channels[3]
		if !ok {
			isChannelClosed[3] = true
		}
		_, ok = <-channels[4]
		if !ok {
			isChannelClosed[4] = true
		}
	}
	for _, x := range isChannelClosed {
		if !x {
			t.Errorf("Тест провален")
		}
	}

}
