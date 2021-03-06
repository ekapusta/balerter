package chart

import (
	"bytes"
	"github.com/balerter/balerter/internal/script/script"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

func ModuleName() string {
	return "chart"
}

func Methods() []string {
	return []string{
		"render",
	}
}

type DataItem struct {
	Timestamp float64
	Value     float64
}

type DataSeries struct {
	Color      string
	LineColor  string
	PointColor string
	Data       []DataItem
}

type Data struct {
	Title  string
	Series []DataSeries
}

type Chart struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Chart {
	l := &Chart{
		logger: logger,
	}

	return l
}

func (ch *Chart) Name() string {
	return ModuleName()
}

func (ch *Chart) Stop() error {
	return nil
}

func (ch *Chart) GetLoader(s *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		var exports = map[string]lua.LGFunction{
			"render": ch.render(s),
		}

		mod := luaState.SetFuncs(luaState.NewTable(), exports)

		luaState.Push(mod)
		return 1 //nolint:mnd
	}
}

func (ch *Chart) render(_ *script.Script) lua.LGFunction {
	return func(luaState *lua.LState) int {
		ch.logger.Debug("Chart.Render")

		chartTitle := luaState.Get(1) //nolint:mnd
		if chartTitle.Type() == lua.LTNil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("title must be defined"))
			return 2 //nolint:mnd
		}

		chartData := luaState.Get(2) //nolint:mnd
		if chartData.Type() != lua.LTTable {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("chart data table must be defined"))
			return 2 //nolint:mnd
		}

		data := &Data{}

		err := gluamapper.Map(chartData.(*lua.LTable), data)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("wrong chart data format, " + err.Error()))
			return 2 //nolint:mnd
		}

		buf := bytes.NewBuffer([]byte{})

		err = ch.Render(chartTitle.String(), data, buf)
		if err != nil {
			luaState.Push(lua.LNil)
			luaState.Push(lua.LString("error render chart, " + err.Error()))
			return 2 //nolint:mnd
		}

		luaState.Push(lua.LString(buf.String()))

		return 1 //nolint:mnd
	}
}
