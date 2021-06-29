package creating

import (
	"context"
	"errors"
	"redistore/internal/domain/factories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"redistore/internal/domain"
	"redistore/internal/domain/ports/mocks"
	"redistore/pkg/yerror"
)

func TestNew(t *testing.T) {
	repository := new(mocks.Repository)
	a, ok := New(repository).(Service)
	assert.True(t, ok, "instance should be of type creating.Service")
	assert.NotNil(t, a, "instance should not be nil")
}
func TestCreateProduct(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	type mockInsertInputs struct {
		ctx     context.Context
		product domain.Product
	}

	type mockInsertOutputs struct {
		product *domain.Product
		err     error
	}

	type createProductInput struct {
		ctx         context.Context
		Title       string
		Description string
		Price       uint
		Category    string
	}
	type expected struct {
		product *domain.Product
		err     error
	}
	repoErr := yerror.E(errors.New("error occurred in repository"))
	product := factories.Product.Create()
	product.ID = 0
	product.CreatedAt = 0

	testCases := []struct {
		name               string
		mockInsertInputs   mockInsertInputs
		mockInsertOutputs  mockInsertOutputs
		createProductInput createProductInput
		expected           expected
	}{
		{
			name: "get error in Save",
			mockInsertInputs: mockInsertInputs{
				ctx:     ctx,
				product: product,
			},
			mockInsertOutputs: mockInsertOutputs{
				product: nil,
				err:     repoErr,
			},
			createProductInput: createProductInput{
				ctx:         ctx,
				Title:       product.Title,
				Description: product.Description,
				Price:       product.Price,
				Category:    string(product.Category),
			},
			expected: expected{
				product: nil,
				err:     repoErr,
			},
		},
		{
			name: "successful test",
			mockInsertInputs: mockInsertInputs{
				ctx:     ctx,
				product: product,
			},
			mockInsertOutputs: mockInsertOutputs{
				product: &product,
				err:     nil,
			},
			createProductInput: createProductInput{
				ctx:         ctx,
				Title:       product.Title,
				Description: product.Description,
				Price:       product.Price,
				Category:    string(product.Category),
			},
			expected: expected{
				product: &product,
				err:     nil,
			},
		},
	}

	repositoryMock := new(mocks.Repository)
	aa := New(repositoryMock)

	for _, tc := range testCases {
		if tc.mockInsertOutputs.product != nil || tc.mockInsertOutputs.err != nil {
			repositoryMock.On("InsertProduct", mock.AnythingOfType("*context.timerCtx"),
				tc.mockInsertInputs.product).Return(tc.mockInsertOutputs.product, tc.mockInsertOutputs.err).Once()
		}
		got, gotErr := aa.CreateProduct(tc.createProductInput.ctx, tc.createProductInput.Title,
			tc.createProductInput.Description, tc.createProductInput.Price, tc.createProductInput.Category)
		if tc.expected.err != nil {
			assert.NotNil(t, gotErr, tc.name)
		} else {
			assert.EqualValues(t, tc.expected.product, got, tc.name)
		}
	}
	repositoryMock.AssertExpectations(t)
}
