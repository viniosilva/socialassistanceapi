package unit

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestHealthServiceHealth(t *testing.T) {
	cases := map[string]struct {
		expectedRes service.Health
		prepareMock func(mock *mock.MockHealthStore)
	}{
		"should return health status up": {
			expectedRes: service.Health{Status: service.HealthStatusUp},
			prepareMock: func(mock *mock.MockHealthStore) {
				mock.EXPECT().Health(gomock.Any()).Return(true)
			},
		},
		"should return health status down": {
			expectedRes: service.Health{Status: service.HealthStatusDown},
			prepareMock: func(mock *mock.MockHealthStore) {
				mock.EXPECT().Health(gomock.Any()).Return(false)
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			healthMock := mock.NewMockHealthStore(ctrl)
			cs.prepareMock(healthMock)

			impl := service.NewHealthService(healthMock)

			// when
			res := impl.Health(ctx)

			// then
			if !reflect.DeepEqual(res, cs.expectedRes) {
				t.Errorf("HealthService.Health() = %v, expected %v", res, cs.expectedRes)
			}
		})
	}
}
