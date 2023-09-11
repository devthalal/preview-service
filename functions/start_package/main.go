package start_package

import (
	common "ab-preview-service/common"
)

func StartPackage(previewName string, unusedPorts []int) common.FunctionReturn {
	var err error

	// Remove once appblock docker is released
	err = common.RunCmd(previewName, "docker", "build", ".", "-t", "appblocks_1.0.0_nodejs")

	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error building appblock base image",
		}
	}

	// run docker compose
	err = common.RunCmd(previewName, "docker", "compose", "up", "-d")

	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error docker compose up",
		}
	}

	return common.FunctionReturn{
		Message: "Success",
	}
}

// Remove
// err = setAndPersistEnvVariable("CLOUDFLARE_API_TOKEN", os.Getenv("CLOUDFLARE_API_TOKEN"))
// if err != nil {
// 	fmt.Println("Error:", err)
// 	return common.FunctionReturn{
// 		Err:     err,
// 		Message: "Error cloudflare api token setup source",
// 	}
// }
// func setAndPersistEnvVariable(variableName, value string) error {
// 	// Set the environment variable within Go
// 	// os.Setenv(variableName, value)

// 	// Append to .bashrc
// 	appendCmd := fmt.Sprintf(`echo 'export %s="%s"' >> ~/.bashrc`, variableName, value)
// 	cmd := exec.Command("bash", "-c", appendCmd)
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("error appending to .bashrc: %s", err)
// 	}

// 	cmdSource := exec.Command("bash", "-c", "source ~/.bashrc")
// 	errSource := cmdSource.Run()
// 	if errSource != nil {
// 		return fmt.Errorf("error appending to .bashrc: %s", errSource)
// 	}

// 	return nil
// }
