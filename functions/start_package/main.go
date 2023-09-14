package start_package

import (
	common "ab-preview-service/common"
)

func StartPackage(previewName string) common.FunctionReturn {
	var err error

	// setup env
	sourcePath := "/home/ubuntu/.env.function"
	destinationPath := "/home/ubuntu/" + previewName + "/.env.function.preview"
	res := setupEnv(sourcePath, destinationPath)
	if res.Err != nil {
		return res
	}

	viewSourcePath := "/home/ubuntu/.env.view"
	viewDestinationPath := "/home/ubuntu/" + previewName + "/.env.view.preview"
	res = setupEnv(viewSourcePath, viewDestinationPath)
	if res.Err != nil {
		return res
	}

	// run docker compose
	err = common.RunCmd(previewName, "docker", "compose", "up", "-d")

	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error docker compose up",
		}
	}

	res = copyBuildOnFinish(previewName)
	if res.Err != nil {
		return res
	}

	return common.FunctionReturn{
		Message: "Success",
	}
}
