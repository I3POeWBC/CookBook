package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	connFlag        *string        = flag.String("conn", "postgres://master@localhost:5432/app", "Строка подключения к серверу в формате postgres://username<:password>@<address><:port>/<database name>.")
	connCfgFlag     *string        = flag.String("conn-cfg", "", "Если определено, то указывает на YAML файл с описание строки соединения.")
	cmdFlag         *string        = flag.String("cmd", "", "Если не пусто, то команда для выполнения.")
	scriptFlag      *string        = flag.String("script", "script.sql", "Файл со скриптом.")
	sepFlag         *string        = flag.String("separator", "\\", "Разделитель запросов в скрипте.")
	timeoutFlag     *time.Duration = flag.Duration("", time.Second*10, "Таймаут отдельных операций.")
	maxRowCountFlag *uint          = flag.Uint("max-row", 0, "Число сохраняемых строк запросов select. Если 0, то без ограничений.")
	outDirFlag      *string        = flag.String("output-dir", "output", "Каталог для сохранения результатов запросов. Очищается перед выполнение скрипта. Запросы нумеруются в порядке очередности.")
)

func main() {

	flag.Parse()

	mustNewLogger("log/console-client.log")

	if *scriptFlag != "" {
		if bs, err := os.ReadFile(*scriptFlag); err != nil {
			log.Printf("ошибка открытия файла [%s] err:[%v]", *scriptFlag, err)
		} else {
			execSript(string(bs))
		}
	}

}
