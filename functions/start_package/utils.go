package start_package

import (
	common "ab-preview-service/common"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func copyDir(src, dst string) common.FunctionReturn {
	// Open the source directory
	srcDir, err := os.Open(src)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error opening src for copy",
		}
	}
	defer srcDir.Close()

	// Create the destination directory if it doesn't exist
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return common.FunctionReturn{
				Err:     err,
				Message: "Error creating dest",
			}
		}
	}

	// Read the contents of the source directory
	fileInfo, err := srcDir.Readdir(-1)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error reading src",
		}
	}

	// Loop through the files and directories in the source directory
	for _, file := range fileInfo {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			// Recursively copy subdirectories
			if res := copyDir(srcPath, dstPath); res.Err != nil {
				return res
			}
		} else {
			// Copy regular files
			if res := copyFile(srcPath, dstPath); res.Err != nil {
				return res
			}
		}
	}

	return common.FunctionReturn{
		Message: "successfully copied",
	}
}

func copyFile(src, dst string) common.FunctionReturn {
	srcFile, err := os.Open(src)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error opening src",
		}
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error creating dest",
		}
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error coping file",
		}
	}

	return common.FunctionReturn{
		Message: "coping file success",
	}
}

func setupEnv(sourcePath, destinationPath string) common.FunctionReturn {
	// Open the source file for reading
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error opening source file",
		}
	}
	defer sourceFile.Close()

	// Open the destination file for writing in append mode
	destinationFile, err := os.OpenFile(destinationPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		// If the destination file doesn't exist, create it
		destinationFile, err = os.Create(destinationPath)
		if err != nil {
			return common.FunctionReturn{
				Err:     err,
				Message: "Error creating destination file",
			}
		}
	}
	defer destinationFile.Close()

	// Copy the content from the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error copying file content",
		}
	}

	return common.FunctionReturn{
		Message: "Setup env successfully",
	}
}

func copyBuildOnFinish(previewName string) common.FunctionReturn {
	success := false
	maxIterations := 20 //  try for 10 mns
	iteration := 0

	for !success {
		cmd := exec.Command("docker", "compose", "logs")
		cmd.Dir = previewName

		// Capture the command's standard output and error
		output, err := cmd.CombinedOutput()
		if err != nil {
			return common.FunctionReturn{
				Err:     err,
				Message: "Error executing compose logs command",
			}
		}

		if strings.Contains(string(output), "Start process completed") || iteration >= maxIterations {
			success = true

			var err error

			containerSrcDir := "/home/ubuntu/" + previewName + "/._bb_/container_build"
			containerDstDir := "/usr/share/container_build"

			if res := copyDir(containerSrcDir, containerDstDir); err != nil {
				return res
			}

			elementsSrcDir := "/home/ubuntu/" + previewName + "/._bb_/elements_emulator/dist"
			elementsDstDir := "/usr/share/elements_build"

			if res := copyDir(elementsSrcDir, elementsDstDir); err != nil {
				return res
			}

		} else {
			fmt.Printf("\nwaiting for bb process to complete %v\n", iteration)
			fmt.Printf(string(output))

			// Sleep for 30 seconds before running the command again
			time.Sleep(30 * time.Second)

			// Increment the iteration counter
			iteration++
		}
	}

	return common.FunctionReturn{
		Message: "Copy success",
	}
}

// func readPackageName(previewName string) common.FunctionReturn {
// 	// Specify the path to your JSON file
// 	jsonFilePath := "/home/ubuntu" + previewName + "block.config.json"

// 	// Read the JSON file
// 	jsonFile, err := os.ReadFile(jsonFilePath)
// 	if err != nil {
// 		return common.FunctionReturn{
// 			Err:     err,
// 			Message: "Error reading JSON",
// 		}
// 	}

// 	// Create a variable to store the parsed JSON data
// 	var data BlockConfigData

// 	// Unmarshal the JSON data into the struct
// 	if err := json.Unmarshal(jsonFile, &data); err != nil {
// 		return common.FunctionReturn{
// 			Err:     err,
// 			Message: "Error un-marshaling JSON",
// 		}
// 	}

// 	return common.FunctionReturn{
// 		Data:    data.Name,
// 		Message: "Package name retrieved successfully",
// 	}
// }
