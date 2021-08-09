package sqlgen

import (
	"fmt"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
)

type SqlGen struct {
	table string
	dg    *datagen.DataGen
}

func NewSqlGen(table string, dg *datagen.DataGen) *SqlGen {
	if table == "" {
		panic("Table cannot be empty")
	}
	if dg == nil {
		panic("DataGen cannot be nil")
	}
	return &SqlGen{
		table: table,
		dg:    dg,
	}
}

func (sqlg *SqlGen) InsertStatement(n int) string {
	rec := sqlg.dg.Record(n)
	stmt := fmt.Sprintf(`INSERT INTO %s (pk, uniq, small_grp, large_grp, fixed_val, seq_num, payload) 
		                        VALUES ('%s', '%s', '%s', '%s', '%s', %d, '%s')`,
		sqlg.table,
		rec.Pk, rec.Uniq, rec.SmallGrp, rec.LargeGrp, rec.FixedValue, rec.SeqNum, rec.Payload)
	return stmt
}

func (sqlg *SqlGen) SelectByPkStatement(n int) string {
	stmt := fmt.Sprintf("SELECT pk, uniq, small_grp, large_grp, fixed_val, seq_num, ts, payload FROM %s where pk='%s'",
		sqlg.table, sqlg.dg.Key(n))
	return stmt
}

func (sqlg *SqlGen) SelectByUkStatement(n int) string {
	stmt := fmt.Sprintf("SELECT pk, uniq, small_grp, large_grp, fixed_val, seq_num, ts, payload FROM %s where uniq='%s'",
		sqlg.table, sqlg.dg.Uniq(n))
	return stmt
}

func (sqlg *SqlGen) UpdatePayloadByPkStatement(n int) string {
	stmt := fmt.Sprintf("UPDATE %s SET payload='%s' where pk='%s'",
		sqlg.table, sqlg.dg.RandomPayload(), sqlg.dg.Key(n))
	return stmt
}

func (sqlg *SqlGen) DeleteByPkStatement(n int) string {
	stmt := fmt.Sprintf("DELETE FROM %s where pk='%s'",
		sqlg.table, sqlg.dg.Key(n))
	return stmt
}

