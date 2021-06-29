package listing

import (
	"context"
	"errors"
	"redistore/internal/domain"
	"redistore/internal/domain/factories"
	"redistore/internal/domain/ports/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"redistore/pkg/yerror"
)

func TestNew(t *testing.T) {
	repository := new(mocks.Repository)
	a, ok := New(repository).(Service)
	assert.True(t, ok, "instance should be of type listing.Service")
	assert.NotNil(t, a, "instance should not be nil")
}
func TestGetProductList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	type mockGetProductListInputs struct {
		ctx context.Context
	}

	type mockGetProductListOutputs struct {
		products []domain.Product
		err      error
	}

	type GetProductListInput struct {
		ctx context.Context
	}
	type expected struct {
		products []domain.Product
		err      error
	}
	repoErr := yerror.E(errors.New("error occurred in repository"))
	products := factories.Product.CreateMany(2)

	testCases := []struct {
		name                      string
		mockGetProductListInputs  mockGetProductListInputs
		mockGetProductListOutputs mockGetProductListOutputs
		GetProductListInput       GetProductListInput
		expected                  expected
	}{
		{
			name: "get error in GetProductList",
			mockGetProductListInputs: mockGetProductListInputs{
				ctx: ctx,
			},
			mockGetProductListOutputs: mockGetProductListOutputs{
				products: nil,
				err:      repoErr,
			},
			GetProductListInput: GetProductListInput{
				ctx: ctx,
			},
			expected: expected{
				products: nil,
				err:      repoErr,
			},
		},
		{
			name: "successful test",
			mockGetProductListInputs: mockGetProductListInputs{
				ctx: ctx,
			},
			mockGetProductListOutputs: mockGetProductListOutputs{
				products: products,
				err:      nil,
			},
			GetProductListInput: GetProductListInput{
				ctx: ctx,
			},
			expected: expected{
				products: products,
				err:      nil,
			},
		},
	}

	repositoryMock := new(mocks.Repository)
	aa := New(repositoryMock)

	for _, tc := range testCases {

		repositoryMock.On("GetProductList", mock.AnythingOfType("*context.timerCtx")).
			Return(tc.mockGetProductListOutputs.products, tc.mockGetProductListOutputs.err).Once()

		got, gotErr := aa.GetProductList(tc.GetProductListInput.ctx)
		if tc.expected.err != nil {
			assert.NotNil(t, gotErr, tc.name)
		} else {
			assert.EqualValues(t, tc.expected.products, got, tc.name)
		}
	}
	repositoryMock.AssertExpectations(t)
}
