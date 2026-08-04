// FILE IS AUTO GENERATED
// DO NOT EDIT MANUALLY
package app_context
import (
	"github.com/Shitcode-Swamp/go-train/src/core/someshi"
	"github.com/Shitcode-Swamp/go-train/src/core/someshi2"
)
type AppContext struct {
	testInjectee *someshi2.TestInjectee
	testInjectee2 *someshi2.TestInjectee
	test2 *someshi.TestStruct
	test *someshi.TestStruct
}

func NewAppContext() *AppContext {
	testInjectee := someshi2.Init()
	testInjectee2 := someshi2.Init2()
	test2 := someshi.Init2(testInjectee)
	test := someshi.Init(testInjectee, testInjectee2)
	return &AppContext{
		testInjectee: testInjectee,
		testInjectee2: testInjectee2,
		test2: test2,
		test: test,
	}
}
func (ctx *AppContext) Testinjectee() *someshi2.TestInjectee {
	return ctx.testInjectee
}
func (ctx *AppContext) Testinjectee2() *someshi2.TestInjectee {
	return ctx.testInjectee2
}
func (ctx *AppContext) Test2() *someshi.TestStruct {
	return ctx.test2
}
func (ctx *AppContext) Test() *someshi.TestStruct {
	return ctx.test
}
