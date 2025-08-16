package performance

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"time"
)

func ProfileCPU(seconds float32) {
	go func() {
		// timestamp for filenames
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var profileFile = fmt.Sprintf("cpu_%s.prof", ts)
		var svgFile = fmt.Sprintf("cpu_%s.svg", ts)

		var f, err = os.Create(profileFile)
		if err != nil {
			log.Println("could not create profile:", err)
			return
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Println("could not start CPU profile:", err)
			return
		}
		log.Printf("CPU profiling started for %.2f seconds...\n", seconds)

		// convert float32 seconds → duration
		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		pprof.StopCPUProfile()
		log.Println("CPU profiling stopped. Profile saved at", profileFile)

		// Generate SVG via `go tool pprof`
		var cmd = exec.Command("go", "tool", "pprof", "-svg", profileFile)
		out, err := cmd.Output()
		if err != nil {
			log.Println("failed to generate svg:", err)
			return
		}

		if err := os.WriteFile(svgFile, out, 0644); err != nil {
			log.Println("failed to save svg:", err)
			return
		}

		log.Println("SVG generated at", svgFile)

		var open *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			open = exec.Command("rundll32", "url.dll,FileProtocolHandler", svgFile)
		case "linux":
			open = exec.Command("xdg-open", svgFile)
		}

		open.Start()
	}()
}
