package integration_test_suite_test

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"go.uber.org/mock/gomock"
)

func TestGoMockController(t *testing.T) {
	ensure := ensure.New(t)

	assertEq(t, ensure.GoMockController(), gomock.NewController(t))
}

func TestInterfaceT(t *testing.T) {
	ensure := ensure.New(t)

	assertEq(t, ensure.InterfaceT(), t)
}

func TestRun(t *testing.T) {
	sharedEnsureRunTests(t, func(ensure ensuring.E) func(string, func(ensuring.E)) {
		return ensure.Run
	})
}

func sharedEnsureRunTests(t *testing.T, prepare func(ensure ensuring.E) func(string, func(ensuring.E))) {
	t.Run("callback is executed in nested scope", func(t *testing.T) {
		namePrefix := t.Name()
		ensure := ensure.New(t)

		name := ""
		run := prepare(ensure)
		run("some name", func(ensure ensuring.E) {
			name = ensure.T().Name()
		})

		// Shows that ensure.Run executed with a nested scope
		assertEq(t, name, namePrefix+"/some_name")
	})
}

func TestRunParallel(t *testing.T) {
	ensure := ensure.New(t)

	var mu sync.Mutex
	callOrder := []int{}

	for i := range 1000 {
		ensure.RunParallel(fmt.Sprintf("parallel %d", i), func(ensure ensuring.E) {
			mu.Lock()
			defer mu.Unlock()
			callOrder = append(callOrder, i)
		})
	}

	t.Cleanup(func() {
		// The execution order won't be sorted when it is called in parallel
		assertEq(t, sort.IsSorted(sort.IntSlice(callOrder)), false)
	})
}

func TestRunTableByIndex(t *testing.T) {
	sharedEnsureRunTableByIndexTests(t, func(ensure ensuring.E) func(table any, fn func(ensure ensuring.E, i int)) {
		return ensure.RunTableByIndex
	})
}

func sharedEnsureRunTableByIndexTests(t *testing.T, prepare func(ensure ensuring.E) func(table any, fn func(ensure ensuring.E, i int))) {
	t.Run("entries are executed in nested scope when entries are not pointers", func(t *testing.T) {
		namePrefix := t.Name()
		ensure := ensure.New(t)

		table := []struct {
			Name string
		}{
			{
				Name: "first one",
			},
			{
				Name: "second one",
			},
		}

		names := []string{}
		fullNames := []string{}
		runTableByIndex := prepare(ensure)
		runTableByIndex(table, func(ensure ensuring.E, i int) {
			entry := table[i]

			names = append(names, entry.Name)
			fullNames = append(fullNames, ensure.T().Name())
		})

		assertEq(t, names, []string{"first one", "second one"})

		// Shows that ensure.RunTableByIndex executed each entry with a nested scope
		assertEq(t, fullNames, []string{namePrefix + "/first_one", namePrefix + "/second_one"})
	})

	t.Run("entries are executed in nested scope when entries are pointers", func(t *testing.T) {
		namePrefix := t.Name()
		ensure := ensure.New(t)

		table := []*struct {
			Name string
		}{
			{
				Name: "first one",
			},
			{
				Name: "second one",
			},
		}

		names := []string{}
		fullNames := []string{}
		ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
			entry := table[i]

			names = append(names, entry.Name)
			fullNames = append(fullNames, ensure.T().Name())
		})

		assertEq(t, names, []string{"first one", "second one"})

		// Shows that ensure.RunTableByIndex executed each entry with a nested scope
		assertEq(t, fullNames, []string{namePrefix + "/first_one", namePrefix + "/second_one"})
	})

	t.Run("hydrates mocks and injects subject", func(t *testing.T) {
		ensure := ensure.New(t)

		type MySubject struct {
			SomeT testctx.T
		}

		type Mocks struct {
			T *mock_testctx.MockT
		}

		loggedMessages := []string{}

		table := []struct {
			Name string

			Location string

			Mocks      *Mocks
			SetupMocks func(*Mocks)
			Subject    *MySubject
		}{
			{
				Name: "about the world",

				Location: "world",

				SetupMocks: func(m *Mocks) {
					m.T.EXPECT().Logf("hello world").Do(func(format string, args ...any) {
						loggedMessages = append(loggedMessages, format)
					})
				},
			},
			{
				Name: "about the universe",

				Location: "universe",

				SetupMocks: func(m *Mocks) {
					m.T.EXPECT().Logf("hello universe").Do(func(format string, args ...any) {
						loggedMessages = append(loggedMessages, format)
					})
				},
			},
		}

		ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
			entry := table[i]

			entry.Subject.SomeT.Logf("hello " + entry.Location)
		})

		assertEq(t, loggedMessages, []string{"hello world", "hello universe"})
	})
}

func TestT(t *testing.T) {
	ensure := ensure.New(t)

	assertEq(t, ensure.T(), t)
}

func assertEq(t *testing.T, actual, expected any) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%+v != %+v", actual, expected)
	}
}
