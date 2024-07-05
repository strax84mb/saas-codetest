// I changed this package to tests_test.
// The logic behind this decision is to block test dependencies from relase build.
// While it is true that files ending with "_test.go" don't have their code compiled
// their imports are still taken into account so all the regular go files from those
// dependencies do end up as part of the build.
//
// Now, what I wrote here may not make sense as go compiler has dead code detection
// to identify unreachable code as part of build process. However, we should keep
// in mind that dead code detection works on package level, NOT on function level
// (like you would see in Rust). So, while test code isn't reached in any way from
// main function it is still reachable and that right there is a key word here.
//
// Now, if we absolutelly must have tests in the same package as the code being tested
// I suggest using build tags. In test file have this line
//
//	//go:build test
//
// and then use "test" tag when building
//
//	go build -tags test .
//
// or if you are using Visual Studio Code add this line to JSON in ".vscode/settings.json"
//
//	"go.buildTags": "test"
//
// in order to be able to run test from those files
//
// Both approaches I stated are valid and work well at making sure that dev depepndencies
// never end up as part of the release binary, keeping our binary nice, neat and light weight.
package tests_test

import (
	"fmt"
	"testing"

	tests "github.com/SnowSoftwareGlobal/saas-codetests"
	"github.com/stretchr/testify/assert"
)

// moved this here as it shouldn't be in code file
// this is purely a test structure
type apiTestCase struct {
	users  []*tests.User
	input  tests.UpdateUserRequest
	output *tests.User
	err    error
}

// moved this here as it shouldn't be in code file
// this is function is only used in test
func pointer(v string) *string {
	return &v
}

func TestEndpoint(t *testing.T) {
	testCases := []apiTestCase{
		// removed &User in users field as it is unnecessary
		{
			users:  []*tests.User{{Id: "6a43df", FullName: "Tom Jefferson", Email: "jefferson999@mirro.com"}},
			input:  tests.UpdateUserRequest{Id: "6a43df", Email: pointer("t.jefferson@mirro.com")},
			output: &tests.User{Id: "6a43df", FullName: "Tom Jefferson", Email: "t.jefferson@mirro.com"},
		},
		{
			users: []*tests.User{{Id: "56781a", FullName: "Eric Nilsson", Email: "eric_fantastic@offtop.com"}},
			input: tests.UpdateUserRequest{Id: "56781c", FullName: pointer("Eric Fantastic")},
			err:   tests.UserNotFound,
		},
		{
			users: []*tests.User{{Id: "556f36", FullName: "Antony Downtown", Email: "antony.downtown@gmail.com"}},
			// Id was 556f37 but should be 556f36
			input:  tests.UpdateUserRequest{Id: "556f36", FullName: pointer("Antony Uptown")},
			output: &tests.User{Id: "556f36", FullName: "Antony Uptown", Email: "antony.downtown@gmail.com"},
		},
		{
			users:  []*tests.User{{Id: "34d35", FullName: "Mickle Now", Email: "m.n@story.com"}},
			input:  tests.UpdateUserRequest{Id: "34d35"},
			output: &tests.User{Id: "34d35", FullName: "Mickle Now", Email: "m.n@story.com"},
		},
		{
			users: []*tests.User{},
			input: tests.UpdateUserRequest{Id: "34d35", FullName: pointer("Nina Mitk"), Email: pointer("m.n@story.com")},
			err:   tests.UserNotFound,
		},
	}
	for ind, test := range testCases {
		t.Run(fmt.Sprint(ind), func(t *testing.T) {
			api, err := tests.NewUserApi(test.users)
			// added testify/assert here as I find it a much less wordy way to write assertions
			// while at the same time being clear and readable
			assert.NoError(t, err)
			res, err := api.Update(test.input)
			assert.Truef(t, assert.ObjectsAreEqual(res, test.output), "actual result is not as expected")
			assert.ErrorIsf(t, err, test.err, "received error is not as expected")
		})
	}
}
