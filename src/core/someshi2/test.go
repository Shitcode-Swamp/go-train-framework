package someshi2

type TestInjectee struct {

}

// @Injection(testInjectee)
func Init() *TestInjectee {
	return &TestInjectee{

	}
} 

// @Injection(testInjectee2)
func Init2() *TestInjectee {
	return &TestInjectee{

	}
}
