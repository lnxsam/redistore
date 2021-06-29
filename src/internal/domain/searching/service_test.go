package searching

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
	assert.True(t, ok, "instance should be of type searching.Service")
	assert.NotNil(t, a, "instance should not be nil")
}
func TestSearchProductsByTitle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	type mockSearchProductsByTitleInputs struct {
		ctx           context.Context
		titleKeywords string
	}

	type mockSearchProductsByTitleOutputs struct {
		products []domain.Product
		err      error
	}

	type GetProductListInput struct {
		ctx           context.Context
		titleKeywords string
	}
	type expected struct {
		products []domain.Product
		err      error
	}
	repoErr := yerror.E(errors.New("error occurred in repository"))
	products := factories.Product.CreateMany(2)
	titleKeywords := "Title"

	testCases := []struct {
		name                             string
		mockSearchProductsByTitleInputs  mockSearchProductsByTitleInputs
		mockSearchProductsByTitleOutputs mockSearchProductsByTitleOutputs
		GetProductListInput              GetProductListInput
		expected                         expected
	}{
		{
			name: "get error in SearchProductsByTitle",
			mockSearchProductsByTitleInputs: mockSearchProductsByTitleInputs{
				ctx:           ctx,
				titleKeywords: titleKeywords,
			},
			mockSearchProductsByTitleOutputs: mockSearchProductsByTitleOutputs{
				products: nil,
				err:      repoErr,
			},
			GetProductListInput: GetProductListInput{
				ctx:           ctx,
				titleKeywords: titleKeywords,
			},
			expected: expected{
				products: nil,
				err:      repoErr,
			},
		},
		{
			name: "successful test",
			mockSearchProductsByTitleInputs: mockSearchProductsByTitleInputs{
				ctx:           ctx,
				titleKeywords: titleKeywords,
			},
			mockSearchProductsByTitleOutputs: mockSearchProductsByTitleOutputs{
				products: products,
				err:      nil,
			},
			GetProductListInput: GetProductListInput{
				ctx:           ctx,
				titleKeywords: titleKeywords,
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

		repositoryMock.On("SearchProductsByTitle", mock.AnythingOfType("*context.timerCtx"),
			tc.mockSearchProductsByTitleInputs.titleKeywords).
			Return(tc.mockSearchProductsByTitleOutputs.products, tc.mockSearchProductsByTitleOutputs.err).Once()

		got, gotErr := aa.SearchProductsByTitle(tc.GetProductListInput.ctx, tc.GetProductListInput.titleKeywords)
		if tc.expected.err != nil {
			assert.NotNil(t, gotErr, tc.name)
		} else {
			assert.EqualValues(t, tc.expected.products, got, tc.name)
		}
	}
	repositoryMock.AssertExpectations(t)
}
