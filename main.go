package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	//"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type IndexData struct {
	Title   string
	Content string
}

type FileInfo interface {
	Name() string // base name of the file
	Size() int64  // length in bytes for regular files; system-dependent for others
	//Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

type Order int64

const (
	Undefined Order = iota
	LastModified
	Size
	FileName
)

func main() {
	router := gin.Default()

	// 注册路由和Handler
	// url为 /welcome?firstname=Jane&lastname=Doe
	router.GET("/file", func(c *gin.Context) {
		orderBy_input := c.DefaultQuery("orderBy", "Undefined")
		orderDirection := c.DefaultQuery("orderByDirection", "Undefined")
		filter := c.DefaultQuery("filterByName", "")
		orderBy := Undefined
		if orderBy_input == "LastModified" {
			orderBy = LastModified
		} else if orderBy_input == "Size" {
			orderBy = Size
		} else if orderBy_input == "FileName" {
			orderBy = FileName
		} else {
			orderBy = Undefined
		}
		var files []string
		files = getAllFile(orderBy, orderDirection, filter)
		//delFile(filename)
		//c.String(http.StatusOK, "Hello %s %s", filename)
		if len(files) == 0 {
			c.JSON(404, gin.H{
				"message": "HTTP code not found",
			})
		}
		c.JSON(200, gin.H{
			"isDirectory": true,
			"files":       files,
		})
	})

	router.Run(":8080")
}

func newFile(fileName string) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte("要寫入的文本內容"))
		fmt.Print(err)
	}
}

func delFile(fileName string) {
	e := os.Remove(fileName)
	if e != nil {
		log.Fatal(e)
	}
}

func getAllFile(orderBy Order, orderDirection string, filter string) []string {
	myfolder := "./"
	var all []string
	files, _ := ioutil.ReadDir(myfolder)
	switch orderBy {
	case LastModified:
		sort.Slice(files, func(i, j int) bool {
			if orderDirection == "Descending" {
				return files[i].ModTime().Unix() > files[j].ModTime().Unix()
			}
			return files[i].ModTime().Unix() < files[j].ModTime().Unix()
		})
	case Size:
		sort.Slice(files, func(i, j int) bool {
			if orderDirection == "Descending" {
				return files[i].Size() > files[j].Size()
			}
			return files[i].Size() < files[j].Size()
		})
	case FileName:
		sort.Slice(files, func(i, j int) bool {
			if orderDirection == "Descending" {
				return files[i].Name() > files[j].Name()
			}
			return files[i].Name() < files[j].Name()
		})
	case Undefined:
	}
	if filter != "" {
		for _, file := range files {
			if strings.Contains(file.Name(), filter) {
				fmt.Println(file.Name())
				all = append(all, file.Name())
			}
		}
		return all
	}
	for _, file := range files {
		fmt.Println(file.Name())
		all = append(all, file.Name())
		// if file.IsDir() {
		// 	continue
		// } else {
		// 	fmt.Println(file.Name())
		// 	all = append(all, file.Name())
		// }
	}
	return all
}
