package process

import "runtime"

func ExecutableName(executable string) string {
	if runtime.GOOS == "windows" {
		return executable + ".exe"
	}

	return executable
}
