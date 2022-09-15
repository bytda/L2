package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/
func main() {
	scanner := bufio.NewScanner(os.Stdin) //объявляем новый сканнер ввода в консоль
	fmt.Println("Введите url")            //обращаемся к пользователю

	ok := scanner.Scan() //проверяем работу сканера
	if !ok {             //если все не ок то
		log.Fatal("Ошибка") //выводим ошибку
	}
	url := scanner.Text() //текст сканера кладем в переменную с ссылкой
	err := wget(url)      //вызываем функцию
	if err != nil {
		log.Fatal(err) //обрабатываем ошибку
	}
}

func wget(url string) error {
	response, err := http.Get(url) //запрашиваем ответ по протоколу http по данному url
	if err != nil {
		return errors.New("Ошибка по данному url") //обрабатываем ошибку
	}
	temp := strings.Split(url, "/")      //разделяем строку на входе с разделителем /
	fileName := temp[len(temp)-1]        //находим имя файла в конце строки заданной пользователем
	saveFile, err := os.Create(fileName) //создаем файл для сохранения скачанного
	if err != nil {
		return errors.New("Ошибка при создании файла") //обрабатываем ошибку
	}
	defer saveFile.Close()                    //после отработки функции закрываем файл
	_, err = io.Copy(saveFile, response.Body) //копируем файл из тела ответа на запрос в файл который мы создали
	if err != nil {
		return errors.New("Ошибка при сохранении файла") //обрабатываем ошибку
	}
	fmt.Println("Файл сохранен") //выводим в консоль
	return nil                   //возращаем нулевую ошибку если функция успешно до сюда дошла
}
