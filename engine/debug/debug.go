package debug

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"pure-kit/engine/data/folder"
	"pure-kit/engine/utility/text"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
)

func PrintLinesOfCode() {
	var directory = folder.PathOfExecutable()
	var cmd = exec.Command("bash", "-c",
		fmt.Sprintf(`find "%s" -name "*.go" -type f -exec wc -l {} +`, directory),
	)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdOut
	var err = cmd.Run()
	if err != nil {
		return
	}

	var results = make(map[string]int)
	var scanner = bufio.NewScanner(&cmdOut)

	for scanner.Scan() {
		var line = strings.TrimSpace(scanner.Text())
		var parts = strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		var count = text.FromNumber(parts[0])
		var path = parts[1]
		var rel, _ = filepath.Rel(directory, path)
		results[rel] = int(count)
	}

	var dirTotals = make(map[string]int)
	for path, count := range results {
		dirTotals[path] = count
		var dir = filepath.Dir(path)
		for dir != path {
			dirTotals[dir] += count
			dir = filepath.Dir(dir)
			if dir == "." {
				break
			}
		}
	}

	var allPaths []string
	for p := range dirTotals {
		if p != "." {
			allPaths = append(allPaths, p)
		}
	}
	sort.Strings(allPaths)

	var out strings.Builder
	fmt.Fprintf(&out, "%s", "Lines of code in:\n")

	var printTree func(path string, prefix string, isLast bool)
	printTree = func(path string, prefix string, isLast bool) {
		var connector string
		if isLast {
			connector = "└"
		} else {
			connector = "├"
		}

		var name = filepath.Base(path)

		var displayCount string
		if _, ok := results[path]; ok { // It's a file
			displayCount = fmt.Sprintf("%d", dirTotals[path])
		} else { // It's a directory
			displayCount = fmt.Sprintf("[%d]", dirTotals[path])
		}

		if name == "." {
			fmt.Fprintf(&out, "[%s] %s\n", displayCount, directory)
		} else {
			fmt.Fprintf(&out, "%6s %s%s%s\n", displayCount, prefix, connector, name)
		}

		var children []string
		for _, p := range allPaths {
			if filepath.Dir(p) == path {
				children = append(children, p)
			}
		}
		sort.Strings(children)

		for i, c := range children {
			var newPrefix string
			if isLast {
				newPrefix = prefix + "  "
			} else {
				newPrefix = prefix + "│ "
			}
			printTree(c, newPrefix, i == len(children)-1)
		}
	}

	var topLevel []string
	for _, p := range allPaths {
		if !strings.Contains(p, string(filepath.Separator)) {
			topLevel = append(topLevel, p)
		}
	}
	sort.Strings(topLevel)
	for i, t := range topLevel {
		printTree(t, "", i == len(topLevel)-1)
	}

	fmt.Printf("%v\n", out.String())
}

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
