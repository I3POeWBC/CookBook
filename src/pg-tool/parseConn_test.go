package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

type Sample struct {
	Value  string
	Assert bool
	Note   string
}

// postgres://username<:password>@<address><:port>/<database name>
func TestParseConnection(t *testing.T) {

	var (
		samples []Sample = []Sample{
			//"postgres://(username:password@)(address:port)(/database_name)"
			//postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
			{"postgres://username:password@address:123/database_name", true, "Вариант со всеми элементами, адрес текст"},
			{"postgres://username:password@127.0.0.1:123/database_name", true, "Вариант со всеми элементами, адрес IP"},
			{"postgres://username@address:345/database_name", true, "Вариант без пароля"},
			{"postgres://username@address:345", true, "Вариант без пароля и без базы"},
			{"postgres://username@address:port", false, "Вариант с ошибкой в описании порта"},
			{"postgres://username@address", true, "Вариант с логином и адресом сервера"},
		}
	)

	mustNewLogger("log/test.log")

	//reStr := `^postgres://(^[:%/()@]+)(.+)@(.+)([:][0-9]+)?([/].+)$`
	// НЕ -:%/()@ [^-:%/()@]
	reStr := `^postgres://(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/](.+))?$`
	re := regexp.MustCompile(reStr)
	_ = samples

	for k, v := range samples {
		r := re.MatchString(v.Value)
		if r != v.Assert {
			t.Fatalf("для cтроки [%s] получен не ожидаемый результат  [%v] вместо [%v]", v.Value, r, v.Assert)
		} else {

			log.Printf("[%d] %v [%s]\n", k, v.Value, v.Value)

			ss := re.FindStringSubmatch(v.Value)
			for kk, vv := range ss {
				log.Printf("%-2d %s\n", kk, vv)
			}

		}
	}

}

func TestParseConfigConnectionW(t *testing.T) {

	var (
		samples []Sample = []Sample{
			//"postgres://(username:password@)(address:port)(/database_name)"
			//postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
			{"postgres://username:password@address:123/database_name", true, "Вариант со всеми элементами, адрес текст"},
			{"postgres://username:password@127.0.0.1:123/database_name", true, "Вариант со всеми элементами, адрес IP"},
			{"postgres://username@address:345/database_name", true, "Вариант без пароля"},
			{"postgres://username@address:345", true, "Вариант без пароля и без базы"},
			{"postgres://username@address:port", false, "Вариант с ошибкой в описании порта"},
			{"postgres://username@address", true, "Вариант с логином и адресом сервера"},
		}
	)

	mustNewLogger("log/test.log")

	//reStr := `^postgres://(^[:%/()@]+)(.+)@(.+)([:][0-9]+)?([/].+)$`
	// НЕ -:%/()@ [^-:%/()@]
	reStr := `^postgres://(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/].+)?$`
	re := regexp.MustCompile(reStr)
	_ = samples

	for k, v := range samples {
		r := re.MatchString(v.Value)
		if r != v.Assert {
			t.Fatalf("для cтроки [%s] получен не ожидаемый результат  [%v] вместо [%v]", v.Value, r, v.Assert)
		} else {

			log.Printf("[%d] %v [%s]\n", k, v.Value, v.Value)

			if _, err := pgx.ParseConfig(v.Value); err != nil {
				log.Printf("\tParseConfig fail: [%v]\n", err)
			} else {
				log.Printf("\tParseConfig OK\n")
			}

		}
	}

}

// https://github.com/gravitational/teleport/discussions/26217
func TestConnection(t *testing.T) {

	var (
		samples []Sample = []Sample{
			//{"postgres://master:xx1234@localhost:5432/postgress", true, "Вариант со всеми элементами, адрес текст"},
			{"postgres://master:xx1234@localhost:5432/app", true, "Вариант со всеми элементами, адрес текст"},
		}
	)

	mustNewLogger("log/test.log")

	reStr := `^postgres://(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/].+)?$`
	re := regexp.MustCompile(reStr)

	timeout := time.Second * 5

	for k, v := range samples {
		if !re.MatchString(v.Value) {
			t.Fatalf("[%d] строка подключения [%s] не прошла проверку", k, v.Value)
		} else {

			//if cfg, err := parseConnectionString(v.Value, timeout); err != nil {
			if cfg, err := pgx.ParseConfig(v.Value); err != nil {
				t.Fatalf("parseConnectionString: [%v]", err)
			} else {
				func() {
					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					cfg.TLSConfig = nil

					if conn, err := pgx.ConnectConfig(ctx, cfg); err != nil {
						t.Fatalf("[%d] (%s) не удалось подключиться строкой [%s] err: [%v]", k, v.Note, v.Value, err)
					} else {
						defer conn.Close(ctx)
						log.Printf("[%d] (%s) OK %s \n", k, v.Note, v.Value)
					}

				}()

			}

		}
	}

}

func checkConnect(ctx context.Context, connection string) (err error) {

	var cfg *pgx.ConnConfig
	if cfg, err = pgx.ParseConfig(connection); err != nil {
		return err
	}
	_ = cfg

	conn, err := pgx.Connect(ctx, connection)
	if err != nil {
		err = fmt.Errorf("Unable to connect to database: %v", err)
		return
	}
	defer conn.Close(ctx)

	return
}
