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

	return common.FunctionReturn{
		Message: "Success",
	}
}
