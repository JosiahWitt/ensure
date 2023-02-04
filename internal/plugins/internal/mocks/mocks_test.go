package mocks_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
)

func TestAllAddMock(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("adds mock", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m1 := m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m2 := m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		ensure(m1.Path).Equals("a")
		ensure(m1.Optional).Equals(true)
		ensure(m.Slice()[0].Path).Equals("a")
		ensure(m.Slice()[0].Optional).Equals(true)

		ensure(m2.Path).Equals("b")
		ensure(m2.Optional).Equals(false)
		ensure(m.Slice()[1].Path).Equals("b")
		ensure(m.Slice()[1].Optional).Equals(false)

		ensure(len(m.Slice())).Equals(2)
	})

	ensure.Run("when path is duplicated", func(ensure ensurepkg.Ensure) {
		defer func() {
			ensure(recover()).Equals(`mock with path "a" was already added: (PREVIOUS TYPE: string, NEW TYPE: bool)`)
		}()

		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(""))
		m.AddMock("b", true, reflect.TypeOf(1))
		m.AddMock("a", true, reflect.TypeOf(false))
	})
}

func TestAllSlice(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when no items", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		ensure(m.Slice()).IsEmpty()
	})

	ensure.Run("when many items", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		ensure(m.Slice()[0].Path).Equals("a")
		ensure(m.Slice()[0].Optional).IsTrue()

		ensure(m.Slice()[1].Path).Equals("b")
		ensure(m.Slice()[1].Optional).IsFalse()

		ensure(len(m.Slice())).Equals(2)
	})
}

func TestAllPathSet(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when no items", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		ensure(m.PathSet()).IsEmpty()
	})

	ensure.Run("when many items", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		ensure(m.PathSet()).Equals(mocks.PathSet{"a": struct{}{}, "b": struct{}{}})
	})
}

func TestMockImplements(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when interface is provided", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		type Greeter interface{ Hello(string) string }
		type Bingo interface{ Bingo(string) bool }

		greeter := reflect.TypeOf((*Greeter)(nil)).Elem()
		bingo := reflect.TypeOf((*Bingo)(nil)).Elem()

		ensure(m.Slice()[0].Implements(greeter)).IsTrue()
		ensure(m.Slice()[0].Implements(bingo)).IsFalse()
		ensure(m.Slice()[1].Implements(greeter)).IsFalse()
		ensure(m.Slice()[1].Implements(bingo)).IsTrue()
	})

	ensure.Run("when interface is not provided", func(ensure ensurepkg.Ensure) {
		defer func() {
			ensure(recover()).Equals("expected an interface to be provided to Implements, got: string")
		}()

		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.Slice()[0].Implements(reflect.TypeOf("not interface"))
	})
}

func TestMockSetValueByEntryIndex(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("sets values correctly", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		m.Slice()[0].SetValueByEntryIndex(5, reflect.ValueOf(&ExampleGreeterMock{"hey"}))
		m.Slice()[1].SetValueByEntryIndex(1, reflect.ValueOf(&ExampleOtherMock{"over"}))
		m.Slice()[0].SetValueByEntryIndex(0, reflect.ValueOf(&ExampleGreeterMock{"there"}))
		ensure(m.Slice()[0].ValueByEntryIndex(0).Interface()).Equals(&ExampleGreeterMock{"there"})
		ensure(m.Slice()[1].ValueByEntryIndex(1).Interface()).Equals(&ExampleOtherMock{"over"})
		ensure(m.Slice()[0].ValueByEntryIndex(5).Interface()).Equals(&ExampleGreeterMock{"hey"})
	})

	ensure.Run("when mock value does not match mock type", func(ensure ensurepkg.Ensure) {
		defer func() {
			ensure(recover()).Equals(`type of value for mock with path "a" was not the expected type: (EXPECTED: *mocks_test.ExampleGreeterMock, GOT: string)`)
		}()

		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		m.Slice()[0].SetValueByEntryIndex(5, reflect.ValueOf("hey"))
	})

	ensure.Run("when mock values are duplicated for an index", func(ensure ensurepkg.Ensure) {
		defer func() {
			ensure(recover()).Equals(`value at index 5 was already added for mock with path "a" and type: *mocks_test.ExampleGreeterMock`)
		}()

		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		m.Slice()[0].SetValueByEntryIndex(5, reflect.ValueOf(&ExampleGreeterMock{"hey"}))
		m.Slice()[0].SetValueByEntryIndex(5, reflect.ValueOf(&ExampleGreeterMock{"there"}))
	})
}

func TestMockValueByEntryIndex(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("gets values correctly", func(ensure ensurepkg.Ensure) {
		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		m.Slice()[0].SetValueByEntryIndex(5, reflect.ValueOf(&ExampleGreeterMock{"hey"}))
		m.Slice()[1].SetValueByEntryIndex(1, reflect.ValueOf(&ExampleOtherMock{"over"}))
		m.Slice()[0].SetValueByEntryIndex(0, reflect.ValueOf(&ExampleGreeterMock{"there"}))
		ensure(m.Slice()[0].ValueByEntryIndex(0).Interface()).Equals(&ExampleGreeterMock{"there"})
		ensure(m.Slice()[1].ValueByEntryIndex(1).Interface()).Equals(&ExampleOtherMock{"over"})
		ensure(m.Slice()[0].ValueByEntryIndex(5).Interface()).Equals(&ExampleGreeterMock{"hey"})
	})

	ensure.Run("when mock value is not set for index", func(ensure ensurepkg.Ensure) {
		defer func() {
			ensure(recover()).Equals(`value at index 5 was not set for mock with path "a" and type: *mocks_test.ExampleGreeterMock`)
		}()

		m := mocks.All{}
		m.AddMock("a", true, reflect.TypeOf(&ExampleGreeterMock{}))
		m.AddMock("b", false, reflect.TypeOf(&ExampleOtherMock{}))

		ensure(m.Slice()[0].ValueByEntryIndex(5).Interface()).Equals(&ExampleGreeterMock{"there"})
	})
}

func TestOnlyOneRequired(t *testing.T) {
	ensure := ensure.New(t)

	buildMocks := func(required ...bool) []*mocks.Mock {
		m := &mocks.All{}

		for i, r := range required {
			m.AddMock(strconv.Itoa(i), !r, nil)
		}

		return m.Slice()
	}

	ensure.Run("when none are required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(false, false, false)...)).IsFalse()
	})

	ensure.Run("when first is required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(true, false, false)...)).IsTrue()
	})

	ensure.Run("when middle is required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(false, true, false)...)).IsTrue()
	})

	ensure.Run("when last is required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(false, false, true)...)).IsTrue()
	})

	ensure.Run("when first and last are required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(true, false, true)...)).IsFalse()
	})

	ensure.Run("when first and middle are required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(true, true, false)...)).IsFalse()
	})

	ensure.Run("when middle and last are required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(false, true, true)...)).IsFalse()
	})

	ensure.Run("when all are required", func(ensure ensurepkg.Ensure) {
		ensure(mocks.OnlyOneRequired(buildMocks(true, true, true)...)).IsFalse()
	})
}

type ExampleGreeterMock struct{ unique string }

func (*ExampleGreeterMock) Hello(s string) string { return "hello, " + s }

type ExampleOtherMock struct{ unique string }

func (*ExampleOtherMock) Bingo(s string) bool { return false }
