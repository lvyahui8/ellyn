package agent

import (
	"encoding/json"
	"github.com/lvyahui8/ellyn/sdk/common/collections"
	"reflect"
)

type MockRule struct {
	ruleId       int
	methodId     uint32
	returnValues []any
	args         []any
	// vars 根据请求或者表达式填充变量，在响应中可以使用
	vars []string
}

type Monkey struct {
	// RuleMap 规则配置
	RuleMap *collections.IntMapWrap[[]*MockRule]
}

func (m *Monkey) FilterMockRule(methodId uint32, methodArgs ...any) *MockRule {
	// 比对参数，看是否全部命中
	mtd := methods[methodId]
	rules, exist := m.RuleMap.Get(int(methodId))
	if !exist {
		return nil
	}
	for _, rule := range rules {
		for i := 0; i < mtd.ArgsList.Count(); i++ {
			argType := mtd.ArgsList.Type(i)
			actualVal := methodArgs[i]
			if reflect.ValueOf(actualVal).Type().String() != argType {
				continue
			}
			filterVal := rule.args[i]
			if filterVal == nil {
				continue
			}
			if reflect.DeepEqual(actualVal, filterVal) {
				return rule
			}
		}
	}

	return nil
}

func (m *Monkey) BuildReturn(methodId uint32, rule *MockRule, rawReturns []any) (methodReturns []any) {
	// 按照mock配置构造返回值
	mtd := methods[methodId]
	methodReturns = make([]any, mtd.ReturnList.Count())
	for i := 0; i < mtd.ReturnList.Count(); i++ {
		returnType := mtd.ReturnList.Type(i)
		returnRawVal := rule.returnValues[i]
		switch n := returnRawVal.(type) {
		case string:
			if returnType == reflect.Ptr.String() && returnType == "*"+reflect.Struct.String() {
				v := reflect.New(reflect.TypeOf(rawReturns[i]))
				_ = json.Unmarshal([]byte(n), v.Interface())
				methodReturns[i] = v.Interface()
			}
		case int:
			if returnType == reflect.Int.String() {
				methodReturns[i] = n
			}
		case int64:
			if returnType == reflect.Int64.String() {
				methodReturns[i] = n
			}
		case float64, float32:
		case uintptr:
		}
	}
	return nil
}
