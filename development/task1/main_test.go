package main

import (
	"testing"
	"time"
)

func TestGetTime(t *testing.T) { //функция теста сверяет время с сервера и время системы
	got := getTime()            //получаем время с сервера
	want := time.Now()          //получаем время системы
	difference := got.Sub(want) //вычитаем чтобы увидеть разницу
	if difference < 0 {         //если разница отрицательная
		difference = -difference //получаем ее абсолютное значение, то есть убираем знак
	}
	if difference > time.Second { //в случае если разница превышает секунду
		t.Errorf("getTime(): %q\ntime.Now(): %q\nDifference: %q", got, want, difference) //выводим отчет
	}

}
