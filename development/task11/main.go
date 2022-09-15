package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

/*
HTTP-сервер

Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать
строго со стандартной HTTP-библиотекой.

В рамках задания необходимо:
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
Реализовать вспомогательные функции для парсинга и валидации параметров
методов /create_event и /update_event.
Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные
функции и объекты доменной области.
Реализовать middleware для логирования запросов

Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON-документ содержащий
либо {"result": "..."} в случае успешного выполнения метода, либо {"error": "..."} в случае
ошибки бизнес-логики.

В рамках задачи необходимо:
Реализовать все методы.
Бизнес логика НЕ должна зависеть от кода HTTP сервера.
В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных
данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок
сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге
и выводить в лог каждый обработанный запрос.
*/
var mutex *sync.Mutex //объявляем мьютекс для защиты записи в переменную

type Event struct { //новый тип структуры с параметрами события
	IdEvent int       `json:"event_id"` // уникальный номер
	Content string    `json: "content"` //содержание
	Date    time.Time `json:"date"`     //дата

}

var events []Event //срез для событий

type Logger struct { // обертка для хандлера
	handler http.Handler
}

func (l *Logger) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	//обертка для переопределения метода ServeHTTP, добавили логгирование
	initTime := time.Now()
	l.handler.ServeHTTP(w, r)
	endTime := time.Now()
	difference := endTime.Sub(initTime)
	log.Printf("%s %s %v", r.Method, r.URL.Path, difference) //вывод время обработки хендлера
}

func main() {
	port := os.Getenv("PORT") //получаем порт из окружения

	httpMultyPlex := http.NewServeMux() //мультиплексор, который выбирает обработчик в зависимости от запроса
	httpMultyPlex.HandleFunc("/create_event", CreateEvent)
	httpMultyPlex.HandleFunc("/update_event", UpdateEvent)
	httpMultyPlex.HandleFunc("/delete_event", DeleteEvent)
	httpMultyPlex.HandleFunc("/events_for_day", EventsDay)
	httpMultyPlex.HandleFunc("/events_for_week", EventsWeek)
	httpMultyPlex.HandleFunc("/events_for_month", EventsMonth)

	midleLogger := WrapHandler(httpMultyPlex)
	log.Fatalln(http.ListenAndServe(port, midleLogger.handler))
}

func WrapHandler(h http.Handler) *Logger { //оборачиваем стандартный handler
	var wrapLoger Logger
	wrapLoger.handler = h
	return &wrapLoger
}

func JSONparsing(r *http.Request) (Event, error) { //парсинг ответа от сервера
	var event Event                               // перемена типа события
	err := json.NewDecoder(r.Body).Decode(&event) //читаем и сохраняем в переменную
	if err != nil {
		return event, errors.New("Некорректный файл Json") //обрабатываем ошибку
	}
	return event, nil //возращаем результат и нулевую ошибку
}

