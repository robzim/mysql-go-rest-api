package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/kr/pretty"

	// "github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     int32   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: 1, Title: "dummy album", Artist: "Coltrane", Price: 21.99},
	{ID: 2, Title: "dummy album 2", Artist: "Cool Man", Price: 24.99},
	{ID: 3, Title: "dummy album 3", Artist: "Young", Price: 21.99},
	{ID: 4, Title: "dummy album 4", Artist: "Oldtrane", Price: 12.99},
}

func getDBConnection() sql.DB {
	var db *sql.DB

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
	return *db
}

func getAlbums(c *gin.Context) {
	albums = myGetAlbumsSQL()
	// c.IndentedJSON(http.StatusOK, albums)
	c.JSON(http.StatusOK, albums)
}

func getAlbumsFromDB(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(c *gin.Context) {
	fmt.Println("in add album")
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(newAlbum)
	albums = append(albums, newAlbum)
	myAddAlbumSQL(newAlbum.Title, newAlbum.Artist, float32(newAlbum.Price))
	fmt.Println(albums)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func myUpdateAlbumSQL(c *gin.Context) {
	theID := c.Query("id")
	theTitle := c.Query("title")
	theArtist := c.Query("artist")
	thePrice := c.Query("price")
	fmt.Println("in udpate album sql")
	fmt.Println(theID, theTitle, theArtist, thePrice)
	myDB := getDBConnection()
	myQuery := fmt.Sprintf("update album set title='%s', artist='%s', price=%v where `id` = %v", theTitle, theArtist, thePrice, theID)
	fmt.Println(myQuery)
	results, err := myDB.Query(myQuery)
	if err != nil {
		fmt.Println("failed to update")
	}
	if results != nil {
		defer results.Close()
	}
}

func myDeleteAlbumSQL(c *gin.Context) {
	theID := c.Param("id")
	fmt.Println("in delete album sql")
	fmt.Println(theID)
	myDB := getDBConnection()
	myQuery := fmt.Sprintf("delete from `album` where `id` = %v", theID)
	fmt.Println(myQuery)
	results, err := myDB.Query(myQuery)
	if err != nil {
		fmt.Println("failed to delete")
	}
	if results != nil {
		defer results.Close()
	}
}

func myAddAlbumSQL(theTitle, theAuthor string, thePrice float32) {
	fmt.Println("in addalbum")
	fmt.Println(theTitle, theAuthor, thePrice)
	myDB := getDBConnection()
	results, err := myDB.Query(fmt.Sprintf("insert into album values ('%s', '%s', %f, %v)", theTitle, theAuthor, thePrice, 0))
	if err != nil {
		fmt.Println("duplicate")
	}
	if results != nil {
		defer results.Close()
	}
}

func myExit(c *gin.Context) {
	c.Done()
	fmt.Println("All Done.")
	os.Exit(0)
}

func getAlbumByID(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	id := int32(idInt)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, a := range albums {
		fmt.Println(a.ID, id)
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"messsage": "album not found"})
}

func getAlbumByTitle(c *gin.Context) {
	title := c.Param("title")
	fmt.Println(fmt.Sprintf("title %s", title))
	var aList []album
	for _, a := range albums {
		if strings.Contains(strings.ToLower(a.Title), strings.ToLower(title)) {
			fmt.Println("appending", a, title)
			aList = append(aList, a)

		}
	}
	if len(aList) > 0 {
		c.IndentedJSON(http.StatusOK, aList)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"messsage": "album not found"})
}

func getAlbumByArtist(c *gin.Context) {
	artist := c.Param("artist")
	fmt.Println(fmt.Sprintf("title %s", artist))
	var aList []album
	for _, a := range albums {
		if strings.Contains(strings.ToLower(a.Artist), strings.ToLower(artist)) {
			fmt.Println("appending", a, artist)
			aList = append(aList, a)
		}
	}
	if len(aList) > 0 {
		// fmt.Println(c.Request.Header)
		c.IndentedJSON(http.StatusOK, aList)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"messsage": "album not found"})
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func myGetAlbumsSQL() []album {
	var theAlbums []album
	myDB := getDBConnection()

	defer myDB.Close()

	results, err := myDB.Query("select * from album")

	if err != nil {
		panic(err.Error())
	}
	for results.Next() {

		var myAlbum album

		err = results.Scan(&myAlbum.Title, &myAlbum.Artist, &myAlbum.Price, &myAlbum.ID)
		if err != nil {
			panic(err.Error())
		}
		theAlbums = append(theAlbums, myAlbum)
	}
	return theAlbums
}

func myDummyAddAlbum(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ho")
	return
}

func main() {

	defer fmt.Println("Exiting. From Defer!")

	albums = myGetAlbumsSQL()
	fmt.Println(albums)

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", getAlbums)
	router.GET("/exit", myExit)
	router.GET("/albums", getAlbums)
	router.GET("/albumid/:id", getAlbumByID)
	router.GET("/albumtitle/:title", getAlbumByTitle)
	router.GET("/albumartist/:artist", getAlbumByArtist)
	router.GET("/addalbum", addAlbum)

	router.GET("/google", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
	})

	router.POST("/addalbum", addAlbum)
	router.DELETE("/deletealbum/:id", myDeleteAlbumSQL)
	router.GET("/editalbum", myUpdateAlbumSQL)

	// var rt gin.RoutesInfo = router.Routes()

	// pretty.Print(rt)
	// // fmt.Println(rt)
	pretty.Print(albums)

	router.Run("localhost:8080")

}
