package ensuring_test

import (
	"strings"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/ensuring/internal/testhelper"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"go.uber.org/mock/gomock"
)

func TestERunTableByIndex(t *testing.T) {
	runTableConfig{
		prepare: func(ensure ensuring.E) func(table any, fn func(ensure ensuring.E, i int)) {
			return ensure.RunTableByIndex
		},
	}.test(t)
}

type runTableTestEntryGroup struct {
	Prefix  string
	Entries []runTableTestEntry
}

type runTableTestEntry struct {
	Name                 string
	ExpectedNames        []string
	ExpectedTableSize    int // Defaults to the length of ExpectedNames
	FatalMessagesContain []string
	Table                any
	CheckEntry           func(t *testing.T, rawEntry any)
}

type runTableConfig struct {
	isSync  bool
	prepare func(ensure ensuring.E) func(table any, fn func(ensure ensuring.E, i int))
}

func (cfg runTableConfig) test(t *testing.T) {
	runTableTests := runTableTests{}

	groups := []runTableTestEntryGroup{
		runTableTests.general(),
		runTableTests.mocksField(),
		runTableTests.setupMocksField(),
		runTableTests.subjectField(),
	}

	table := []runTableTestEntry{}
	for _, group := range groups {
		for _, entry := range group.Entries {
			entry.Name = group.Prefix + ": " + entry.Name
			table = append(table, entry)
		}
	}

	for _, entry := range table {
		t.Run(entry.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			outerMockT := setupMockTWithCleanupCheck(t)
			outerMockT.EXPECT().Helper().MinTimes(1)

			outerMockCtx := mock_testctx.NewMockContext(ctrl)
			outerMockCtx.EXPECT().T().Return(outerMockT).AnyTimes()
			testhelper.SetTestContext(t, outerMockT, outerMockCtx)

			expectedRunCalls := []any{}
			innerMockTs := []*mock_testctx.MockT{} //lint:ignore ST1003 mockTs not mockTS

			actualFatalMessages := []string{}
			fatalMessagesRecorder := func(msg string, args ...any) {
				actualFatalMessages = append(actualFatalMessages, msg)
			}

			for _, name := range entry.ExpectedNames {
				innerMockT := setupMockT(t)
				innerMockT.EXPECT().Helper().MinTimes(cfg.expectedMinHelperCalls(entry))
				innerMockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()
				innerMockT.EXPECT().Fatalf(gomock.Any()).Do(fatalMessagesRecorder).AnyTimes()
				innerMockTs = append(innerMockTs, innerMockT)

				innerMockCtx := mock_testctx.NewMockContext(ctrl)
				innerMockCtx.EXPECT().T().Return(innerMockT).AnyTimes()
				innerMockCtx.EXPECT().GoMockController().Return(gomock.NewController(innerMockT)).AnyTimes()
				innerMockCtx.EXPECT().Ensure().Return(ensure.New(innerMockT)).AnyTimes()
				testhelper.SetTestContext(t, innerMockT, innerMockCtx)

				if cfg.isSync {
					preSyncInnerMockT := setupMockT(t)
					preSyncInnerMockT.EXPECT().Helper().MinTimes(cfg.expectedMinHelperCalls(entry))
					preSyncInnerMockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()
					preSyncInnerMockT.EXPECT().Fatalf(gomock.Any()).Do(fatalMessagesRecorder).AnyTimes()

					preSyncInnerMockCtx := mock_testctx.NewMockSyncableContext(ctrl)
					preSyncInnerMockCtx.EXPECT().T().Return(preSyncInnerMockT).AnyTimes()
					preSyncInnerMockCtx.EXPECT().GoMockController().Return(gomock.NewController(preSyncInnerMockT)).AnyTimes()
					preSyncInnerMockCtx.EXPECT().Ensure().Return(ensure.New(preSyncInnerMockT)).AnyTimes()
					testhelper.SetTestContext(t, preSyncInnerMockT, preSyncInnerMockCtx)

					expectedRunCalls = append(expectedRunCalls,
						outerMockCtx.EXPECT().Run(name, gomock.Any()).
							Do(execFuncParamWithName(preSyncInnerMockCtx)),
					)

					preSyncInnerMockCtxSyncCall := preSyncInnerMockCtx.EXPECT().Sync(gomock.Any()).
						Do(execFuncParam(innerMockCtx))

					if len(entry.FatalMessagesContain) > 0 {
						preSyncInnerMockCtxSyncCall.AnyTimes()
					}

					expectedRunCalls = append(expectedRunCalls, preSyncInnerMockCtxSyncCall)
				} else {
					expectedRunCalls = append(expectedRunCalls,
						outerMockCtx.EXPECT().Run(name, gomock.Any()).
							Do(execFuncParamWithName(innerMockCtx)),
					)
				}
			}

			// Run calls should be in order
			gomock.InOrder(expectedRunCalls...)

			outerMockT.EXPECT().Fatalf(gomock.Any()).Do(func(msg string, args ...any) {
				actualFatalMessages = append(actualFatalMessages, msg)
			}).AnyTimes()

			type entryCall struct {
				ensure ensuring.E
				i      int
			}

			// Run table and save call details
			actualEntryCalls := []entryCall{}
			ensure := ensure.New(outerMockT)
			runTableByIndex := cfg.prepare(ensure)
			runTableByIndex(entry.Table, func(ensure ensuring.E, i int) {
				actualEntryCalls = append(actualEntryCalls, entryCall{ensure: ensure, i: i})
			})

			// Verify entry calls
			{
				expectedTableSize := entry.ExpectedTableSize
				if expectedTableSize == 0 {
					expectedTableSize = len(entry.ExpectedNames)
				}

				// Verify call count
				if len(actualEntryCalls) != expectedTableSize {
					t.Fatalf("len(actualParams) != expectedTableSize: %d != %d", len(actualEntryCalls), expectedTableSize)
				}

				// Verify parameters are correct
				for i, actualParam := range actualEntryCalls {
					if actualParam.i != i {
						t.Fatalf("actualParams[%d].i != %d", i, i)
					}

					// Show the correct mock is paired with the correct call
					innerMockT := innerMockTs[i]
					innerMockT.EXPECT().Fatalf("failing %d", i)
					actualParam.ensure.Failf("failing %d", i)
				}
			}

			// Verify fatal messages
			{
				if len(actualFatalMessages) != len(entry.FatalMessagesContain) {
					t.Fatalf("Expected %d fatal message(s), got %d fatal message(s): %v", len(entry.FatalMessagesContain), len(actualFatalMessages), actualFatalMessages)
				}

				for i, expectedMessageContains := range entry.FatalMessagesContain {
					if !strings.Contains(actualFatalMessages[i], expectedMessageContains) {
						t.Fatalf("Error message expected to contain %q, got: %s", expectedMessageContains, actualFatalMessages[i])
					}
				}
			}

			if entry.CheckEntry != nil {
				entry.CheckEntry(t, entry.Table)
			}
		})
	}
}