func errorResponse(w http.ResponseWriter, e string, status int) {
	errorResponse := struct {
		Error string `json:"error"`
	}{Error: e}
	json, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func resultResponse(w http.ResponseWriter, r string, e []Event, status int) {
	resultResponse := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{Result: r, Events: e}

	json, err := json.Marshal(resultResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ValidEvent(event Event) error { //валидация события на содержание данных
	if event.IdEvent <= 0 || event.Content == "" { //если номер меньше или равен 0 или содержание пустое
		return errors.New("Некорректное событие") //выводим ошибку
	}
	return nil //а так выыводим пустую ошибку
}

func CreateNewEvent(event Event) error { //создание нового события
	mutex.Lock()               //блокируем доступ к данным для остальных процессов
	defer mutex.Unlock()       //по завершении работы функции разблокируем
	for _, x := range events { //перебираем срез событий
		if x.IdEvent == event.IdEvent { //если номера событий совпали
			return errors.New("Событие уже существует") //выводим ошибку
		}
	}
	events = append(events, event) //а так добавляем событие в срез
	return nil                     //возращаем пустую ошибку
}

func CreateEvent(w http.ResponseWriter, r *http.Request) { //обработчик для create_event
	if r.Method != http.MethodPost { //если это не post-метод, не обрабатываем его
		errorResponse(w, "не POST-метод", http.StatusBadRequest)
		return
	}
	newEvent, err := JSONparsing(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
		return
	}
	err = ValidEvent(newEvent) //проходим валидацию
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
	}
	if err := CreateNewEvent(newEvent); err != nil { //пройдя проверки, создаем событие
		errorResponse(w, err.Error(), http.StatusBadRequest)
	}
	resultResponse(w, "Успешно", []Event{newEvent}, http.StatusCreated)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) { //обработчик обновления события
	if r.Method != http.MethodPost { //если это не POST-метод не обрабатываем его
		errorResponse(w, "не POST-метод", http.StatusBadRequest)
		return
	}
	newEvent, err := JSONparsing(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
		return
	}
	err = ValidEvent(newEvent) //проводим валидацию
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
		return
	}
	mutex.Lock()               //блокируем доступ чтобы другие процессы не вносили
	defer mutex.Unlock()       //как функция отработает разблокируем доступ
	for i, x := range events { //перебираем события
		if x.IdEvent == newEvent.IdEvent { //если нашли событие которое нам надо
			events[i] = newEvent //меняем это событие на новое
			return
		}
	}
	return
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //если это не post-метод,  не обрабатываем его
		errorResponse(w, "not POST-method", http.StatusBadRequest) //обрабатываем ошибку
		return
	}
	newEvent, err := JSONparsing(r) //парсим json
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
		return
	}
	err = ValidEvent(newEvent) //проводим валидацию
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest) //обрабатываем ошибку
		return
	}

	mutex.Lock()               //пока меняем слайс, другие потоки не должны туда влезать
	defer mutex.Unlock()       //как функция отработает разблокируем доступ
	for i, x := range events { //перебираем события
		if x.IdEvent == newEvent.IdEvent { //если нашли которое нам нужно
			events = append(events[:i], events[i+1:]...) //вырезаем которое нужно удалить, а
			//точнее не берем его в диапазон нового среза
		}
	}
	return
}

func EventsByDay(date time.Time) []Event { //список событий cj с полным совпадением дат
	var result []Event //срез для результата
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events { //перебираем события
		if x.Date.Year() == date.Year() && x.Date.Month() == date.Month() && x.Date.Day() == date.Day() {
			//если полное совпадение даты
			result = append(result, x) //добавляем в итогой срез
		}
	}
	return result //выводим результат
}

func EventsByWeek(date time.Time) []Event { //список событий с совпадением до недели
	var result []Event //срез для результата
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events { //перебираем события
		difference := date.Sub(x.Date) //определяем разницу между датами
		if difference < 0 {            //если отрицательное
			difference = -difference //возращаем абсолютное значение
		}
		if difference <= time.Duration(7*24)*time.Hour { //если разница меньше чем неделя (в часах)
			result = append(result, x) //то добавляем
		}
	}
	return result //возращаем результат

}

func EventsByMonth(date time.Time) []Event { //список событий с совпадением дат вплоть до месяца
	var result []Event //срез итоговый
	mutex.Lock()
	mutex.Unlock()
	for _, x := range events { //перребираем события
		if x.Date.Year() == date.Year() && x.Date.Month() == date.Month() { //если совпадает по году и месяцу
			result = append(result, x) //добавляем
		}
	}
	return result //выводим результат

}

func EventsDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "Не GET-метод", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := EventsByDay(date)
	resultResponse(w, "Успешно", result, http.StatusOK)
}

func EventsWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "Не GET-метод", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := EventsByWeek(date)
	resultResponse(w, "Успешно", result, http.StatusOK)
}
func EventsMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "Не GET-метод", http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2019-09-09", r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := EventsByMonth(date)
	resultResponse(w, "Успешно", result, http.StatusOK)
}
