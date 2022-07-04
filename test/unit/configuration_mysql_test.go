package unit

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
			if !reflect.DeepEqual(fields, cs.expectedFields) {
				t.Errorf("MySQLConfiguration.GetNotEmptyFields() fields: = %v, expected %v", fields, cs.expectedFields)
			}
			if !reflect.DeepEqual(values, cs.expectedValues) {
				t.Errorf("MySQLConfiguration.GetNotEmptyFields() values: = %v, expected %v", values, cs.expectedValues)
			}
		})
	}
}