func (cfg runTableConfig) expectedMinHelperCalls(entry runTableTestEntry) int {
	if len(entry.FatalMessagesContain) > 0 {
		return 0
	}

	return 2
}

type runTableTests struct{}

func (runTableTests) general() runTableTestEntryGroup {
	return runTableTestEntryGroup{
		Prefix: "general",
		Entries: []runTableTestEntry{
			{
				Name:          "with valid table: slice with non-pointers",
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
				Name:          "with valid table: slice with pointers",
				ExpectedNames: []string{"name 1", "name 2", "name 3"},
				Table: []*struct {
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
				Name:          "with valid table: array with non-pointers",
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
				Name:          "with valid table: array with pointers",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: [2]*struct {
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
				FatalMessagesContain: []string{"Expected a slice or array for the table, got string"},
				Table:                "my table",
			},

			{
				Name:                 "with invalid table type: not array or slice of structs",
				FatalMessagesContain: []string{"Expected entry in table to be a struct or a pointer to a struct, got string"},
				Table: []string{
					"item 1",
					"item 2",
				},
			},

			{
				Name:                 "with missing name",
				FatalMessagesContain: []string{"Required Name field does not exist on struct in table"},
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
				FatalMessagesContain: []string{"Required Name field in struct in table is not a string"},
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
				ExpectedNames:        []string{"name 1", ""},
				ExpectedTableSize:    1,
				FatalMessagesContain: []string{"table[1].Name is empty"},
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
				ExpectedNames:        []string{"name 1", "name 2", "name 1"},
				ExpectedTableSize:    2,
				FatalMessagesContain: []string{"table[2].Name duplicates table[0].Name: name 1"},
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
				ExpectedNames:        []string{"name 1", "name 2", "name 1", "name 1"},
				ExpectedTableSize:    2,
				FatalMessagesContain: []string{"table[2].Name duplicates table[0].Name: name 1", "table[3].Name duplicates table[0].Name: name 1"},
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
		},
	}
}

func (runTableTests) mocksField() runTableTestEntryGroup {
	type (
		//lint:ignore U1000 Present for testing purposes
		TwoValidMocksWithUnexported struct {
			Valid1      *ExampleMockValid1
			notExported string //nolint:structcheck // Present for the test
			Valid2      *ExampleMockValid2
		}

		Embedable struct {
			Valid1 *ExampleMockValid1
		}

		TwoValidMocksWithEmbedded struct {
			Embedable
			Valid2 *ExampleMockValid2
		}

		TwoValidMocksWithEmbeddedPtr struct {
			*Embedable
			Valid2 *ExampleMockValid2
		}

		OneMockNEWMethodZeroParams struct {
			Valid1     *ExampleMockValid1
			ZeroParams *ExampleMockNEWMethodZeroParams
			Valid2     *ExampleMockValid2
		}

		BrokenEmbedable struct {
			Valid1 ExampleMockValid2 // Not a pointer
		}

		BrokenEmbedded struct {
			BrokenEmbedable
			Valid2 *ExampleMockValid2
		}

		OneMockMissingNEWMethod struct {
			Valid1  *ExampleMockValid1
			Invalid *struct{ Nothing bool }
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodExtraParam struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodExtraParam
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodIncorrectParam struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodIncorrectParam
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodZeroReturns struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodZeroReturns
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodIncorrectReturn struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodIncorrectReturn
			Valid2  *ExampleMockValid2
		}

		OneMockNotPointer struct {
			Valid1  *ExampleMockValid1
			Invalid ExampleMockValid1
			Valid2  *ExampleMockValid2
		}

		DuplicateMocks struct {
			Valid1          *ExampleMockValid1
			Valid1Duplicate *ExampleMockValid1
		}
	)

	return runTableTestEntryGroup{
		Prefix: "Mocks field",
		Entries: []runTableTestEntry{
			{
				Name:          "when valid",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *TwoValidMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *TwoValidMocks
					})

					for _, entry := range table {
						entry.Mocks.check(t)
					}
				},
			},

			{
				Name:          "when valid with unexported field",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *TwoValidMocksWithUnexported
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *TwoValidMocksWithUnexported
					})

					for _, entry := range table {
						checkTwoValidMocks(t, entry.Mocks.Valid1, entry.Mocks.Valid2)
					}
				},
			},

			{
				Name:          "when valid with embedded field",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *TwoValidMocksWithEmbedded
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *TwoValidMocksWithEmbedded
					})

					for _, entry := range table {
						checkTwoValidMocks(t, entry.Mocks.Valid1, entry.Mocks.Valid2)
					}
				},
			},

			{
				Name:          "when valid with NEW method with no params",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodZeroParams
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *OneMockNEWMethodZeroParams
					})

					for _, entry := range table {
						checkTwoValidMocks(t, entry.Mocks.Valid1, entry.Mocks.Valid2)
						isTrue(t, entry.Mocks.ZeroParams.WasInitialized)
					}
				},
			},

			{
				Name:          "when embedded field is not struct",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *TwoValidMocksWithEmbeddedPtr
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *TwoValidMocksWithEmbeddedPtr
					})

					for _, entry := range table {
						checkTwoValidMocks(t, entry.Mocks.Valid1, entry.Mocks.Valid2)
					}
				},
			},

			{
				Name:                 "when embedded field has error",
				FatalMessagesContain: []string{"Mocks.BrokenEmbedable.Valid1 is expected to be a pointer to a struct"},
				Table: []struct {
					Name  string
					Mocks *BrokenEmbedded
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when not pointer to mock struct",
				FatalMessagesContain: []string{"expected Mocks field to be a pointer to a struct"},
				Table: []struct {
					Name  string
					Mocks TwoValidMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when pointer to non struct",
				FatalMessagesContain: []string{"expected Mocks field to be a pointer to a struct"},
				Table: []struct {
					Name  string
					Mocks *string
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when missing NEW method",
				FatalMessagesContain: []string{"Mocks.Invalid (*struct { Nothing bool }) is missing a NEW method"},
				Table: []struct {
					Name  string
					Mocks *OneMockMissingNEWMethod
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has an extra param",
				FatalMessagesContain: []string{"Mocks.Invalid (*ensuring_test.ExampleMockNEWMethodExtraParam) must have a NEW method matching one of the following signatures"},
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodExtraParam
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has incorrect param",
				FatalMessagesContain: []string{"Mocks.Invalid (*ensuring_test.ExampleMockNEWMethodIncorrectParam) must have a NEW method matching one of the following signatures"},
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodIncorrectParam
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has zero returns",
				FatalMessagesContain: []string{"Mocks.Invalid (*ensuring_test.ExampleMockNEWMethodZeroReturns) must have a NEW method matching one of the following signatures"},
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodZeroReturns
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has incorrect return",
				FatalMessagesContain: []string{"Mocks.Invalid (*ensuring_test.ExampleMockNEWMethodIncorrectReturn) must have a NEW method matching one of the following signatures"},
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodIncorrectReturn
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when mock is not a pointer",
				FatalMessagesContain: []string{"Mocks.Invalid is expected to be a pointer to a struct"},
				Table: []struct {
					Name  string
					Mocks *OneMockNotPointer
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:          "with duplicate mock",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *DuplicateMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name  string
						Mocks *DuplicateMocks
					})

					for _, entry := range table {
						isTrue(t, entry.Mocks.Valid1.WasInitialized)
						isTrue(t, entry.Mocks.Valid1Duplicate.WasInitialized)
					}
				},
			},
		},
	}
}

