package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type book struct{
	ID		string `json:"id"`
	Title	string `json:"title"`
	Author	string `json:"author"`
	Quantity  int  `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context){
	c.IndentedJSON(http.StatusOK,books) //get a JSON with a proper Indentation
}//return the json version of the books (GET ENDPOINT)

func createBook(c *gin.Context){
	var newBook book

	if err := c.BindJSON(&newBook); err!=nil{ 
		return //sending the error response
	}

	books = append(books,newBook)
	c.IndentedJSON(http.StatusCreated,newBook)

}//get the data of the book that we're going to be creating

func bookByID(c *gin.Context){
	id := c.Param("id") //path parameter
	book,err := getBookByID(id)

	if err !=nil { //without using BindJSON we have to specify our error message
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK,book)
}

func getBookByID(id string)(*book, error){
	for i,b:= range books{
		if b.ID == id{
			return &books[i],nil
		} // returning nil because we found no error
	}

	return nil, errors.New("Book not found...")
}

func checkoutBook(c *gin.Context){
	id,ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Missing id query parameter..."})
		return
	}//do we have that id?

	book,err := getBookByID(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}// checking the book's existing

	if book.Quantity <=0 {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Book not available..."})
		return
	} //checking the quantity of the book

	book.Quantity-=1
	c.IndentedJSON(http.StatusOK,gin.H{"message":"Book checked out successfully!"})
}

func returnBook(c *gin.Context){
	id,ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message":"Missing id query parameter..."})
		return
	}//do we have that id?

	book,err := getBookByID(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}// checking the book's existing

	//update the quality
	book.Quantity +=1
	c.IndentedJSON(http.StatusOK,gin.H{"message": "Book returned successfully."})
}


func main(){
	router := gin.Default()
	router.GET("/books",getBooks) 
	router.GET("/books/:id",bookByID)
	router.POST("/books",createBook) //add
	router.PATCH("/checkout",checkoutBook) //update
	router.PATCH("/return",returnBook)
	router.Run("localhost:8080") //running our web server
}