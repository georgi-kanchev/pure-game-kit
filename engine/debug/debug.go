package debug

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"pure-kit/engine/data/path"
	"pure-kit/engine/utility/text"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
)

func PrintLinesOfCode() {
	directory := path.Folder(path.Executable())
	cmd := exec.Command("bash", "-c",
		fmt.Sprintf(`find "%s" -name "*.go" -type f -exec wc -l {} +`, directory),
	)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdOut
	if err := cmd.Run(); err != nil {
		return
	}

	results := make(map[string]int)
	scanner := bufio.NewScanner(&cmdOut)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		count := text.ToNumber(parts[0])
		path := parts[1]
		rel, _ := filepath.Rel(directory, path)
		results[rel] = int(count)
	}

	dirTotals := make(map[string]int)
	for path, count := range results {
		dirTotals[path] = count
		dir := filepath.Dir(path)
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
	fmt.Fprintf(&out, "%s\n", "Lines of code in:")

	var printTree func(path, prefix string, isLast bool)
	printTree = func(path, prefix string, isLast bool) {
		connector := "├"
		if isLast {
			connector = "└"
		}

		name := filepath.Base(path)
		displayCount := ""
		if _, ok := results[path]; ok {
			displayCount = fmt.Sprintf("%d", dirTotals[path])
		} else {
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
			newPrefix := prefix
			if isLast {
				newPrefix += "  "
			} else {
				newPrefix += "│ "
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

	fmt.Print(out.String())
}

func PrintDependencies() {
	var out strings.Builder
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}} -> {{.Imports}}", "./...")
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Run()

	lines := strings.Split(strings.TrimSpace(cmdOut.String()), "\n")
	deps := make(map[string][]string)

	for _, line := range lines {
		parts := strings.Split(line, "->")
		if len(parts) != 2 {
			continue
		}
		pkg := strings.TrimSpace(parts[0])
		imports := strings.Fields(strings.TrimSpace(parts[1]))
		deps[pkg] = imports
	}

	var pkgs []string
	for k := range deps {
		pkgs = append(pkgs, k)
	}
	sort.Strings(pkgs)

	for _, pkg := range pkgs {
		fmt.Fprintf(&out, "%s\n", pkg)
		imports := deps[pkg]
		sort.Strings(imports)
		for _, imp := range imports {
			imp = text.Remove(imp, "[", "]")
			fmt.Fprintf(&out, "\t%s\n", imp)
		}
		fmt.Fprintln(&out)
	}

	fmt.Print(out.String())
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
