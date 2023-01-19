package infra_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/socialassistanceapi/internal/infra"
)

func Test_MySQLBuild_UpdateData(t *testing.T) {
	cases := map[string]struct {
		inputData      map[string]interface{}
		expectedFields []string
		expectedValues []interface{}
	}{
		"should return": {
			inputData: map[string]interface{}{
				"string":      "string",
				"emptystring": "",
				"int":         0,
				"float64":     1.0,
				"float64_no":  0.0,
				"boolean":     false,
			},
			expectedFields: []string{
				"string = ?",
				"int = ?",
				"float64 = ?",
				"boolean = ?",
			},
			expectedValues: []interface{}{"string", 0, 1.0, false},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			impl := infra.MySQLConfigure("", 0, "", "", "", time.Duration(0), 0, 0)

			// when
			fields, values := impl.BuildUpdateData(cs.inputData)

			// then
			assert.Equal(t, len(cs.expectedFields), len(fields))
			assert.Equal(t, len(cs.expectedValues), len(values))

			for _, e := range cs.expectedFields {
				for i, f := range fields {
					if reflect.DeepEqual(e, f) {
						break
					}

					assert.LessOrEqual(t, i, len(fields)-1)
				}
			}

			for _, e := range cs.expectedValues {
				for i, v := range values {
					if reflect.DeepEqual(e, v) {
						break
					}

					assert.LessOrEqual(t, i, len(values)-1)
				}
			}
		})
	}
}
