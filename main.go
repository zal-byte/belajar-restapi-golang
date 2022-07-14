package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
	"github.com/google/uuid"
)



type Customer struct {
	ID       string	`json:"ID" validate:"nonzero"`
	Fullname string `json:"Fullname" validate:"min=5"`
	Address  string `json:"Address" validate:"nonzero"`
	Email    string `json:"Email" validate:"nonzero"`
	Phone    uint16 `json:"Phone" validate:"nonzero"`
	Age      uint8  `json:"Age" validate:"min=0, max=80"`
}

var customers = []Customer{}



func newCustomer(c *gin.Context) {


	sampleUUID := uuid.New()
	fmt.Println(sampleUUID.String())

	var newCustomer Customer


	err := c.ShouldBindJSON( &newCustomer )
	if errs := validator.Validate(&newCustomer); errs != nil{
		c.JSON( http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	if err != nil{
		fmt.Println(err)
		return
	}


	for _, value := range customers{
		if value.ID == newCustomer.ID{
			c.JSON( http.StatusConflict, gin.H{"message": "This ID has been used by another customer !"})
			return
		}
	}



	customers = append(customers, newCustomer)

	c.JSON(http.StatusCreated, gin.H{"message": "New customer has been added"})
}

func getCustomerById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range customers {
		if a.ID == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})

}

func showCustomer(c *gin.Context) {
	c.JSON(http.StatusOK, customers)
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")

	var updateCustomer Customer

	if err := c.ShouldBindJSON(&updateCustomer); err != nil {
		fmt.Println(err)
		return
	}

	for i, a := range customers {
		if a.ID == id {
			updateCustomer.ID = a.ID
			customers[i] = updateCustomer

			customers = append(customers, customers[i])
			c.JSON( http.StatusOK, updateCustomer )
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})

}

func delete(array []Customer, index int) []Customer {
	return append(array[:index], array[index+1:]...)
}

func deleteCustomer(c *gin.Context) {

	id := c.Param("id")

	for index, value := range customers{
		if value.ID == id{
			customers = delete(customers, index)
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
}

func main() {
	fmt.Println("TEST")

	router := gin.Default()

	router.GET("/customers", showCustomer)
	router.GET("/customers/:id", getCustomerById)
	router.POST("/customers", newCustomer)
	router.DELETE("/customers/:id", deleteCustomer)
	router.PUT("/customers/:id", updateCustomer)

	router.Run("localhost:8080")
}
