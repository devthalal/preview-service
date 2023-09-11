package pull_code

import (
	common "ab-preview-service/common"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func PullCode(previewName, s3PreSignedURL string) common.FunctionReturn {

	// Create a HTTP client
	client := &http.Client{}

	// Perform GET request
	resp, err := client.Get(s3PreSignedURL)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error pulling s3 api",
			Err:     err,
		}

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("\n resp %v\n", resp)
		return common.FunctionReturn{
			Message: "Error pulling,s3 api not success status",
			Err:     errors.New("pulling s3 api, not success"),
			Status:  resp.StatusCode,
		}
	}

	// Create a new file to save the downloaded content
	file, err := os.Create(previewName + ".zip")
	if err != nil {
		return common.FunctionReturn{
			Message: "Error creating file",
			Err:     err,
		}
	}
	defer file.Close()

	// Copy the downloaded content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return common.FunctionReturn{
			Message: "File copy error",
			Err:     err,
		}
	}

	// Create a directory to extract zip contents
	err = os.MkdirAll(previewName, 0755)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error creating extraction directory:",
			Err:     err,
		}
	}

	// Open a zip archive for reading.
	zipReader, err := zip.OpenReader(previewName + ".zip")
	if err != nil {
		return common.FunctionReturn{
			Message: "Impossible to open zip reader",
			Err:     err,
		}
	}
	defer zipReader.Close()

	// Extract zip contents
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return common.FunctionReturn{
				Message: "Error opening zip file",
				Err:     err,
			}
		}
		defer fileReader.Close()

		extractedFilePath := fmt.Sprintf("%s/%s", previewName, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			extractedFile, err := os.Create(extractedFilePath)
			if err != nil {
				return common.FunctionReturn{
					Message: "Error creating extracted file",
					Err:     err,
				}
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, fileReader)
			if err != nil {
				return common.FunctionReturn{
					Message: "Error extracting zip file",
					Err:     err,
				}
			}
		}
	}

	// Create a directory to extract zip contents
	err = os.Remove(previewName + ".zip")
	if err != nil {
		fmt.Println("Error removing directory:", err)
	}

	return common.FunctionReturn{
		Message: "Pulled successfully",
	}
}