func (runTableTests) setupMocksField() runTableTestEntryGroup {
	return runTableTestEntryGroup{
		Prefix: "SetupMocks field",
		Entries: []runTableTestEntry{
			{
				Name:          "with valid function with one param",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name: "name 1",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 1"
						},
					},
					{
						Name: "name 2",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 2"
						},
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name       string
						Mocks      *TwoValidMocks
						SetupMocks func(*TwoValidMocks)
					})

					for _, entry := range table {
						entry.Mocks.check(t)
						isTrue(t, entry.Mocks.Valid1.CustomField == "updated "+entry.Name)
					}
				},
			},
			{
				Name:                 "with valid function with two params",
				ExpectedNames:        []string{"name 1", "name 2"},
				FatalMessagesContain: []string{"first SetupMocks", "second SetupMocks"}, // Not actual failures; only to show ensure is passed in correctly
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks, ensuring.E)
				}{
					{
						Name: "name 1",
						SetupMocks: func(tvm *TwoValidMocks, ensure ensuring.E) {
							tvm.Valid1.CustomField = "updated name 1"
							ensure.Failf("first SetupMocks")
						},
					},
					{
						Name: "name 2",
						SetupMocks: func(tvm *TwoValidMocks, ensure ensuring.E) {
							tvm.Valid1.CustomField = "updated name 2"
							ensure.Failf("second SetupMocks")
						},
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name       string
						Mocks      *TwoValidMocks
						SetupMocks func(*TwoValidMocks, ensuring.E)
					})

					for _, entry := range table {
						entry.Mocks.check(t)
						isTrue(t, entry.Mocks.Valid1.CustomField == "updated "+entry.Name)
					}
				},
			},

			{
				Name:          "with function not present for one",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name: "name 1",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 1"
						},
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name       string
						Mocks      *TwoValidMocks
						SetupMocks func(*TwoValidMocks)
					})

					isTrue(t, table[0].Mocks.Valid1.CustomField == "updated name 1")
					isTrue(t, table[1].Mocks.Valid1.CustomField == "")
				},
			},

			{
				Name:                 "without Mocks field",
				FatalMessagesContain: []string{"Mocks field must be set on the table to use SetupMocks"},
				Table: []struct {
					Name       string
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*TwoValidMocks) {},
					},
					{
						Name:       "name 2",
						SetupMocks: func(*TwoValidMocks) {},
					},
				},
			},

			{
				Name:                 "function missing param",
				FatalMessagesContain: []string{"expected SetupMocks field to be one of the following:"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func()
				}{
					{
						Name:       "name 1",
						SetupMocks: func() {},
					},
					{
						Name:       "name 2",
						SetupMocks: func() {},
					},
				},
			},

			{
				Name:                 "function with invalid param",
				FatalMessagesContain: []string{"expected SetupMocks field to be one of the following:"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*string)
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*string) {},
					},
					{
						Name:       "name 2",
						SetupMocks: func(*string) {},
					},
				},
			},

			{
				Name:                 "function with a return",
				FatalMessagesContain: []string{"expected SetupMocks field to be one of the following:"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks) error
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*TwoValidMocks) error { return nil },
					},
					{
						Name:       "name 2",
						SetupMocks: func(*TwoValidMocks) error { return nil },
					},
				},
			},
		},
	}
}

