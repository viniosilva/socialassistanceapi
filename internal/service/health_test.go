package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/socialassistanceapi/internal/service"
	"github.com/viniosilva/socialassistanceapi/mock"
)

func TestHealthServicePing(t *testing.T) {
	cases := map[string]struct {
		expectedRes service.HealthResponse
		prepareMock func(mockHealthRepository *mock.MockHealthRepository)
	}{
		"should return health status up": {
			expectedRes: service.HealthResponse{Status: service.HealthStatusUp},
			prepareMock: func(mockHealthRepository *mock.MockHealthRepository) {
				mockHealthRepository.EXPECT().Ping(gomock.Any()).Return(nil)
			},
		},
		"should return health status down": {
			expectedRes: service.HealthResponse{Status: service.HealthStatusDown},
			prepareMock: func(mockHealthRepository *mock.MockHealthRepository) {
				mockHealthRepository.EXPECT().Ping(gomock.Any()).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			mockHealthRepository := mock.NewMockHealthRepository(ctrl)
			cs.prepareMock(mockHealthRepository)

			impl := &service.HealthServiceImpl{HealthRepository: mockHealthRepository}

			// when
			res := impl.Ping(ctx)

			// then
			assert.Equal(t, cs.expectedRes, res)
		})
	}
}
