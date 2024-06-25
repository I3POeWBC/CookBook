package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
	"time"

	"log"

	"github.com/jackc/pgx/v5"
)

const (
	SEL  = "SEL"
	EXEC = "EXEC"
)

type OutputExec struct {
	Index    int           `json:"index"`
	Cmd      string        `json:"cmd"`
	Error    error         `json:"error"`
	Duration time.Duration `json:"duration"`
	Type     string        `json:"type"`
	Affected string        `json:"affected"`
}

var (
	conn   *pgx.Conn
	cancel context.CancelFunc
)

func execSript(script string) {
	now := time.Now()

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("перехват паники [%v]", r)
		}
		log.Printf("Выполнение скрипта завершено [%v]", time.Since(now))
	}()

	defer func() {
		if cancel != nil {
			cancel()
			conn = nil
		}
	}()

	if script == "" {
		log.Printf("Пустой скрипт. Обработка завершена")
		return
	}

	if err := cleanOutput(); err != nil {
		log.Printf("ошибка очистки каталога для результатов запроса err: [%v]", err)
		return
	}

	chunks := strings.Split(script, "/")

	reIsSelect := regexp.MustCompile(`^select.+`)

	var err error
	if conn == nil {
		if err = openConnection(*connFlag); err != nil {
			log.Printf("ошибка подключения к базе err: [%v]", err)
			return
		}

	}

	for k, v := range chunks {

		var queryType string = EXEC
		if reIsSelect.MatchString(v) {
			queryType = SEL
		}

		v = strings.Trim(v, " \t\n\r")
		log.Printf("Выполняется [%d:%s] элемент скрипта: [%s]", k, queryType, Crop(v, 30))

		switch queryType {
		case SEL:
		case EXEC:
			execCmd(v)
		}
	}

}

func execCmd(cmd string) (err error) {

	ctxExec, cancelExe := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancelExe()

	tag, err := conn.Exec(ctxExec, cmd)
	if err != nil {
		log.Printf("ошибка выполнения [%v]", err)
		return
	}

	txt := ""
	if tag.Insert() {
		txt = "Insert"
	}
	if tag.Delete() {
		txt = "Delete"
	}
	if tag.Select() {
		txt = "Select"
	}
	if tag.Update() {
		txt = "Update"
	}

	log.Printf("OK [%s:%d] [%v]", txt, tag.RowsAffected(), tag)

	return
}

func openConnection(connString string) (err error) {

	if cancel != nil {
		cancel()
		conn = nil
	}

	var cfg *pgx.ConnConfig

	if cfg, err = pgx.ParseConfig(connString); err != nil {
		err = fmt.Errorf("openConnection: [%v]", err)
		return
	}

	var ctx context.Context
	ctx, cancel = context.WithTimeout(context.Background(), *timeoutFlag)

	cfg.TLSConfig = nil

	if conn, err = pgx.ConnectConfig(ctx, cfg); err != nil {
		err = fmt.Errorf("не удалось подключиться строкой [%s] err: [%v]", connString, err)
		return
	}

	return

}

func cleanOutput() (err error) {

	//log.Printf("cleanOutput fs.ErrExist): [%v], errors.Is(err, fs.ErrNotExist: [%v]\n", errors.Is(err, fs.ErrExist), errors.Is(err, fs.ErrNotExist))

	if _, err = os.Stat(*outDirFlag); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(*outDirFlag, 0766); err != nil {
			return
		}
	}
	log.Printf("cleanOutput fs.ErrExist): [%v], errors.Is(err, fs.ErrNotExist: [%v]\n", errors.Is(err, fs.ErrExist), errors.Is(err, fs.ErrNotExist))

	if de, errRead := os.ReadDir(*outDirFlag); err != nil {
		err = errRead
		return
	} else {
		for _, v := range de {
			log.Printf("cleanOutput %s\n", v.Name())
		}

	}

	return
}
