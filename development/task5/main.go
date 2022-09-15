package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
Утилита grep

Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки

*/

func main() {
	var a, b, c int
	var count, ignore, invert, fixed, lnum bool

	flag.IntVar(&a, "A", 0, "печатать +N строк после совпадения") //определяем флаги
	flag.IntVar(&b, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&c, "C", 0, "печатать ±N строк вокруг совпадения")
	flag.BoolVar(&count, "c", false, "количество строк")
	flag.BoolVar(&lnum, "n", false, "печатать номер строки")

	flag.BoolVar(&fixed, "F", false, "точное совпадение со строкой")
	flag.BoolVar(&invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&ignore, "i", false, "игнорировать регистр")

	flag.Parse()

	str := flag.Arg(0) //считываем строку на входе
	if len(str) == 0 {
		panic("Некорректный ввод строки")
	}
	input := flag.Arg(1) // считываем файл для чтения
	if len(input) == 0 {
		panic("Некорректный ввод файла")
	}
	data := readFileToStrings(input) // считываем данные из файла

	prepareStrings := func(str, substr string) (string, string) { // возвращает неизмененные значения
		return str, substr
	}
	checkContain := func(str, substr string) bool { // проверяет, содержит ли строка подстроку
		return strings.Contains(str, substr)
	}
	handleCheck := func(check bool) bool { // возвращает неизмененное значение
		return check
	}

	if ignore { // если активен ключ i
		prepareStrings = func(str, substr string) (string, string) {
			return strings.ToLower(str), strings.ToLower(substr) // приводит обе строки к нижнему регистру
		}
	}
	if fixed { // если активен ключ F
		checkContain = func(str, substr string) bool { // проводит прямое сравнение строк
			if str == substr {
				return true
			}
			return false
		}
	}
	if invert { // если активен ключ v
		handleCheck = func(check bool) bool {
			return !check // инвертирует значение
		}
	}
	var conIdx []int // инициализируем слайс, который будет хранить значения индексов строк
	// удовлетворяющих условиям

	for i, line := range data { // обходим все строки
		if handleCheck(checkContain(prepareStrings(line, str))) { // если строка удовлетворяет всем условиям добавляем ее индекс в слайс
			conIdx = append(conIdx, i)
		}
	}
	switch true {
	// обработка ключа A
	case a > 0:
		for _, idx := range conIdx { // обходим все сохраненные индексы
			for i := 0; i < a; i++ { // выводим в STDOUT A строк стоящих перед строкой с сохраненным индексом
				if idx-a+i >= 0 {
					fmt.Println(data[idx-a+i])
				}
			}
		}
	case b > 0: // обработка ключа B
		for _, idx := range conIdx {
			for i := 0; i < b; i++ {
				if idx+i < len(data) {
					fmt.Println(data[idx+i])
				}
			}
		}
	case c > 0:
		for _, idx := range conIdx {
			for i := 0; i < c; i++ {
				if idx-c+i >= 0 {
					fmt.Println(data[idx-c+i])
				}
				if idx+i < len(data) {
					fmt.Println(data[idx+i])

				}
			}
		}
	case count:
		fmt.Println(len(conIdx))
	case lnum:
		for _, idx := range conIdx {
			fmt.Println(idx + 1)
		}
	default:
		for _, idx := range conIdx {
			fmt.Println(data[idx])
		}
	}

}

func readFileToStrings(dir string) (result []string) { //функция для чтения файла и записи его в срез строк
	file, err := os.Open(dir) //открываем файл
	if err != nil {           //проверяем ошибку
		panic(err)
	}
	defer file.Close()                //закрываем подключение по завершении
	scanner := bufio.NewScanner(file) //объявляем новый сканер
	for scanner.Scan() {              //во время сканирования
		result = append(result, scanner.Text()) //добавляем в срез
	}
	return result //возращаем результат
}
