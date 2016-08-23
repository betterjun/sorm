package sorm

type tbs struct {
	SId   int    `orm:"pk=1;fn=id"`
	Name  string `orm:"_"`
	Dummy string `orm:"fn=dummy"`
}

const (
	CONN_STRING = "root:root@tcp(127.0.0.1:3306)/world"
)
