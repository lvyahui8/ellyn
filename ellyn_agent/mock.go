package ellyn_agent

import (
	"encoding/json"
	"github.com/lvyahui8/ellyn/ellyn_common/collections"
	"reflect"
)

type MockRule struct {
	ruleId       int
	methodId     uint32
	returnValues []any
	// vars 根据请求或者表达式填充变量，在响应中可以使用
	vars []string
}

type Monkey struct {
	// RuleMap 规则配置
	RuleMap *collections.IntMapWrap[[]*MockRule]
}

func (m *Monkey) FilterMockRule(methodId uint32, methodArgs ...any) *MockRule {
	// 比对参数，看是否全部命中
	return nil
}

func (m *Monkey) BuildReturn(methodId uint32, rule *MockRule) (methodReturns []any) {
	// 按照mock配置构造返回值
	mtd := allMethods[methodId]
	list := mtd.returnTypeList
	methodReturns = make([]any, len(list))
	for i := 0; i < len(list); i++ {
		returnType := list[i]
		returnRawVal := rule.returnValues[i]
		switch n := returnRawVal.(type) {
		case string:
			if returnType.Kind() == reflect.Ptr && returnType.Elem().Kind() == reflect.Struct {
				v := reflect.New(returnType.Elem())
				_ = json.Unmarshal([]byte(n), v.Interface())
				methodReturns[i] = v.Interface()
			}
		case int:
			if returnType.Kind() == reflect.Int {
				methodReturns[i] = n
			}
		case int64:
			if returnType.Kind() == reflect.Int64 {
				methodReturns[i] = n
			}
		case float64, float32:
		case uintptr:
		}
	}
	return nil
}