package main

import "testing"

func TestUnpacking(t *testing.T) {
	input := []string{`a4bc2d5e`, `abcd`, ``, `qwe\4\5`, `qwe\45`, `qwe\\5`}             //задаем вводные строки
	wantOutput := []string{`aaaabccddddde`, `abcd`, ``, `qwe45`, `qwe44444`, `qwe\\\\\`} //задаем верные ожидаемые строки
	for i, str := range input {                                                          //перебираем строки
		got, err := Unpacking(str) //вызываем нашу функцию
		if err != nil {
			t.Errorf("Получена некорректная строка") //обрабатываем ошибку
		}
		want := wantOutput[i] //вытаскиваем одну строку из ожидаемых
		if got != want {      //сравниваем
			t.Errorf("got: %s want: %s\n", got, want) //в случае не совпадения выводим ошибку
		}
	}
}
