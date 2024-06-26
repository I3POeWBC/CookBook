package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

const (
	CFG_CONNECTION        = "CONNECTION"
	HOST_INDEX     int    = 7
	PORT_INDEX     int    = 9
	DB_INDEX       int    = 11
	USER_INDEX     int    = 3
	PASSWORD_INDEX int    = 5
	DB_PASSWORD    string = "DB_PASSWORD"
)

type ConnDetail struct {
	DetailIndex int
	Diag        string
}

type Connection struct {
	Connection string `yaml:"connection"`
}

var (
	connDetails []ConnDetail = []ConnDetail{
		{USER_INDEX, "в строке подключения отсутствует логин"},
		//{PASSWORD_INDEX, "в строке подключения отсутствует пароль"},
		{HOST_INDEX, "в строке подключения отсутствует адрес сервера postgresql"},
		//{PORT_INDEX, "в строке подключения отсутствует порт, будет использован порт по-молчанию (5432)"},
		{DB_INDEX, "в строке подключения отсутствует имя базы"},
	}
)

func loadYaml(filePath string, a any) (err error) {

	var bs []byte
	if bs, err = os.ReadFile(filePath); err != nil {
		return
	} else if err = yaml.Unmarshal(bs, a); err != nil {
		return
	}

	return
}

func parseConnectionString(a string) (ret string, err error) {

	reStr := `^(postgres://)(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/](.+))?$`
	re := regexp.MustCompile(reStr)

	if !re.MatchString(a) {
		err = fmt.Errorf("строка [%s] не соответствует формату", a)
		return
	}

	ss := re.FindStringSubmatch(a)
	_ = ss

	for _, v := range connDetails {
		if ss[v.DetailIndex] == "" {
			err = fmt.Errorf("%v", v.Diag)
			return
		}
	}

	if ss[PASSWORD_INDEX] == "" {
		log.Printf("в строке подключения отсутствует пароль. \n")

		if envPass := os.Getenv(DB_PASSWORD); envPass == "" {
			err = fmt.Errorf("в окружении отсутствует/пустая переменная среды [%s] для пароля к базе данных", DB_PASSWORD)
			return

			/* Не работает, по-краней мере, в windows
			if pass, errPass := term.ReadPassword(1); errPass != nil || string(pass) == "" {
				err = fmt.Errorf("с клавиатуры введен пустой пароль или произошли ошибки err: [%v]", errPass)
				return
			} else {
				ss[PASSWORD_INDEX] = string(pass)
			}
			*/
		} else {
			ss[PASSWORD_INDEX] = envPass
		}
	}

	if ss[PORT_INDEX] == "" {
		log.Printf("в строке подключения отсутствует порт, будет использован порт по-молчанию (5432)")
		ss[PORT_INDEX] = "5432"
	}

	var port int
	if port, err = strconv.Atoi(ss[PORT_INDEX]); err != nil {
		err = fmt.Errorf("ошибочное значение порта [%v]", err)
		return
	}

	//fmt.Printf("src:[%s]\n", a)
	//for k, v := range ss {
	//	fmt.Printf("\t[%d]: [%s]\n", k, v)
	//}

	// `^(postgres://)(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/](.+))?$`
	ret = ss[1] + ss[3] + ":" + ss[5] + "@" + ss[7] + fmt.Sprintf(":%d", port) + ss[10]
	//fmt.Printf("%s\n\n", ret)

	return
}
