package entable_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/experiment/entable"
)

func TestTRun(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	ensure.Run("when table has no entries", func(ensure ensuring.E) {
		table := entable.T[Entry]{}

		calls := 0
		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			calls++
		})

		ensure(calls).Equals(0)
	})

	ensure.Run("when table has entries", func(ensure ensuring.E) {
		table := entable.T[Entry]{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		}

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/first",
			ensure.T().Name() + "/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})
}

func TestTRunWithIndex(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	ensure.Run("when table has no entries", func(ensure ensuring.E) {
		table := entable.T[Entry]{}

		calls := 0
		table.RunWithIndex(ensure, func(ensure ensuring.E, i int, entry *Entry) {
			calls++
		})

		ensure(calls).Equals(0)
	})

	ensure.Run("when table has entries", func(ensure ensuring.E) {
		table := entable.T[Entry]{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		}

		indices := []int{}
		names := []string{}
		calls := []*Entry{}

		table.RunWithIndex(ensure, func(ensure ensuring.E, i int, entry *Entry) {
			indices = append(indices, i)
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(indices).Equals([]int{0, 1})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/first",
			ensure.T().Name() + "/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})
}

func TestGRun(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	ensure.Run("when group has no entries", func(ensure ensuring.E) {
		group := entable.G[Entry]{}

		calls := 0
		group.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			calls++
		})

		ensure(calls).Equals(0)
	})

	ensure.Run("when group has entries", func(ensure ensuring.E) {
		group := entable.G[Entry]{
			Name: "grp",
			Table: entable.T[Entry]{
				{
					Name:  "first",
					Value: "hello",
				},
				{
					Name:  "second",
					Value: "world",
				},
			},
		}

		names := []string{}
		calls := []*Entry{}

		group.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/grp/first",
			ensure.T().Name() + "/grp/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})
}

func TestGRunWithIndex(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	ensure.Run("when group has no entries", func(ensure ensuring.E) {
		group := entable.G[Entry]{}

		calls := 0
		group.RunWithIndex(ensure, func(ensure ensuring.E, i int, entry *Entry) {
			calls++
		})

		ensure(calls).Equals(0)
	})

	ensure.Run("when group has entries", func(ensure ensuring.E) {
		group := entable.G[Entry]{
			Name: "grp",
			Table: entable.T[Entry]{
				{
					Name:  "first",
					Value: "hello",
				},
				{
					Name:  "second",
					Value: "world",
				},
			},
		}

		indices := []int{}
		names := []string{}
		calls := []*Entry{}

		group.RunWithIndex(ensure, func(ensure ensuring.E, i int, entry *Entry) {
			indices = append(indices, i)
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(indices).Equals([]int{0, 1})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/grp/first",
			ensure.T().Name() + "/grp/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})
}

func TestNewBuilder(t *testing.T) {
	ensure := ensure.New(t)

	//nolint:unused // It's used for a type parameter
	type Entry struct {
		Name  string
		Value string
	}

	ensure(entable.NewBuilder[Entry]()).Equals(&entable.Builder[Entry]{})
}

func TestBuilderRun(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	builder := entable.NewBuilder[Entry]()

	builder.Append(
		&Entry{
			Name:  "first entry",
			Value: "qwerty",
		},
	)

	builder.AppendTable(entable.T[Entry]{
		{
			Name:  "first in first table",
			Value: "hello",
		},
		{
			Name:  "second in first table",
			Value: "world",
		},
	})

	builder.AppendGroup(&entable.G[Entry]{
		Name: "grp1",
		Table: entable.T[Entry]{
			{
				Name:  "first in first group",
				Value: "1+2",
			},
			{
				Name:  "second in first group",
				Value: "3",
			},
		},
	})

	builder.Append(
		&Entry{
			Name:  "second entry",
			Value: "asdf",
		},
		&Entry{
			Name:  "third entry",
			Value: "zxcv",
		},
	)

	builder.AppendTable(entable.T[Entry]{
		{
			Name:  "first in second table",
			Value: "ice",
		},
		{
			Name:  "second in second table",
			Value: "cream",
		},
	})

	builder.AppendGroup(&entable.G[Entry]{
		Name: "grp2",
		Table: entable.T[Entry]{
			{
				Name:  "first in second group",
				Value: "bingo",
			},
			{
				Name:  "second in second group",
				Value: "time",
			},
		},
	})

	names := []string{}
	calls := []*Entry{}

	builder.Run(ensure, func(ensure ensuring.E, entry *Entry) {
		names = append(names, ensure.T().Name())
		calls = append(calls, entry)
	})

	ensure(names).Equals([]string{
		ensure.T().Name() + "/first_entry",
		ensure.T().Name() + "/first_in_first_table",
		ensure.T().Name() + "/second_in_first_table",
		ensure.T().Name() + "/grp1/first_in_first_group",
		ensure.T().Name() + "/grp1/second_in_first_group",
		ensure.T().Name() + "/second_entry",
		ensure.T().Name() + "/third_entry",
		ensure.T().Name() + "/first_in_second_table",
		ensure.T().Name() + "/second_in_second_table",
		ensure.T().Name() + "/grp2/first_in_second_group",
		ensure.T().Name() + "/grp2/second_in_second_group",
	})

	ensure(calls).Equals([]*Entry{
		{
			Name:  "first entry",
			Value: "qwerty",
		},
		{
			Name:  "first in first table",
			Value: "hello",
		},
		{
			Name:  "second in first table",
			Value: "world",
		},
		{
			Name:  "first in first group",
			Value: "1+2",
		},
		{
			Name:  "second in first group",
			Value: "3",
		},
		{
			Name:  "second entry",
			Value: "asdf",
		},
		{
			Name:  "third entry",
			Value: "zxcv",
		},
		{
			Name:  "first in second table",
			Value: "ice",
		},
		{
			Name:  "second in second table",
			Value: "cream",
		},
		{
			Name:  "first in second group",
			Value: "bingo",
		},
		{
			Name:  "second in second group",
			Value: "time",
		},
	})
}
