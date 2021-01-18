package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRunTableByIndex(t *testing.T) {
	table := []struct {
		Name                 string
		ExpectedNames        []string
		ExpectedFatalMessage string
		Table                interface{}
	}{
		{
			Name:          "with valid table: slice",
			ExpectedNames: []string{"name 1", "name 2", "name 3"},
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
					Name:  "name 3",
					Value: "item 3",
				},
			},
		},

		{
			Name:          "with valid table: array",
			ExpectedNames: []string{"name 1", "name 2"},
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
			ExpectedFatalMessage: "Expected a slice or array for the table, got string",
			Table:                "my table",
		},

		{
			Name:                 "with invalid table type: not array or slice of stucts",
			ExpectedFatalMessage: "Expected entry in table to be a struct, got string",
			Table: []string{
				"item 1",
				"item 2",
			},
		},

		{
			Name:                 "with missing name",
			ExpectedFatalMessage: "Name field does not exist on struct in table",
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
			ExpectedFatalMessage: "Name field in struct in table is not a string",
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
			ExpectedFatalMessage: "Errors encountered while building table:\n - table[1]: Name not set for item",
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
			ExpectedFatalMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1",
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
			ExpectedFatalMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1\n - table[3]: duplicate Name found; first occurrence was table[0].Name: name 1",
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
			mockT := setupMockT(t)
			expectedTableSize := len(entry.ExpectedNames)

			mockT.EXPECT().Helper().Times(3 * expectedTableSize) // 3 = RunTableByIndex + run + before Cleanup call
			mockT.EXPECT().Cleanup(gomock.Any()).Times(expectedTableSize)

			// Build expected Run calls only if there's no expected error
			expectedTestingInputs := []*testing.T{}
			if entry.ExpectedFatalMessage == "" {
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
				gomock.InOrder(
					mockT.EXPECT().Helper(),
					mockT.EXPECT().Cleanup(gomock.Any()),
					mockT.EXPECT().Helper(),
					mockT.EXPECT().Fatalf(entry.ExpectedFatalMessage),
				)
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

			// Verify call count
			if len(actualParams) != expectedTableSize {
				t.Fatalf("len(actualParams) != expectedTableSize")
			}

			// Verify parameters are correct
			for i, actualParam := range actualParams {
				if actualParam.ensure.T() != expectedTestingInputs[i] {
					t.Fatalf("actualParams[%d].ensure.T() != expectedTestingInputs[%d]", i, i)
				}

				if actualParam.i != i {
					t.Fatalf("actualParams[%d].i != %d", i, i)
				}
			}
		})
	}
}
