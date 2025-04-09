package main

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

type Car struct {
	CarId        string `json:"carId"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Date         string `json:"dateOfManufacture"`
	Manufacturer string `json:"manufacturerName"`
	IPFSHash     string `json:"ipfsHash,omitempty"` // Store IPFS hash
}

const IPFSAddress = "localhost:5001" // Ensure IPFS daemon is running

var ipfs *shell.Shell

func init() {
	ipfs = shell.NewShell(IPFSAddress)
}

// Upload file to IPFS and return its hash
func UploadFileToIPFS(file multipart.File) (string, error) {
	hash, err := ipfs.Add(file)
	if err != nil {
		return "", err
	}
	fmt.Println("File uploaded to IPFS with hash:", hash)
	return hash, nil
}



func main() {
 
	router := gin.Default()

    router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

    router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Manufacturer Dashboard",
		})
	})

	router.POST("/api/car", func(ctx *gin.Context) {
		// var req Car
		// if err := ctx.BindJSON(&req); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		// 	return
		// }

		carId := ctx.DefaultPostForm("carId", "")
		carMake := ctx.DefaultPostForm("make", "")
		model := ctx.DefaultPostForm("model", "")
		color := ctx.DefaultPostForm("color", "")
		date := ctx.DefaultPostForm("dateOfManufacture", "")
		manufacturer := ctx.DefaultPostForm("manufacturerName", "")

		file, _, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}
		defer file.Close()

		 // Upload file to IPFS
		ipfsHash, err := UploadFileToIPFS(file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to IPFS"})
			return
		}
		fmt.Printf("IPFSHAsh:%v",ipfsHash)

        fmt.Printf("CarId:%v, Make:%v,Model:%v,Color:%v,Date:%v,Manufactureer:%v,File:%v",carId,carMake,model,color,date,manufacturer,file)
        // fmt.Println("request", req)

		req := Car{
			CarId:        carId,
			Make:         carMake,
			Model:        model,
			Color:        color,
			Date:         date,
			Manufacturer: manufacturer,
			IPFSHash:     ipfsHash,
		}

		result := submitTxnFn(
			"manufacturer",
			"autochannel",
			"KBA-Automobile",
			"CarContract",
			"invoke",
			make(map[string][]byte),
			"CreateCar",
			req.CarId,
			req.Make,
			req.Model,
			req.Color,
			req.Manufacturer,
			req.Date,
			req.IPFSHash,
		)
		ctx.JSON(http.StatusOK, gin.H{"message": "Created new car", "result":result})
	})

    router.GET("/api/car/:id", func(ctx *gin.Context) {
		carId := ctx.Param("id")

		result := submitTxnFn("manufacturer", "autochannel", "KBA-Automobile", "CarContract", "query", make(map[string][]byte), "ReadCar", carId)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

    router.Run(":3000")
}
