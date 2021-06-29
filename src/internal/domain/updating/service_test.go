package updating

import (
	"context"
	"errors"
	"redistore/internal/domain/factories"
	"redistore/internal/domain/ports/mocks"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"redistore/internal/domain"
	"redistore/pkg/yerror"
)

func TestNew(t *testing.T) {
	productRepository := new(mocks.Repository)
	a, ok := New(productRepository).(Service)
	assert.True(t, ok, "instance should be of type updating.Service")
	assert.NotNil(t, a, "instance should not be nil")
}
func TestAddProductToCard(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	type mockGetCardByIDInputs struct {
		ctx    context.Context
		cardID string
	}

	type mockGetCardByIDOutputs struct {
		card *domain.Card
		err  error
	}

	type mockGetProductByIDInputs struct {
		ctx       context.Context
		productID string
	}

	type mockGetProductByIDOutputs struct {
		product *domain.Product
		err     error
	}

	type mockUpdateCardInputs struct {
		ctx  context.Context
		card domain.Card
	}

	type mockUpdateCardOutputs struct {
		err error
	}
	type AddProductToCardInput struct {
		ctx       context.Context
		cardID    string
		productID string
		count     uint
	}
	type expected struct {
		err error
	}
	productId := uint(1)
	product := factories.Product.Create()
	product.ID = productId
	baseCard := factories.Card.Create()
	card := baseCard
	count := uint(1)
	updatedCard := card
	updatedCard.CardItems = make(map[string]*domain.CardItem)
	updatedCard.CardItems[strconv.FormatUint(uint64(product.ID), 10)] = &domain.CardItem{Product: &product, Count: count}
	updatedCard.Price += count * product.Price

	argsErr := yerror.E(errors.New("invalid input"))
	repoErr := yerror.E(errors.New("error occurred in repository"))

	testCases := []struct {
		name                      string
		mockGetCardByIDInputs     mockGetCardByIDInputs
		mockGetCardByIDOutputs    mockGetCardByIDOutputs
		mockGetProductByIDInputs  mockGetProductByIDInputs
		mockGetProductByIDOutputs mockGetProductByIDOutputs
		mockUpdateCardInputs      mockUpdateCardInputs
		mockUpdateCardOutputs     mockUpdateCardOutputs
		AddProductToCardInput     AddProductToCardInput
		expected                  expected
	}{
		{
			name:                      "invalid input",
			mockGetCardByIDInputs:     mockGetCardByIDInputs{},
			mockGetCardByIDOutputs:    mockGetCardByIDOutputs{},
			mockGetProductByIDInputs:  mockGetProductByIDInputs{},
			mockGetProductByIDOutputs: mockGetProductByIDOutputs{},
			mockUpdateCardInputs:      mockUpdateCardInputs{},
			mockUpdateCardOutputs:     mockUpdateCardOutputs{},
			AddProductToCardInput: AddProductToCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "",
				count:     1,
			},
			expected: expected{
				err: argsErr,
			},
		},
		{
			name: "get error in GetCardByID",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: nil,
				err:  repoErr,
			},
			mockGetProductByIDInputs:  mockGetProductByIDInputs{},
			mockGetProductByIDOutputs: mockGetProductByIDOutputs{},
			mockUpdateCardInputs:      mockUpdateCardInputs{},
			mockUpdateCardOutputs:     mockUpdateCardOutputs{},
			AddProductToCardInput: AddProductToCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
				count:     1,
			},
			expected: expected{
				err: repoErr,
			},
		},
		{
			name: "get error in GetProductByID",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: &card,
				err:  nil,
			},
			mockGetProductByIDInputs: mockGetProductByIDInputs{
				ctx:       ctx,
				productID: strconv.FormatUint(uint64(product.ID), 10),
			},
			mockGetProductByIDOutputs: mockGetProductByIDOutputs{
				product: nil,
				err:     repoErr,
			},
			mockUpdateCardInputs:  mockUpdateCardInputs{},
			mockUpdateCardOutputs: mockUpdateCardOutputs{},
			AddProductToCardInput: AddProductToCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
				count:     1,
			},
			expected: expected{
				err: repoErr,
			},
		},
		{
			name: "get error in Update",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: &card,
				err:  nil,
			},
			mockGetProductByIDInputs: mockGetProductByIDInputs{
				ctx:       ctx,
				productID: strconv.FormatUint(uint64(product.ID), 10),
			},
			mockGetProductByIDOutputs: mockGetProductByIDOutputs{
				product: &product,
				err:     nil,
			},
			mockUpdateCardInputs: mockUpdateCardInputs{
				ctx:  ctx,
				card: updatedCard,
			},
			mockUpdateCardOutputs: mockUpdateCardOutputs{
				err: repoErr,
			},
			AddProductToCardInput: AddProductToCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
				count:     1,
			},
			expected: expected{
				err: repoErr,
			},
		},
		{
			name: "successful test",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: &card,
				err:  nil,
			},
			mockGetProductByIDInputs: mockGetProductByIDInputs{
				ctx:       ctx,
				productID: strconv.FormatUint(uint64(product.ID), 10),
			},
			mockGetProductByIDOutputs: mockGetProductByIDOutputs{
				product: &product,
				err:     nil,
			},
			mockUpdateCardInputs: mockUpdateCardInputs{
				ctx:  ctx,
				card: updatedCard,
			},
			mockUpdateCardOutputs: mockUpdateCardOutputs{
				err: nil,
			},
			AddProductToCardInput: AddProductToCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
				count:     1,
			},
			expected: expected{
				err: nil,
			},
		},
	}

	repositoryMock := new(mocks.Repository)
	aa := New(repositoryMock)

	for _, tc := range testCases {
		if tc.mockGetCardByIDOutputs.card != nil || tc.mockGetCardByIDOutputs.err != nil {
			repositoryMock.On("GetCardByID", mock.AnythingOfType("*context.timerCtx"),
				tc.mockGetCardByIDInputs.cardID).Return(tc.mockGetCardByIDOutputs.card,
				tc.mockGetCardByIDOutputs.err).Once()
		}

		if tc.mockGetProductByIDOutputs.product != nil || tc.mockGetProductByIDOutputs.err != nil {
			repositoryMock.On("GetProductByID", mock.AnythingOfType("*context.timerCtx"),
				tc.mockGetProductByIDInputs.productID).Return(tc.mockGetProductByIDOutputs.product,
				tc.mockGetProductByIDOutputs.err).Once()
		}

		if tc.mockGetProductByIDOutputs.product != nil {
			repositoryMock.On("UpdateCard", mock.AnythingOfType("*context.timerCtx"),
				tc.mockUpdateCardInputs.card).Return(tc.mockUpdateCardOutputs.err).Once()
		}
		gotErr := aa.AddProductToCard(tc.AddProductToCardInput.ctx, tc.AddProductToCardInput.cardID,
			tc.AddProductToCardInput.productID, tc.AddProductToCardInput.count)
		card = baseCard
		if tc.expected.err != nil {
			assert.NotNil(t, gotErr, tc.name)
		}
	}
	repositoryMock.AssertExpectations(t)
}