func (runTableTests) subjectField() runTableTestEntryGroup {
	type (
		OneValidMock struct {
			Valid1 *ExampleMockValid1
		}

		TwoValidMocksWithIgnoreUnusedTag struct {
			Valid1 *ExampleMockValid1
			Valid2 *ExampleMockValid2 `ensure:"ignoreunused"`
		}

		IntAdder interface {
			Add(a, b int) int
		}

		IntMultiplier interface {
			Multiply(a, b int) int
		}

		AdderSubject struct {
			Adder IntAdder
		}

		MultiInterfaceSubject struct {
			Adder      IntAdder
			Multiplier IntMultiplier
		}

		AdderSubjectWithDuplicate struct {
			Adder1 IntAdder
			Adder2 IntAdder
		}

		AdderSubjectWithExtraField struct {
			Adder      IntAdder
			ExtraField string
		}

		AdderSubjectWithUnmockedInterface struct {
			Adder             IntAdder
			UnmockedInterface IntMultiplier
		}

		SubjectMatchingMultipleMocks struct {
			Subber interface{ Sub(a, b int) int }
		}
	)

	return runTableTestEntryGroup{
		Prefix: "Subject field",
		Entries: []runTableTestEntry{
			{
				Name:          "when valid",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject *AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *OneValidMock
						Subject *AdderSubject
					})

					for _, entry := range table {
						isTrue(t, entry.Mocks.Valid1.WasInitialized)
						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
					}
				},
			},

			{
				Name:          "when valid with no mocks",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Subject *MultiInterfaceSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Subject *MultiInterfaceSubject
					})

					for _, entry := range table {
						isTrue(t, entry.Subject != nil)
					}
				},
			},

			{
				Name:          "when duplicate interfaces",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject *AdderSubjectWithDuplicate
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *OneValidMock
						Subject *AdderSubjectWithDuplicate
					})

					for _, entry := range table {
						isTrue(t, entry.Mocks.Valid1.WasInitialized)

						isTrue(t, entry.Subject.Adder1.Add(1, 2) == 3)
						isTrue(t, entry.Subject.Adder2.Add(1, 2) == 3)
						isTrue(t, entry.Subject.Adder1 == entry.Subject.Adder2) // Should point to the same mock
					}
				},
			},

			{
				Name:                 "when not pointer to struct",
				FatalMessagesContain: []string{"expected Subject field to be a pointer to a struct"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when pointer to non struct",
				FatalMessagesContain: []string{"expected Subject field to be a pointer to a struct"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject *string
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:          "when field contains a non-interface field",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject *AdderSubjectWithExtraField
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *OneValidMock
						Subject *AdderSubjectWithExtraField
					})

					for _, entry := range table {
						isTrue(t, entry.Mocks.Valid1.WasInitialized)

						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
						isTrue(t, entry.Subject.ExtraField == "")
					}
				},
			},

			{
				Name:          "when field contains a non-mocked interface",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *OneValidMock
					Subject *AdderSubjectWithUnmockedInterface
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *OneValidMock
						Subject *AdderSubjectWithUnmockedInterface
					})

					for _, entry := range table {
						isTrue(t, entry.Mocks.Valid1.WasInitialized)

						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
						isTrue(t, entry.Subject.UnmockedInterface == nil)
					}
				},
			},

			{
				Name:                 "when entry matches multiple mocks",
				FatalMessagesContain: []string{"Subject.Subber is satisfied by more than one mock: Mocks.Valid1, Mocks.Valid2."},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *SubjectMatchingMultipleMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when mock is unused",
				FatalMessagesContain: []string{"Mocks.Valid2 was required but not matched by any interfaces in Subject."},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:          "when mock is unused but has ignoreunused tag",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocksWithIgnoreUnusedTag
					Subject *AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable any) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *TwoValidMocksWithIgnoreUnusedTag
						Subject *AdderSubject
					})

					for _, entry := range table {
						checkTwoValidMocks(t, entry.Mocks.Valid1, entry.Mocks.Valid2)
						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
					}
				},
			},
		},
	}
}

