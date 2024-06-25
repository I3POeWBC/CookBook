package main

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

const (
	CFG_CONNECTION = "CONNECTION"
)

type Connection struct {
	Connection string `yaml:"connection"`
}

func loadYaml(filePath string, a any) (err error) {

	var bs []byte
	if bs, err = os.ReadFile(filePath); err != nil {
		return
	} else if err = yaml.Unmarshal(bs, a); err != nil {
		return
	}

	return
}

/*
func parseConnectionString(a string, connectTimeout time.Duration) (ret *pgx.ConnConfig, err error) {

	reStr := `^postgres://(([^:@]+)(:([^@]+))?[@])(([^:]+)([:]([0-9]+))?)([/](.+))?$`
	re := regexp.MustCompile(reStr)

	if !re.MatchString(a) {
		err = fmt.Errorf("строка [%s] не соотвествует формату", a)
		return
	}

	ss := re.FindStringSubmatch(a)
	_ = ss

	HOST_INDEX := 6
	PORT_INDEX := 8
	DB_INDEX := 10
	USER_INDEX := 2
	PASSWORD_INDEX := 4

	port, err := strconv.Atoi(ss[PORT_INDEX])
	if err != nil {
		err = fmt.Errorf("ошибочное значение порта [%v]", err)
		return
	}

	cfg := &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:           ss[HOST_INDEX],
			Port:           uint16(port),
			Database:       ss[DB_INDEX],
			User:           ss[USER_INDEX],
			Password:       ss[PASSWORD_INDEX],
			TLSConfig:      nil,
			ConnectTimeout: connectTimeout,
		},
		Tracer:                   nil,
		StatementCacheCapacity:   0,
		DescriptionCacheCapacity: 0,
		DefaultQueryExecMode:     0,
	}

	ret = cfg

	return
}
*/
