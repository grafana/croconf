package croconf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type testCaseGroup struct {
	name      string
	field     func(testSources) Field
	testCases []fieldTestCase
}

type testSources struct {
	json *SourceJSON
	env  *SourceEnvVars
	cli  *SourceCLI
}

type fieldTestCase struct {
	json           string
	env            []string
	cli            []string
	expectedValue  interface{}
	expectedErrors []string
}

var testCaseGroups = []testCaseGroup{ //nolint:gochecknoglobals
	{
		name: "simple int64 field",
		field: func(sources testSources) Field {
			var dest int64
			return NewInt64Field(
				&dest,
				DefaultIntValue(1),
				sources.json.From("vus"),
				sources.env.From("K6_VUS"),
				sources.cli.FromNameAndShorthand("vus", "u"),
			)
		},
		testCases: []fieldTestCase{
			{
				expectedValue: int64(1), // default, no sources
			},
			{
				json: `{"vus": "foo"}`,
				// TODO: improve this error message, something like `"foo" is not a valid integer value` would be much better
				expectedErrors: []string{`strconv.ParseInt: parsing "\"foo\"": invalid syntax`},
			},
			{
				json:          `{"vus": 2}`,
				expectedValue: int64(2),
			},
			{
				json:          `{"vus": 2}`,
				env:           []string{"K6_VUS=3"},
				expectedValue: int64(3),
			},
			{
				json:           `{"vus": 2}`,
				env:            []string{"K6_VUS=foo"},
				expectedErrors: []string{`BindInt64ValueTo: parsing "foo": invalid syntax`}, // TODO: better error message
			},
			{
				json: `{"vus": "foo"}`,
				env:  []string{"K6_VUS=bar"},
				expectedErrors: []string{ // TODO: better error messages
					`strconv.ParseInt: parsing "\"foo\"": invalid syntax`,
					`strconv.ParseInt: parsing "bar": invalid syntax`,
				},
			},
			{
				json:          `{"vus": 2}`,
				env:           []string{"K6_VUS=3"},
				cli:           []string{"--vus", "4"},
				expectedValue: int64(4),
			},
		},
	},
	{
		name: "simple string field",
		field: func(sources testSources) Field {
			var dest string
			return NewStringField(
				&dest,
				sources.json.From("fieldName"),
				sources.env.From("FIELD_NAME"),
				sources.cli.FromNameAndShorthand("field-name", "f"),
			)
		},
		testCases: []fieldTestCase{
			{
				expectedValue: "", // default, no sources
			},
			{
				json:          `{"fieldName": "foo"}`,
				expectedValue: "foo",
			},
			// TODO: add more test cases for this field
		},
	},
	{
		name: "int8 default",
		field: func(sources testSources) Field {
			var dest int8
			return NewInt8Field(
				&dest,
				DefaultIntValue(129),
			)
		},
		testCases: []fieldTestCase{
			{
				expectedErrors: []string{"invalid value 129, has to be between -128 and 127"},
			},
		},
	},
	{
		name: "int8 field",
		field: func(sources testSources) Field {
			var dest int8
			return NewInt8Field(
				&dest,
				DefaultIntValue(127),
				sources.json.From("tiny"),
				sources.env.From("K6_TINY"),
				sources.cli.FromName("tiny"),
			)
		},
		testCases: []fieldTestCase{
			{
				expectedValue: int8(127),
			},
			{
				json:          `{"tiny": -128}`,
				expectedValue: int8(-128),
			},
			{
				cli:            []string{"--tiny=-129"},
				expectedErrors: []string{`strconv.ParseInt: parsing "-129": value out of range`}, // TODO: better error message
			},
		},
	},
	{
		name: "int8 array",
		field: func(sources testSources) Field {
			var dest []int8
			return NewInt8SliceField(
				&dest,
				sources.json.From("tinyArr"),
				sources.env.From("TINY_ARR"),
			)
		},
		testCases: []fieldTestCase{
			// TODO: test defaults and null values
			{
				json:          `{"tinyArr": [1, 2]}`,
				expectedValue: []int8{1, 2},
			},
			{
				json:          `{"tinyArr": [1, 1, 2]}`,
				env:           []string{`TINY_ARR=3,5,8`},
				expectedValue: []int8{3, 5, 8},
			},
			{
				json: `{"tinyArr": [1, 255]}`,
				expectedErrors: []string{
					`strconv.ParseInt: parsing "255": value out of range`, // TODO: better error
				},
			},
		},
	},
	{
		name: "int64 array",
		field: func(sources testSources) Field {
			var dest []int64
			return NewInt64SliceField(
				&dest,
				sources.env.From("BIG_ARR"),
				sources.json.From("bigArr"),
			)
		},
		testCases: []fieldTestCase{
			// TODO: test defaults and null values
			{
				env:           []string{`BIG_ARR=1,2,1337`},
				expectedValue: []int64{1, 2, 1337},
			},
			{
				env:           []string{`BIG_ARR=1,2,1337`},
				json:          `{"bigArr": [0, 2, 4, 8]}`,
				expectedValue: []int64{0, 2, 4, 8},
			},
			{
				env:  []string{`BIG_ARR=1,2,foo`},
				json: `{"bigArr": [1, 2, null]}`,
				expectedErrors: []string{
					// TODO: better errors
					`strconv.ParseInt: parsing "foo": invalid syntax`,
					`strconv.ParseInt: parsing "null": invalid syntax`,
				},
			},
		},
	},
	{
		name: "nested config",
		field: func(sources testSources) Field {
			var dest string
			return NewStringField(
				&dest,
				sources.json.From("parent").From("child"),
			)
		},
		testCases: []fieldTestCase{
			{
				json:          `{"parent": {"child": "data"}}`,
				expectedValue: "data",
			},
			// TODO: add more test cases for this field?
		},
	},
	// TODO: add a lot more like these...
}

