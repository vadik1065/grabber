package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// выводим ошибку
func shownError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// делает валидное имя
func doValidName(nameString *string) {
	delSymbols := [4]string{"www.", "https://", "http://", ".html"}
	for _, delStr := range delSymbols {
		*nameString = strings.ReplaceAll(*nameString, delStr, "")
	}
	*nameString = strings.ReplaceAll(*nameString, "/", "-")
	*nameString = strings.TrimSuffix(*nameString, "-")

}

// скачиваем страницу
func downHtm(namePage string, c chan string) {
	// fmt.Println("start " + namePage)
	http, err := http.Get(namePage)
	if err == nil {
		doValidName(&namePage)
		body, err := ioutil.ReadAll(http.Body)
		direct := flag.Arg(1)
		formatF := "html"
		err = ioutil.WriteFile(direct+namePage+"."+formatF, body, 0644)
		shownError(err)
	}
	shownError(err)
	// fmt.Println("end " + namePage)
	c <- "end"
}

// основная функция
func main() {
	//парсим флаги
	flag.Parse()

	// fmt.Println(flag.Arg(0))
	// fmt.Println(flag.Arg(1))

	// чтение файла
	file, err := os.Open(flag.Arg(0))
	shownError(err)

	fileScaner := bufio.NewScanner(file)
	fileScaner.Split(bufio.ScanLines)

	// пробегаем по всем строкам
	var c chan string // объявляем канал
	for fileScaner.Scan() {
		c = make(chan string) // переприсваивание канала
		puthPage := fileScaner.Text()
		go downHtm(puthPage, c)
	}
	<-c // дожидаемся выполнение последнего, т.к он пересвоен

	defer file.Close()
}