func isTrue(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Errorf("value is not true")
	}
}

type TwoValidMocks struct {
	Valid1 *ExampleMockValid1
	Valid2 *ExampleMockValid2
}

func (tvm *TwoValidMocks) check(t *testing.T) {
	t.Helper()
	checkTwoValidMocks(t, tvm.Valid1, tvm.Valid2)
}

func checkTwoValidMocks(t *testing.T, valid1 *ExampleMockValid1, valid2 *ExampleMockValid2) {
	t.Helper()

	isTrue(t, valid1.WasInitialized)
	isTrue(t, valid2.WasInitialized)
	isTrue(t, valid1.GoMockController == valid2.GoMockController) // Ensure GoMock Controller is memoized
}

type ExampleMockValid1 struct {
	WasInitialized   bool
	GoMockController *gomock.Controller
	CustomField      string
}

func (m *ExampleMockValid1) NEW(ctrl *gomock.Controller) *ExampleMockValid1 {
	if ctrl == nil {
		panic("GoMock controller is nil")
	}

	return &ExampleMockValid1{WasInitialized: true, GoMockController: ctrl}
}

func (m *ExampleMockValid1) Add(a, b int) int {
	return a + b
}

func (m *ExampleMockValid1) Sub(a, b int) int {
	return a - b
}

