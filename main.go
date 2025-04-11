package main

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// Configuring AWS S3 details
const (
	awsRegion    = "us-east-1"                          // Your AWS region
	bucketName   = "verifyhire-document-storage-bucket" // Your S3 bucket name
	awsAccessKey = "myaccesskey"
	awsSecretKey = "mysecretkey"
)

func main() {
	r := gin.Default()

	// POST endpoint for uploading a file to S3
	r.POST("/upload", uploadFileHandler)

	r.Run(":8080") // Run the Gin server
}

func uploadFileHandler(c *gin.Context) {
	// Get the file from the form
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}

	// Initialize an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create AWS session: " + err.Error()})
		return
	}

	// Create a new S3 service client
	s3Client := s3.New(sess)

	// Upload the file to S3
	err = uploadToS3(s3Client, file, "uploaded-file.jpg") // Choose a unique filename if needed
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload file to S3: " + err.Error()})
		return
	}

	// Success
	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

func uploadToS3(s3Client *s3.S3, file multipart.File, fileName string) error {
	// Creating the file upload parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"), // Set to "private" if you want the file to be private
	}

	// Upload the file
	_, err := s3Client.PutObject(params)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return nil
}
