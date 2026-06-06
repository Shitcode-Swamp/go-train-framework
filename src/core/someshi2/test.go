package someshi2

type TestInjectee struct {

}

// @Injection(testInjectee)
func Init() *TestInjectee {
	return &TestInjectee{

	}
} 
