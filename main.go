package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Estruturas de Dados
type MenuItem struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

type Order struct {
	gorm.Model
	TableNumber  int         `json:"table_number"`
	CustomerName string      `json:"customer_name"` 
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Status       string      `json:"status"` 
}

type OrderItem struct {
	gorm.Model
	OrderID    uint    `json:"order_id"`
	MenuItemID uint    `json:"menu_item_id"`
	Name       string  `json:"name"` 
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

var db *gorm.DB

var menuItemsSeed = []MenuItem{
	{Name: "Carne Seca com Catupiry", Description: "Carne seca desfiada, catupiry, cebola roxa e muçarela ", Price: 12.99, Category: "Pizza"},
	{Name: "Queijo Brie com Damasco", Description: "Queijo brie derretido, muçarela e geleia de damasco", Price: 8.99, Category: "Pizza"},
	{Name: "Portuguesa Light", Description: "Versão com menos queijo e ingredientes reduzidos", Price: 14.99, Category: "Pizza"},
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("file:restaurant.db?cache=shared&_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&MenuItem{}, &Order{}, &OrderItem{})

	var count int64
	db.Model(&MenuItem{}).Count(&count)
	if count == 0 {
		db.Create(&menuItemsSeed)
	}

	r := gin.Default()

    // Servir os arquivos HTML que vamos criar
    r.StaticFile("/", "./index.html")
    r.StaticFile("/cozinha", "./cozinha.html")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		api.GET("/menu", getMenu)
		api.POST("/orders", createOrder)
		api.GET("/orders/kitchen", getKitchenOrders)
		api.PUT("/orders/:id/status", updateOrderStatus)
		api.GET("/orders/history", getHistory) 
		api.DELETE("/orders/clear", clearOrders)
	}

	r.Run(":8080")
}

// Handlers
func getMenu(c *gin.Context) {
	var items []MenuItem
	db.Find(&items)
	c.JSON(200, items)
}

// Rota para deletar/limpar todos os pedidos concluídos
func clearOrders(c *gin.Context) {
	// Deleta permanentemente os pedidos com status completo
	db.Unscoped().Where("status = ?", "completed").Delete(&Order{})
	db.Unscoped().Where("status = ?", "completed").Delete(&OrderItem{})
	c.JSON(200, gin.H{"message": "Histórico limpo"})
}

// Função de criar pedido com validação obrigatória
func createOrder(c *gin.Context) {
	var input struct {
		TableNumber  int    `json:"table_number" binding:"required"`
		CustomerName string `json:"customer_name" binding:"required"`
		Items        []struct {
			MenuItemID uint `json:"menu_item_id" binding:"required"`
			Quantity   int  `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
	}

	// Se faltar nome ou mesa, o Gin retorna erro 400 automaticamente aqui
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Nome, mesa e itens são obrigatórios!"})
		return
	}

	newOrder := Order{
		TableNumber:  input.TableNumber,
		CustomerName: input.CustomerName,
		Status:       "pending",
	}
	db.Create(&newOrder)

	for _, item := range input.Items {
		var mi MenuItem
		// Busca o nome do prato no cardápio usando o ID enviado
		if err := db.First(&mi, item.MenuItemID).Error; err == nil {
			db.Create(&OrderItem{
				OrderID:    newOrder.ID,
				MenuItemID: mi.ID,
				Name:       mi.Name, // Salvando o nome no banco para a cozinha ler
				Quantity:   item.Quantity,
				Price:      mi.Price,
			})
		}
	}

	// Recarrega o pedido com os itens para confirmar
	db.Preload("Items").First(&newOrder, newOrder.ID)
	c.JSON(201, newOrder)
}

// Lista apenas o que NÃO está concluído
func getKitchenOrders(c *gin.Context) {
	var orders []Order
	// Busca ordens onde o status é diferente de 'completed'
	db.Preload("Items").Where("status != ?", "completed").Find(&orders)
	c.JSON(200, orders)
}

// Nova rota para ver o histórico de concluídos
func getHistory(c *gin.Context) {
	var orders []Order
	db.Preload("Items").Where("status = ?", "completed").Find(&orders)
	c.JSON(200, orders)
}

func updateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	// Atualiza o banco de dados
	result := db.Model(&Order{}).Where("id = ?", id).Update("status", body.Status)
	
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Falha ao atualizar banco"})
		return
	}

	c.JSON(200, gin.H{"message": "Status atualizado para " + body.Status})
}