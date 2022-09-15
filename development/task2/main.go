package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

/*
Задача на распаковку
Создать Go-функцию, осуществляющую примитивную распаковку строки,
содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)


В случае если была передана некорректная строка,
функция должна возвращать ошибку. Написать unit-тесты.
*/

func main() {
	result, err := Unpacking(`qwe\4\5 `)
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println("Ошибка")
		log.Fatal(err)
	}
}

func Unpacking(str string) (string, error) { //объявляем функцию распаковки
	var strBuild strings.Builder //используем билдер, он быстрее чем конкатенация с помощью +(вопрос из L1)
	strSlice := []rune(str)      //разбиваем строку на слайс рун, чтобы работать с символами
	fmt.Println(strSlice)
	var currentRune *rune                //создаем указатель где будем хранить текущий символ слайса
	currentRune = nil                    //определяем начальное состояние
	for i := 0; i < len(strSlice); i++ { //обходим слайс рун(символов)
		currentChar := strSlice[i] //получаем текущий символ
		if currentChar == '\\' {   //если текущий символ это символ экранирования, записываем как \\ но там \
			if i+1 >= len(strSlice) { //если \ - последний символ в строке
				return "", errors.New("Некорретная строка") //то выводим ошибку
			}
			i++                     //а так пропускаем этот символ,переходим к следующему значимому символу
			if currentRune != nil { // если у нас уже есть сохраненный символ
				strBuild.WriteRune(*currentRune) //прописываем его один раз
			} else {
				currentRune = new(rune) //иначе выделяем память под руну
			}
			*currentRune = strSlice[i] //запоминаем новый символ после символа экранирования
			continue                   //переходим к следующей итерации цикла,
			// то есть в начало, не проходя проверки ниже
		}
		if !unicode.IsDigit(currentChar) { //если попался символ
			if currentRune != nil { // если у нас уже есть сохраненный символ
				strBuild.WriteRune(*currentRune) //записываем его один раз
			} else {
				currentRune = new(rune) //выделяем память под руну
			}
			*currentRune = currentChar // запоминаем текущий символ
		} else { //если попалась цифра
			count, _ := strconv.Atoi(string(currentChar)) //переводим символ(кол-во символов) в строку а потом строку в число
			if currentRune == nil {                       //так как nil может быть если до этого тоже было число или это первый символ
				return "", errors.New("Неверная строка") //это некорректная строка (как в случае с примером "45")
			}
			for j := 0; j < count; j++ { //число определяет количество повторений буквы
				strBuild.WriteRune(*currentRune) //повторяем букву кол-во раз равное числу
			}
			currentRune = nil //обнуляем
		}

	}
	if currentRune != nil { //если остался символ, который мы сохранили ("abcd" - d будет таким символом)
		strBuild.WriteRune(*currentRune) //запишем один раз
	}
	return strBuild.String(), nil //выводим итоговую строку и пустую ошибку
}
