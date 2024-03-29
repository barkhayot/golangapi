package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	

	"strings"

	"goapi/testing/models"
	"goapi/testing/routes"
	

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

//func homePage(c *gin.Context){
//	c.JSON(200, gin.H{"message": "Hello! Please use Postman to Test to test this API"})
//}

func main() {

	conn, err := connectDB()
	if err != nil {
		return
	}

	router := gin.New()

	router.Use(dbMiddleware(*conn))

	router.GET("home", func(c *gin.Context){
		c.JSON(200, gin.H{"message": "Hello. Welcome to test golang API",
							"please" : "Please Use The Postman to test this API"})
	})

	usersGroup := router.Group("auth")
	{
		usersGroup.POST("signup", routes.UsersRegister)
		usersGroup.POST("signin", routes.UsersLogin)
	}

	itemsGroup := router.Group("items")
	{
		itemsGroup.GET("index", authMiddleWare(), routes.ItemsIndex)
		itemsGroup.POST("create", authMiddleWare(), routes.ItemsCreate)
		itemsGroup.GET("sold_by_user", authMiddleWare(), routes.ItemsForSaleByCurrentUser)
		itemsGroup.PUT("update", authMiddleWare(), routes.ItemsUpdate)
	}

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, nil))
	//router.Run(":3000")
}

func connectDB() (c *pgx.Conn, err error) {
	conn, err := pgx.Connect(context.Background(), "Database://databaseuser:databasehost@password:5432(port)/databasename")
	if err != nil || conn == nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")
		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
			return
		}
		token := split[1]
		//fmt.Printf("Bearer (%v) \n", token)
		isValid, userID := models.IsTokenValid(token)
		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated."})
			c.Abort()
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
