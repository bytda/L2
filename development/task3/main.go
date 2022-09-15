package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Утилита sort
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание
и основные параметры): на входе подается файл из несортированными строками,
 на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов

*/

func main() {
	var n, r, u, m, c, b, h bool //объявляем булевые переменные для флагов
	var k int                    //переменная для флага с указанием строки
	flag.IntVar(&k, "k", 0, "Указание колонки для сортировки")
	flag.BoolVar(&n, "n", false, "Сортировать по числовому значению")
	flag.BoolVar(&r, "r", false, "Сортировать в обратном порядке")
	flag.BoolVar(&u, "u", false, "Не выводить повторяющиеся строки")
	//дополнительные ключи
	flag.BoolVar(&m, "M", false, "Сортировать по названию месяца")
	flag.BoolVar(&b, "b", false, "Игнорировать хвостовые пробелы")
	flag.BoolVar(&c, "c", false, "Проверять отсортированы ли данные")
	flag.BoolVar(&h, "h", false, "Сортировать по числовому значению с учетом суффиксов")

	flag.Parse() //ловим флаги

	input := flag.Arg(0)  //считываем первый аргумент на входе - исходный файл
	output := flag.Arg(1) //считываем второй аргумент на входе - отсортированный файл

	if input == "" || output == "" { //проверяем ввод пользователя
		panic("Укажите файл для чтения и файл для записи")
	}
	if k < 1 { //если указанная колонка меньше 1 то ее не указывали
		k = 0 //значит индекс будет нулевой то есть первая колонка
	} else { //если больше 1 то ее задали
		k-- //вычитаем единицу так как индексы начинаются с 0 а пользователь задает номер начиная с 1
	}

	var data [][]string              //двумерный слайс для данных из файла
	data = readFile(input)           //записываем данные
	var SortFunc func(i, j int) bool //переменная для функции сортировки
	switch true {                    //выбираем сортировку по ключам
	case n: //если ключ "n"
		SortFunc = func(i, j int) bool { //один из методов сортировки в го с помощью функции func(i, j int) bool
			//дает возможность задавать свои условия сортировки
			a, _ := strconv.ParseFloat(getDataElement(data, i, k), 64) //переводим string в float для сортировки по числовому значению
			b, _ := strconv.ParseFloat(getDataElement(data, j, k), 64)
			if r { //если активен ключ "r"
				return a > b //возращаем в обратном порядке
			}
			return a < b //возращаем сортировку
		}
	case m: //если активен ключ "m"
		SortFunc = func(i, j int) bool {
			if r { //если активен ключ "r" выводим в обратном порядке
				return getMonth(getDataElement(data, j, k)).Before(getMonth(getDataElement(data, i, k)))
			}
			return getMonth(getDataElement(data, i, k)).Before(getMonth(getDataElement(data, j, k))) //сортируем по месяцам
		}
	case h: //если активен ключ "h"
		SortFunc = func(i, j int) bool { //то сортируем по количеству символов в строке
			if r { //если активен ключ "r" выводим в обратном порядке
				return getLen(data[i][k:]) > getLen(data[j][k:]) //по убыванию
			}
			return getLen(getDataElementSlice(data, i, k)) < getLen(getDataElementSlice(data, j, k)) //по возрастанию
		}
	default: //обычная сортировка
		SortFunc = func(i, j int) bool {
			if r {
				return getDataElement(data, i, k) > getDataElement(data, j, k)
			}
			return getDataElement(data, i, k) < getDataElement(data, j, k)
		}
	}
	sort.Slice(data, SortFunc) //сортируем согласно логике в зависимости от ключа

	if c { // если активен ключ проверки сортировки, то проверяем файл на упорядоченность
		// согласно указанному правилу (ключу), SortFunc получает логику сортировку в зависимости от ключа сортировки

		isSorted := sort.SliceIsSorted(data, SortFunc)
		fmt.Println("Отсортировано?", isSorted)
		return
	}

	WriteFile(data, output) //записываем отсортированные данные в файл
}

func readFile(directory string) (result [][]string) { //считывает данные из файла и преобразует в двумерный слайс (строки, колонки)
	file, err := os.Open(directory) //открываем файл по пути указанному пользователем
	if err != nil {
		panic(err) //обрабатываем ошибку
	}
	defer file.Close()                //если выйдет паника перед завершением программы закроем файл
	scanner := bufio.NewScanner(file) //объявляем новый сканер
	for scanner.Scan() {
		var words []string                         //срез для слов
		words = strings.Split(scanner.Text(), " ") // разделяем строки на слова (колонки)
		result = append(result, words)             //записываем в итоговый срез
	}
	return result //возращаем итоговый срез

}
func getDataElement(data [][]string, i, k int) string { //возвращает элемент с индексом i,k если он есть, если его нет, возвращает пустую строку
	if k < len(data[i]) { //если индекс "к" меньше длинны среза i то он входит в этот срез следовательно сущетсвует
		return data[i][k] //возращаем этот элемент
	}
	return "" //в противном случае выдаем пустую строку
}

func getMonth(month string) time.Time { //определяем месяц по строке
	if t, err := time.Parse("January", month); err == nil { //парсим строку, если встретили месяц январь
		return t //возращаем его
	}
	if t, err := time.Parse("Jan", month); err == nil {
		return t
	}
	if t, err := time.Parse("1", month); err == nil {
		return t
	}
	if t, err := time.Parse("01", month); err == nil {
		return t
	}
	return time.Time{} //возращаем пустую строку если не было совпадений
}

func getDataElementSlice(data [][]string, i, k int) []string { //возращает элементы строки "i" начиная с индекса "k"
	if k < len(data[i]) { //если они есть
		return data[i][k:] // "к:" означает что берем элементы начиная с индекса "к" и до конца среза
	}
	return []string{} //если нет возращает пустой срез строк
}
func getLen(str []string) int { //возращает количество символов в строке
	var result int        //переменная для результата
	result = len(str) - 1 // учитываем пробелы между слов, если например слов 5,
	// то длина массива 5, а пробелов между слов на один меньше,то есть из длины вычитаем один
	for _, s := range str { //перебираем символы слов в строке
		result += len(s) //количество символов суммируем
	}
	return result
}

func WriteFile(data [][]string, directory string) { //записывает данные в файл
	file, err := os.Create(directory) // создаем (или перезаписываем) файл и открываем по указанному пути
	if err != nil {                   //обрабатываем ошбику
		panic(err)
	}
	defer file.Close()                 //при заверешнии функции закрываем файл
	lines := make([]string, len(data)) //срез длинной как входные данные
	for i, da := range data {          //перебираем входные данные
		lines[i] = strings.Join(da, " ") //объдиняем слова в строки разделяя пробелом
	}
	_, err = file.WriteString(strings.Join(lines, "\n")) //записываем строчки в текст разделяя переносом строки
	if err != nil {                                      //обрабатываем ошибку
		panic(err)
	}
}
