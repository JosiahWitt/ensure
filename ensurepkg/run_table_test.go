package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/tests/mocks/mock_ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRunTableByIndex(t *testing.T) {
	table := []struct {
		Name                 string
		ExpectedHelperCalls  int
		ExpectedNames        []string
		ExpectedPanicMessage string
		Table                interface{}
	}{
		{
			Name:                "with valid table: slice",
			ExpectedHelperCalls: 4,
			ExpectedNames:       []string{"name 1", "name 2"},
			Table: []struct {
				Name  string
				Value string
			}{
				{
					Name:  "name 1",
					Value: "item 1",
				},
				{
					Name:  "name 2",
					Value: "item 2",
				},
			},
		},

		{
			Name:                "with valid table: array",
			ExpectedHelperCalls: 4,
			ExpectedNames:       []string{"name 1", "name 2"},
			Table: [2]struct {
				Name  string
				Value string
			}{
				{
					Name:  "name 1",
					Value: "item 1",
				},
				{
					Name:  "name 2",
					Value: "item 2",
				},
			},
		},

		{
			Name:                 "with invalid table type: not array or slice",
			ExpectedPanicMessage: "Expected a slice or array for the table, got string",
			Table:                "my table",
		},

		{
			Name:                 "with invalid table type: not array or slice of stucts",
			ExpectedPanicMessage: "Expected entry in table to be a struct, got string",
			Table: []string{
				"item 1",
				"item 2",
			},
		},

		{
			Name:                 "with missing name",
			ExpectedPanicMessage: "Name field does not exist on struct in table",
			Table: []struct {
				Value string
			}{
				{
					Value: "item 1",
				},
				{
					Value: "item 2",
				},
			},
		},

		{
			Name:                 "with name with invalid type",
			ExpectedPanicMessage: "Name field in struct in table is not a string",
			Table: []struct {
				Name  int
				Value string
			}{
				{
					Name:  1,
					Value: "item 1",
				},
				{
					Name:  2,
					Value: "item 2",
				},
			},
		},

		{
			Name:                 "with missing name for one item",
			ExpectedPanicMessage: "Errors encountered while building table:\n - table[1]: Name not set for item",
			Table: []struct {
				Name  string
				Value string
			}{
				{
					Name:  "name 1",
					Value: "item 1",
				},
				{
					Name:  "",
					Value: "item 2",
				},
			},
		},

		{
			Name:                 "with duplicate name",
			ExpectedPanicMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1",
			Table: []struct {
				Name  string
				Value string
			}{
				{
					Name:  "name 1",
					Value: "item 1",
				},
				{
					Name:  "name 2",
					Value: "item 2",
				},
				{
					Name:  "name 1",
					Value: "item 3",
				},
			},
		},

		{
			Name:                 "with double duplicate name",
			ExpectedPanicMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1\n - table[3]: duplicate Name found; first occurrence was table[0].Name: name 1",
			Table: []struct {
				Name  string
				Value string
			}{
				{
					Name:  "name 1",
					Value: "item 1",
				},
				{
					Name:  "name 2",
					Value: "item 2",
				},
				{
					Name:  "name 1",
					Value: "item 3",
				},
				{
					Name:  "name 1",
					Value: "item 4",
				},
			},
		},
	}

	for _, entry := range table {
		entry := entry // Pin range variable

		t.Run(entry.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockT := mock_ensurepkg.NewMockT(ctrl)
			mockT.EXPECT().Helper().Times(entry.ExpectedHelperCalls)

			// Build expected Run calls only if there's no expected error
			expectedTestingInputs := []*testing.T{}
			if entry.ExpectedPanicMessage == "" {
				expectedRunCalls := []*gomock.Call{}
				for _, name := range entry.ExpectedNames {
					providedTestingInput := &testing.T{}
					expectedTestingInputs = append(expectedTestingInputs, providedTestingInput)

					expectedRunCalls = append(expectedRunCalls,
						mockT.EXPECT().Run(name, gomock.Any()).
							Do(func(name string, fn func(t *testing.T)) {
								fn(providedTestingInput)
							}),
					)
				}

				// Run calls should be in order
				gomock.InOrder(expectedRunCalls...)
			} else {
				// Setup panic recovery
				defer func() {
					if msg := recover(); msg.(string) != entry.ExpectedPanicMessage {
						t.Errorf("Expected panic message '%s', got: %v", entry.ExpectedPanicMessage, msg)
					}
				}()
			}

			type Params struct {
				ensure ensurepkg.Ensure
				i      int
			}

			// Run table and save parameters
			actualParams := []Params{}
			ensure := ensure.New(mockT)
			ensure.RunTableByIndex(entry.Table, func(ensure ensurepkg.Ensure, i int) {
				actualParams = append(actualParams, Params{ensure: ensure, i: i})
			})

			// This should not be reached on panic
			if entry.ExpectedPanicMessage != "" {
				t.Errorf("Expected panic, got none")
			}

			// Verify call count
			if len(actualParams) != len(entry.ExpectedNames) {
				t.Errorf("len(actualParams) != len(entry.ExpectedNames)")
			}

			// Verify parameters are correct
			for i, actualParam := range actualParams {
				if actualParam.ensure.T() != expectedTestingInputs[i] {
					t.Errorf("actualParams[%d].ensure.T() != expectedTestingInputs[%d]", i, i)
				}

				if actualParam.i != i {
					t.Errorf("actualParams[%d].i != %d", i, i)
				}
			}
		})
	}
}
