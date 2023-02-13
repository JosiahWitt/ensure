package entable_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/experiments/entable"
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
