package perfcheck

import (
	"fmt"
	"os"
	"os/exec"
)

func RunK6(benchmark *os.File, outFile string) bool {
	_, err := exec.LookPath("k6")
	if err != nil {
		fmt.Fprintln(os.Stderr, "k6 binary is not present in your path. Install it or add it to your path.")
		panic(err)
	}

	cmd := exec.Command("k6", "run", benchmark.Name(), "--out", fmt.Sprintf("json=%s", outFile))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	return err == nil
}
