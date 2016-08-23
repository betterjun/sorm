package sorm

type tbs struct {
	SId   int    `sorm:"fn=id"`
	Name  string `sorm:"_"`
	Dummy string `sorm:"fn=dummy"`
}

const (
	CONN_STRING = "root:root@tcp(127.0.0.1:3306)/world"
)
