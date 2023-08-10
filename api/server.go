package api

import (
	db "github.com/alvinahb/supply-me/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves all HTTP requests for our application
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/company", server.createCompany)
	router.GET("/company/:id", server.getCompany)
	router.GET("/companies", server.listCompanies)

	router.POST("/product", server.createProduct)
	router.GET("/product/:id", server.getProduct)
	router.GET("/products", server.listProducts)

	router.POST("/inventory", server.createInventory)
	router.GET("/inventory/:company_id/:product_id", server.getCompanyProductInventory)
	router.GET("/inventories", server.listCompanyInventories)

	router.POST("/order", server.createOrder)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
