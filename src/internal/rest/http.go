package rest

import (
	"github.com/gin-gonic/gin"
	"redistore/internal/domain/creating"
	"redistore/internal/domain/listing"
	"redistore/internal/domain/searching"
	"redistore/internal/domain/updating"
)

type HTTPHandler struct {
	creatingService  creating.Service
	updatingService  updating.Service
	searchingService searching.Service
	listingService   listing.Service
}

func New(creatingService creating.Service, updatingService updating.Service, searchingService searching.Service, listingService listing.Service) *HTTPHandler {
	return &HTTPHandler{
		creatingService:  creatingService,
		searchingService: searchingService,
		updatingService:  updatingService,
		listingService:   listingService,
	}
}

func (hdl *HTTPHandler) CreateProduct(c *gin.Context) {
	body := ProductCreateDTO{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	product, err := hdl.creatingService.CreateProduct(c, body.Title, body.Description, body.Price, body.Category)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (hdl *HTTPHandler) CreateCard(c *gin.Context) {
	body := CardCreateDTO{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	card, err := hdl.creatingService.CreateCard(c, body.UserID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, card)
}

func (hdl *HTTPHandler) AddProductToCard(c *gin.Context) {
	body := AddProductToCardDTO{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	err = hdl.updatingService.AddProductToCard(c, body.CardID, body.ProductID, body.Count)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "done!"})
}

func (hdl *HTTPHandler) RemoveCardItem(c *gin.Context) {
	body := RemoveCardItemDTO{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	err = hdl.updatingService.RemoveProductFromCard(c, body.CardID, body.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "done!"})
}

func (hdl *HTTPHandler) SearchProductsByTitle(c *gin.Context) {
	body := SearchProductDTO{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	products, err := hdl.searchingService.SearchProductsByTitle(c, body.Title)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, products)

}

func (hdl *HTTPHandler) GetProductList(c *gin.Context) {

	products, err := hdl.listingService.GetProductList(c)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, products)

}
