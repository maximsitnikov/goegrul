# goegrul
Package for obtaining data about an organization from EGRUL (Russian register of legal organizations). 

Пакет предназначен для скачивания данных о юрлицах и ИП по ИНН с сайта ЕГРЮЛ https://egrul.nalog.ru/index.html

## Установка

```bash
$ go get -u github.com/maximsitnikov/goegrul
```

## Использование

```go
package main

import (
  "fmt"
  
  "github.com/maximsitnikov/goegrul"
)

func main() {
	var inn string = "7707083893"

	f, err := goegrul.GetDataFromEGRUL(inn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if f.Yurik {
		//если юрлицо
		fmt.Println("Полное наименование: " + f.FullName)
		fmt.Println("Адрес: " + f.Address)
		fmt.Println("КПП: " + f.KPP)
		fmt.Println("Краткое наименование: " + f.Name)
	} else {
		//если физлицо
		fmt.Println("Полное наименование: " + f.FullName)
	}

}
```
