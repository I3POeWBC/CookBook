package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	SEL  = "SEL"
	EXEC = "EXEC"
)

type IIndex interface {
	GetIndex() int
}

type OutputBase struct {
	Index    int           `json:"index"`
	Cmd      string        `json:"cmd"`
	Error    error         `json:"error"`
	Duration time.Duration `json:"duration"`
	Type     string        `json:"type"`
}

func (base *OutputBase) GetIndex() int {
	return base.Index
}

type OutputExec struct {
	OutputBase
	Affected string `json:"affected"`
}

type OutputSelect struct {
	OutputBase
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
			execCmd(k, v)
		}
	}

}

func execCmd(index int, cmd string) (err error) {

	ctxExec, cancelExe := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancelExe()

	var (
		now time.Time = time.Now()
		tag pgconn.CommandTag
	)

	defer func() {
		output := &OutputExec{
			OutputBase: OutputBase{
				Index:    index,
				Cmd:      cmd,
				Error:    err,
				Duration: time.Since(now),
				Type:     EXEC,
			},
			Affected: fmt.Sprint(tag.RowsAffected()),
		}

		saveCmdResult(output)

	}()

	tag, err = conn.Exec(ctxExec, cmd)
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

	if connString, err = parseConnectionString(connString); err != nil {
		return
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

func saveCmdResult(a any) (err error) {

	ii := a.(IIndex)

	filePath := filepath.Join(*outDirFlag, fmt.Sprint(ii.GetIndex())) + ".json"

	var bs []byte
	if bs, err = json.MarshalIndent(a, "", "  "); err != nil {
		err = fmt.Errorf("saveCmdResult ошибка маршалинга err: [%v]", err)
		return
	} else if err = os.WriteFile(filePath, bs, 0644); err != nil {
		err = fmt.Errorf("saveCmdResult ошибка записи файла err: [%v]", err)
		return
	}
	return
}
