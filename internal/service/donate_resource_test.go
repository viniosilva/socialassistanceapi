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

func TestDonateResourceServiceDonate(t *testing.T) {
	cases := map[string]struct {
		inputDto    service.DonateResourceDonateDto
		expectedErr error
		prepareMock func(mockDonateResourceRepository *mock.MockDonateResourceRepository)
	}{
		"should donate resource": {
			inputDto: service.DonateResourceDonateDto{
				ResourceID: 1,
				AddressID:  1,
				Quantity:   1,
			},
			prepareMock: func(mockDonateResourceRepository *mock.MockDonateResourceRepository) {
				mockDonateResourceRepository.EXPECT().Donate(gomock.Any(), 1, 1, 1.0).Return(nil)
			},
		},
		"should throw error": {
			inputDto: service.DonateResourceDonateDto{
				ResourceID: 1,
				AddressID:  1,
				Quantity:   1,
			},
			expectedErr: fmt.Errorf("error"),
			prepareMock: func(mockDonateResourceRepository *mock.MockDonateResourceRepository) {
				mockDonateResourceRepository.EXPECT().Donate(gomock.Any(), 1, 1, 1.0).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockDonateResourceRepository := mock.NewMockDonateResourceRepository(ctrl)
			cs.prepareMock(mockDonateResourceRepository)

			impl := &service.DonateResourceServiceImpl{DonateResourceRepository: mockDonateResourceRepository}

			// when
			err := impl.Donate(ctx, cs.inputDto)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestDonateResourceServiceReturn(t *testing.T) {
	cases := map[string]struct {
		inputResourceID int
		expectedErr     error
		prepareMock     func(mockDonateResourceRepository *mock.MockDonateResourceRepository)
	}{
		"should return resource": {
			inputResourceID: 1,
			prepareMock: func(mockDonateResourceRepository *mock.MockDonateResourceRepository) {
				mockDonateResourceRepository.EXPECT().Return(gomock.Any(), 1).Return(nil)
			},
		},
		"should throw error": {
			inputResourceID: 1,
			expectedErr:     fmt.Errorf("error"),
			prepareMock: func(mockDonateResourceRepository *mock.MockDonateResourceRepository) {
				mockDonateResourceRepository.EXPECT().Return(gomock.Any(), 1).Return(fmt.Errorf("error"))
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			mockDonateResourceRepository := mock.NewMockDonateResourceRepository(ctrl)
			cs.prepareMock(mockDonateResourceRepository)

			impl := &service.DonateResourceServiceImpl{DonateResourceRepository: mockDonateResourceRepository}

			// when
			err := impl.Return(ctx, cs.inputResourceID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
