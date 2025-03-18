package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
    "github.com/joho/godotenv"
	"os"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	// "log"

)


// ✅ Upload Image to ImageKit.io
func UploadToImageKit(file io.Reader, fileName string) (string, error) {

	if err := godotenv.Load(); err != nil {
        return "", fmt.Errorf("failed to load .env file")
    }
	// log.Println("ENV - imgkit" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - imgkit")
    //     }
    // }
    

	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		PublicKey: os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		UrlEndpoint: os.Getenv("IMAGEKIT_URL_ENDPOINT"),
	})


	// ✅ Read file into bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// ✅ Convert file bytes to Base64 string
	fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)

	// ✅ Use `ik.Uploader.Upload()` to upload the file
	uploadRes, err := ik.Uploader.Upload(context.TODO(), fileBase64, uploader.UploadParam{
		FileName: fileName,
		Folder:   "/uploads",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to ImageKit: %v", err)
	}

	return uploadRes.Data.Url, nil // ✅ Return uploaded image URL
}
