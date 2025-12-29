/*
Lua modding. A way to execute complex code from a string/file, as well as that string/file
executing any provided functions. Some very basic Lua modules are included - everything else
is excluded for security reasons, especially file manipulation.
*/
package script

import (
	"pure-game-kit/debug"
	"pure-game-kit/utility/text"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type Script struct{ state *lua.LState }

func New() *Script {
	var state = lua.NewState(lua.Options{SkipOpenLibs: true})
	lua.OpenBase(state)
	lua.OpenMath(state)
	lua.OpenString(state)
	lua.OpenTable(state)
	return &Script{state: state}
}

func (s *Script) AddFunction(functionName string, function any) {
	var rv = reflect.ValueOf(function)
	var rt = rv.Type()

	if rt.Kind() != reflect.Func {
		return
	}

	var luaFn = func(L *lua.LState) int {
		var numArgs = L.GetTop()
		if numArgs != rt.NumIn() {
			return 0
		}

		var in = make([]reflect.Value, rt.NumIn())
		for i := 0; i < rt.NumIn(); i++ {
			var arg = L.Get(i + 1)
			var result reflect.Value
			var t reflect.Type = rt.In(i)
			switch t.Kind() {
			case reflect.String:
				result = reflect.ValueOf(arg.String()).Convert(t)
			case reflect.Bool:
				result = reflect.ValueOf(lua.LVAsBool(arg)).Convert(t)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result = reflect.ValueOf(int64(lua.LVAsNumber(arg))).Convert(t)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				result = reflect.ValueOf(uint64(lua.LVAsNumber(arg))).Convert(t)
			case reflect.Float32, reflect.Float64:
				result = reflect.ValueOf(float64(lua.LVAsNumber(arg))).Convert(t)
			default:
				result = reflect.Zero(t)
			}

			in[i] = result
		}

		var out = rv.Call(in)
		for _, o := range out {
			L.Push(valueToLuaType(L, o.Interface()))
		}
		return len(out)
	}

	s.state.SetGlobal(functionName, s.state.NewFunction(luaFn))
}
func (s *Script) ExecuteCode(code string) bool {
	var err = s.state.DoString(code)
	if err != nil {
		debug.LogError("Failed to execute code!\n", err)
	}
	return err == nil
}
func (s *Script) ExecuteFunction(functionName string, parameters ...any) any {
	var fn = s.state.GetGlobal(functionName)
	if fn.Type() != lua.LTFunction {
		debug.LogError("Failed to find function: \"", functionName, "\"")
		return nil
	}

	var luaArgs = make([]lua.LValue, len(parameters))
	for i, arg := range parameters {
		luaArgs[i] = valueToLuaType(s.state, arg)
	}

	var err = s.state.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}, luaArgs...)
	if err != nil {
		debug.LogError("Failed to call function: \"", functionName, "\"\n", err)
		return nil
	}

	var ret = s.state.Get(-1)
	s.state.Pop(1)
	return luaToGoValue(ret)
}

func (s *Script) Close() {
	s.state.Close()
}

// =================================================================
// private

func valueToLuaType(L *lua.LState, val any) lua.LValue {
	switch v := val.(type) {
	case nil:
		return lua.LNil
	case string:
		return lua.LString(v)
	case bool:
		return lua.LBool(v)
	case int, int8, int16, int32, int64:
		return lua.LNumber(reflect.ValueOf(v).Convert(reflect.TypeOf(int64(0))).Int())
	case uint, uint8, uint16, uint32, uint64:
		return lua.LNumber(reflect.ValueOf(v).Convert(reflect.TypeOf(uint64(0))).Uint())
	case float32, float64:
		return lua.LNumber(reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0))).Float())
	case []any:
		var tbl = L.NewTable()
		for i, e := range v {
			tbl.RawSetInt(i+1, valueToLuaType(L, e))
		}
		return tbl
	case map[any]any:
		var tbl = L.NewTable()
		for key, val := range v {
			tbl.RawSetString(text.New(key), valueToLuaType(L, val))
		}
		return tbl
	default:
		return lua.LString(text.New(v))
	}
}
func luaToGoValue(val lua.LValue) any {
	switch v := val.(type) {
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		var arr = []any{}
		var m = map[any]any{}
		var isArray = true
		var i = 1
		v.ForEach(func(key, value lua.LValue) {
			var goKey = luaToGoValue(key)
			var goValue = luaToGoValue(value)

			if key.Type() == lua.LTNumber && int(lua.LVAsNumber(key)) == i {
				arr = append(arr, goValue)
				i++
			} else {
				isArray = false
				m[goKey] = goValue
			}
		})

		if isArray {
			return arr
		}
		return m
	case *lua.LFunction:
		return v // or nil, depending on what you want
	case *lua.LUserData:
		return v.Value
	default:
		return nil
	}
}
