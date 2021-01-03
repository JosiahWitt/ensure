package ensurepkg_test

import "testing"

func TestEnsureRun(t *testing.T) {}

// func Test(t *testing.T) {
// 	ensure := ensure.New(t)

// 	ensure(map[string]interface{}{"a": "B", "c": 123, "d": "abc"}).Equals(map[string]interface{}{"a": "B", "c": 234, "d": "abb", "x": "sss"})
// 	ensure.Run("abc", func(ensure ensurepkg.Ensure) {
// 		ensure("abc").Equals("abc")
// 	})

// 	err := errors.New("abc")
// 	ensure(err).IsError(err)
// 	ensure(true).IsTrue()

// 	table := []struct {
// 		Name     string
// 		Actual   string
// 		Expected string
// 	}{
// 		{
// 			Name:     "thing 1",
// 			Actual:   "abc",
// 			Expected: "abc",
// 		},
// 		{
// 			Name:     "thing 2",
// 			Actual:   "abc",
// 			Expected: "abcx",
// 		},
// 	}

// 	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
// 		entry := table[i]

// 		ensure(entry.Actual).Equals(entry.Expected)
// 	})
// }
