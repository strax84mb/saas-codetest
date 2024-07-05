package tests_test

import (
	"testing"

	tests "github.com/SnowSoftwareGlobal/saas-codetests"
	"github.com/stretchr/testify/assert"
)

type letterTestCase struct {
	input    string
	expected rune
}

func TestGetDeepestLetter(t *testing.T) {
	testCases := []letterTestCase{
		{input: "a(b)c", expected: 'b'},
		// expected was "G" but that is not lowercase so I changed it to "f"
		{input: "((a))(((M)))(c)(D)(e)(((f))(((G))))h(i)", expected: 'f'},
		// expected was "c" but it should be "?" as there is a closing parenthesis missing
		// this is clearly a malformed input
		{input: "((A)(b)c", expected: '?'},
		// expected was "g" but that is not lowercase
		// "a" and "c" are on the same level but "a" was encountered first
		// therefore I will expect "a"
		{input: "(a)((G)c)", expected: 'a'},
		{input: "(8)", expected: '?'},
		{input: "(!)", expected: '?'},
		// added this one as the test is missing this malformed check
		// without this I could simply count number of opened and closed parenthesis
		// and make sure they are equal which is clearly a faulty algorythm
		{input: "((q)w))(", expected: '?'},
	}

	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			actual := tests.GetDeepestLetter(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