func runTestCase(t *testing.T, tcg testCaseGroup, tc fieldTestCase) {
	// TODO: actually test source failures as well?
	sources := testSources{
		json: NewJSONSource([]byte(tc.json)),
		env:  NewSourceFromEnv(tc.env),
		cli:  NewSourceFromCLIFlags(tc.cli),
	}

	field := tcg.field(sources)

	for _, s := range []Source{sources.json, sources.env, sources.cli} {
		if err := s.Initialize(); err != nil {
			t.Fatalf("unexpected error when initializing %s source: %s", s.GetName(), err)
		}
	}

	errs := field.Consolidate()
	if len(tc.expectedErrors) != len(errs) {
		t.Fatalf("Expected %d errors but got %d: %#v", len(tc.expectedErrors), len(errs), errs)
	}

	for i, expErr := range tc.expectedErrors {
		if !strings.Contains(errs[i].Error(), expErr) {
			t.Errorf("Expected error #%d to contain '%s' but it is '%s'", i, expErr, errs[i].Error())
		}
	}
	if len(tc.expectedErrors) == 0 {
		destPointer := field.Destination()
		value := reflect.Indirect(reflect.ValueOf(destPointer)).Interface()

		if !reflect.DeepEqual(tc.expectedValue, value) {
			t.Errorf(
				"Expected to get value '%#v' (%T), but got '%#v' (%T)",
				tc.expectedValue, tc.expectedValue, value, value)
		}
	}
}

func TestBuiltInFileds(t *testing.T) {
	t.Parallel()
	for i, tcg := range testCaseGroups {
		i, tcg := i, tcg
		t.Run(fmt.Sprintf("%03d: %s", i, tcg.name), func(t *testing.T) {
			t.Parallel()
			for j, tc := range tcg.testCases {
				j, tc := j, tc
				t.Run(fmt.Sprintf("TC#%03d", j), func(t *testing.T) {
					t.Parallel()
					runTestCase(t, tcg, tc)
				})
			}
		})
	}
}
