package someshi

type TestStruct struct {
	TestField interface{}
}

// @Injection(test)
func Init(
	testField interface{} /* @Inject(testInjectee) */,
) *TestStruct {
	return &TestStruct{
		TestField: testField,
	}
}

// @Injection(test2)
func Init2(testField interface{} /* @Inject(testInjectee) */) *TestStruct {
	return &TestStruct{
		TestField: nil,
	}
}
