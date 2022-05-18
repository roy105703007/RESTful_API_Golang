package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

type FileData struct {
	Filedata string `json: "filedata"`
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
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

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

	router.POST("/file/*path", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// 上傳檔案失敗時的錯誤處理
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		}
		log.Println("file.Filename", file.Filename)
		if file == nil {
			fmt.Println("file", file)
			return
		}
		filename := filepath.Base(file.Filename)
		log.Println("filename", filename)

		// 將檔案上傳到特定位置，這裏上傳的檔案會放到 public 資料夾中
		path := c.Param("path")
		// 檢查檔案是否存在
		if _, err := os.Stat(path + "/" + filename); !os.IsNotExist(err) {
			c.String(http.StatusBadRequest, fmt.Sprintf("file is exist"))
			return
		}
		if err := c.SaveUploadedFile(file, path+"/"+filename); err != nil {
			// 存檔失敗時的錯誤處理
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
	})

	router.PATCH("/file/*path", func(c *gin.Context) {
		file, err := c.FormFile("file")

		// 上傳檔案失敗時的錯誤處理
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		}
		log.Println("file.Filename", file.Filename)
		if file == nil {
			fmt.Println("file", file)
			return
		}
		filename := filepath.Base(file.Filename)
		log.Println("filename", filename)

		// 將檔案上傳到特定位置，這裏上傳的檔案會放到 public 資料夾中
		path := c.Param("path")
		// 檢查檔案是否存在
		if _, err := os.Stat(path + "/" + filename); os.IsNotExist(err) {
			c.String(http.StatusBadRequest, fmt.Sprintf("file is not exist"))
			return
		}
		if err := c.SaveUploadedFile(file, path+"/"+filename); err != nil {
			// 存檔失敗時的錯誤處理
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
	})

	router.DELETE("/file/*path", func(c *gin.Context) {
		fmt.Println(c.FullPath())
		path := c.Param("path")
		delSuccess := delFile(path)
		if delSuccess {
			c.JSON(200, gin.H{
				"delFileSuccess": true,
				"files":          path,
			})
		} else {
			c.JSON(200, gin.H{
				"delFileSuccess": false,
				"files":          path,
			})
		}
	})

	router.Run(":8080")
}

func delFile(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	e := os.Remove(fileName)
	if e != nil {
		log.Fatal(e)
	}
	return true
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