type ExampleMockValid2 struct {
	WasInitialized   bool
	GoMockController *gomock.Controller
	CustomField      string
}

func (m *ExampleMockValid2) NEW(ctrl *gomock.Controller) *ExampleMockValid2 {
	if ctrl == nil {
		panic("GoMock controller is nil")
	}

	return &ExampleMockValid2{WasInitialized: true, GoMockController: ctrl}
}

func (m *ExampleMockValid2) Sub(a, b int) int {
	return a - b
}

type ExampleMockNEWMethodZeroParams struct {
	WasInitialized bool
}

func (m *ExampleMockNEWMethodZeroParams) NEW() *ExampleMockNEWMethodZeroParams {
	return &ExampleMockNEWMethodZeroParams{WasInitialized: true}
}

type ExampleMockNEWMethodExtraParam struct{}

func (m *ExampleMockNEWMethodExtraParam) NEW(ctrl *gomock.Controller, extra string) *ExampleMockNEWMethodExtraParam {
	return nil
}

type ExampleMockNEWMethodIncorrectParam struct{}

func (m *ExampleMockNEWMethodIncorrectParam) NEW(notGoMockCtrl string) *ExampleMockNEWMethodIncorrectParam {
	return nil
}

type ExampleMockNEWMethodZeroReturns struct{}

func (m *ExampleMockNEWMethodZeroReturns) NEW(ctrl *gomock.Controller) {}

type ExampleMockNEWMethodIncorrectReturn struct{}

func (m *ExampleMockNEWMethodIncorrectReturn) NEW(ctrl *gomock.Controller) string { return "" }
