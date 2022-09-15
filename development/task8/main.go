package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/
func main() {
	scanner := bufio.NewScanner(os.Stdin) //чтение со стандартного потока ввода
	fmt.Println(`Для выхода введите команду "exit"`)
	for { //бесконечный цикл с "прослушкой" команд терминала
		currentDir, _ := os.Getwd()    //получаем текущую директорию
		fmt.Printf("%s> ", currentDir) //выводим предложение к вводу команд с указанием текущей директории
		ok := scanner.Scan()
		if !ok {
			log.Fatal(errors.New("Ошибка при вводе команды"))
		}
		input := scanner.Text() //считываем ввод пользователя
		if input == "exit" {    //если ввели строку exit - выходим из всей программы
			fmt.Println("Завершаем программу...")
			os.Exit(0)
		}
		comms := strings.Split(input, "|") // разделитель для использования множественных команд в одной строке
		err := Commands(comms)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Commands(comand []string) error {
	for _, x := range comand { //обходим список команд
		args := strings.Split(x, " ") //каждая команда может содержать список аргументов через пробел
		commandName := args[0]        //первое слово в строке - имя самой команды
		switch commandName {
		case "cd":
			if len(args) < 2 { //если agrs меньше 2 элементов, значит, там только команда
				return errors.New("Не хватает аргументов")
			}
			dir := args[1]                             //новая директория
			fmt.Println("Сменяем директорию на ", dir) //выводим в консоль
			os.Chdir(dir)                              //меняем текущую директорию на указанную
		case "pwd":
			dir, _ := os.Getwd()                     //получаем путь к директории
			fmt.Println("Текущая директория: ", dir) //выводим директорию
		case "echo":
			if len(args) < 2 { //если agrs меньше 2 элементов, значит, там только команда
				return errors.New("Не хватает аргументов")
			}
			for i := 0; i < len(args); i++ {
				fmt.Fprintf(os.Stdout, args[i+1]+" ") //записываем аргументы в поток вывода
			}
		case "kill":
			if len(args) < 2 { //если agrs меньше 2 элементов, значит, там только команда
				return errors.New("Не хватает аргументов")
			}
			err := exec.Command("kill", args[1]).Run()
			if err != nil {
				return err
			}
		case "ps":
			procs, err := ps.Processes() //возращаем процессы ОС
			if err != nil {
				return err
			}
			fmt.Println("Процессы: ")
			for _, v := range procs { //перебираем процессы
				fmt.Println(v.Pid()) //и выводим по Id
			}
		default: //неопознанная команда, невошедшая в список выше
			return errors.New("Неизвестная команда")

		}
	}
	return nil
}