package app

import "testing"

func TestInfoSchema(t *testing.T) {
	is := &InfoSchema{}
	wnt := "information_schema.TABLES"
	got := is.TableName()
	if got != wnt {
		t.Errorf("Wanted %s, got: %s", wnt, got)
	}
}

func TestTableSchema(t *testing.T) {
	tb := &TableBench{}
	wnt := "temp_bench_table"
	got := tb.TableName()
	if got != wnt {
		t.Errorf("Wanted `%s`,  got: %s", wnt, got)
	}
}
