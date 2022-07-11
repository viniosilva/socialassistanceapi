package configuration_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/viniosilva/socialassistanceapi/internal/configuration"
)

func TestConfigurationMySQLBuildUpdateData(t *testing.T) {
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
				"float32":     1.0,
				"boolean":     false,
			},
			expectedFields: []string{
				"string = ?",
				"int = ?",
				"float32 = ?",
				"boolean = ?",
			},
			expectedValues: []interface{}{"string", 0, 1.0, false},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			impl := configuration.NewMySQL("", time.Duration(0), 0, 0)

			// when
			fields, values := impl.BuildUpdateData(cs.inputData)

			// then
			if len(fields) != len(cs.expectedFields) {
				t.Errorf("MySQLConfiguration.GetNotEmptyFields() fields: = %v, expected %v", fields, cs.expectedFields)
			}
			if len(values) != len(cs.expectedValues) {
				t.Errorf("MySQLConfiguration.GetNotEmptyFields() values: = %v, expected %v", values, cs.expectedValues)
			}

			for _, e := range cs.expectedFields {
				for i, f := range fields {
					if reflect.DeepEqual(e, f) {
						break
					}

					if i == len(fields)-1 {
						t.Errorf("MySQLConfiguration.GetNotEmptyFields() fields: = %v, expected %v", fields, cs.expectedFields)
					}
				}
			}

			for _, e := range cs.expectedValues {
				for i, v := range values {
					if reflect.DeepEqual(e, v) {
						break
					}

					if i == len(values)-1 {
						t.Errorf("MySQLConfiguration.GetNotEmptyFields() values: = %v, expected %v", values, cs.expectedValues)
					}
				}
			}
		})
	}
}
