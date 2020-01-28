package postgres

import (
	"context"
	"github.com/balerter/balerter/internal/datasource/converter"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"time"
)

func (m *Postgres) query(L *lua.LState) int {

	q := L.Get(1).String()

	m.logger.Debug("call postgres query", zap.String("query", q))

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3) // todo: timeout to settings
	defer ctxCancel()

	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		m.logger.Error("error postgres query", zap.String("query", q), zap.Error(err))
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()

	cct, _ := rows.ColumnTypes()

	dest := make([]interface{}, 0)
	ffs := make([]func(v interface{}) string, 0)

	for _, c := range cct {
		switch c.DatabaseTypeName() {
		case "NUMERIC":
			dest = append(dest, new(float64))
			ffs = append(ffs, converter.FromFloat64)
		case "TIMESTAMPTZ":
			dest = append(dest, new(time.Time))
			ffs = append(ffs, converter.FromDateTime)
		case "VARCHAR":
			dest = append(dest, new(string))
			ffs = append(ffs, converter.FromString)
		case "INT4":
			dest = append(dest, new(int))
			ffs = append(ffs, converter.FromInt)
		case "BOOL":
			dest = append(dest, new(bool))
			ffs = append(ffs, converter.FromBoolean)
		default:
			m.logger.Error("error scan type", zap.String("typename", c.DatabaseTypeName()))
			L.Push(lua.LNil)
			L.Push(lua.LString("error database type"))
			return 2
		}
	}

	result := &lua.LTable{}

	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			m.logger.Error("error scan", zap.Error(err))
			L.Push(lua.LNil)
			L.Push(lua.LString("error scan: " + err.Error()))
			return 2
		}

		row := &lua.LTable{}

		for idx, c := range cct {
			v := ffs[idx](dest[idx])
			row.RawSet(lua.LString(c.Name()), lua.LString(v))
		}

		result.Append(row)
	}

	L.Push(result)
	L.Push(lua.LNil)
	return 2
}