func TestRemoveProductFromCard(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	type mockGetCardByIDInputs struct {
		ctx    context.Context
		cardID string
	}

	type mockGetCardByIDOutputs struct {
		card *domain.Card
		err  error
	}

	type mockUpdateCardInputs struct {
		ctx  context.Context
		card domain.Card
	}

	type mockUpdateCardOutputs struct {
		err error
	}
	type RemoveProductFromCardInput struct {
		ctx       context.Context
		cardID    string
		productID string
	}
	type expected struct {
		err error
	}
	productId := uint(1)
	product := factories.Product.Create()
	product.ID = productId
	baseCard := factories.Card.Create()
	updatedCard := baseCard
	updatedCard.CardItems = make(map[string]*domain.CardItem)

	count := uint(1)
	baseCard.CardItems = make(map[string]*domain.CardItem)
	baseCard.CardItems[strconv.FormatUint(uint64(product.ID), 10)] = &domain.CardItem{Product: &product, Count: count}
	baseCard.Price += count * product.Price

	argsErr := yerror.E(errors.New("invalid input"))
	repoErr := yerror.E(errors.New("error occurred in repository"))

	testCases := []struct {
		name                   string
		mockGetCardByIDInputs  mockGetCardByIDInputs
		mockGetCardByIDOutputs mockGetCardByIDOutputs
		mockUpdateCardInputs   mockUpdateCardInputs
		mockUpdateCardOutputs  mockUpdateCardOutputs
		AddProductToCardInput  RemoveProductFromCardInput
		expected               expected
	}{
		{
			name:                   "invalid input",
			mockGetCardByIDInputs:  mockGetCardByIDInputs{},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{},
			mockUpdateCardInputs:   mockUpdateCardInputs{},
			mockUpdateCardOutputs:  mockUpdateCardOutputs{},
			AddProductToCardInput: RemoveProductFromCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "",
			},
			expected: expected{
				err: argsErr,
			},
		},
		{
			name: "get error in GetCardByID",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: nil,
				err:  repoErr,
			},
			mockUpdateCardInputs:  mockUpdateCardInputs{},
			mockUpdateCardOutputs: mockUpdateCardOutputs{},
			AddProductToCardInput: RemoveProductFromCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
			},
			expected: expected{
				err: repoErr,
			},
		},

		{
			name: "get error in Update",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: &baseCard,
				err:  nil,
			},
			mockUpdateCardInputs: mockUpdateCardInputs{
				ctx:  ctx,
				card: updatedCard,
			},
			mockUpdateCardOutputs: mockUpdateCardOutputs{
				err: repoErr,
			},
			AddProductToCardInput: RemoveProductFromCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
			},
			expected: expected{
				err: repoErr,
			},
		},
		{
			name: "successful test",
			mockGetCardByIDInputs: mockGetCardByIDInputs{
				ctx:    ctx,
				cardID: "1",
			},
			mockGetCardByIDOutputs: mockGetCardByIDOutputs{
				card: &baseCard,
				err:  nil,
			},
			mockUpdateCardInputs: mockUpdateCardInputs{
				ctx:  ctx,
				card: updatedCard,
			},
			mockUpdateCardOutputs: mockUpdateCardOutputs{
				err: nil,
			},
			AddProductToCardInput: RemoveProductFromCardInput{
				ctx:       ctx,
				cardID:    "1",
				productID: "1",
			},
			expected: expected{
				err: nil,
			},
		},
	}

	repositoryMock := new(mocks.Repository)
	aa := New(repositoryMock)

	for _, tc := range testCases {
		if tc.mockGetCardByIDOutputs.card != nil || tc.mockGetCardByIDOutputs.err != nil {
			repositoryMock.On("GetCardByID", mock.AnythingOfType("*context.timerCtx"),
				tc.mockGetCardByIDInputs.cardID).Return(tc.mockGetCardByIDOutputs.card,
				tc.mockGetCardByIDOutputs.err).Once()
		}

		if tc.mockGetCardByIDOutputs.card != nil {
			repositoryMock.On("UpdateCard", mock.AnythingOfType("*context.timerCtx"),
				tc.mockUpdateCardInputs.card).Return(tc.mockUpdateCardOutputs.err).Once()
		}
		gotErr := aa.RemoveProductFromCard(tc.AddProductToCardInput.ctx, tc.AddProductToCardInput.cardID,
			tc.AddProductToCardInput.productID)
		updatedCard = baseCard
		if tc.expected.err != nil {
			assert.NotNil(t, gotErr, tc.name)
		}
	}
	repositoryMock.AssertExpectations(t)
}
