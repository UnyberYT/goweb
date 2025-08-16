package main

import (

	"strconv"
	"net/http"
	"gorm.io/driver/sqlite"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

var db *gorm.DB 
func main() {
	var err error
	//Connect to Sqlite
	db, err = gorm.Open(sqlite.Open("todos.db"),&gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//Auto create table if not exists
	db.AutoMigrate(&Todo{})
	// Seed data if table is empty
	var count int64
	db.Model(&Todo{}).Count(&count)
	if count == 0 {
		db.Create(&Todo{Title: "Learn Go", Completed: false})
		db.Create(&Todo{Title: "Build a Todo API", Completed: false})
		db.Create(&Todo{Title: "Master GORM", Completed: false})
	}


	r := gin.Default()



	//Get all todos
	r.GET("/todos", func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		c.JSON(http.StatusOK, todos)
	})

	//Get a single todo by ID
	r.GET("/todos/:id",func(c *gin.Context) {
		id, _ :=strconv.Atoi(c.Param("id"))
		var todo Todo
		if err := db.First(&todo, id).Error; err !=nil {
			c.JSON(http.StatusNotFound, gin.H{"error" : "Todo not found"})
			return
		}
		c.JSON(http.StatusOK, todo)
			})
		


	//Create a new todo
	r.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&newTodo)
		c.JSON(http.StatusCreated, newTodo)
	})

	//Update a todo

r.PUT("/todos/:id", func(c *gin.Context) {
	id, _ :=strconv.Atoi( c.Param("id"))
	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.ID = uint(id)
	db.Save(&todo)
	c.JSON(http.StatusOK, todo)
	
})
	
	//Delete a todo

r.DELETE("/todos/:id", func(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db.Delete(&Todo{}, id)
	c.JSON(http.StatusOK, gin.H{"message" : "Deleted"})
})

r.Run((":8082"))
}
