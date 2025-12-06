package debug

// this package shouldn't have any engine dependencies
// because every other package should be able to use it (avoid circular dependency)
import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

var LoggingDisabled = false
var PrintLogs = true
var LogPrints = false

func Log(message ...any) {
	if LoggingDisabled {
		return
	}

	var content = "\n" + elements(message...)
	appendFile(content)

	if PrintLogs {
		fmt.Println(content)
	}
}
func LogError(message ...any) {
	if LoggingDisabled {
		return
	}

	var content = "\nERROR!\n" + callInfo(elements(message...)) + "\n"
	appendFile(content)

	if PrintLogs {
		fmt.Println(content)
	}
}

func Print(message ...any) {
	fmt.Println(elements(message...))

	if !LoggingDisabled && LogPrints {
		appendFile("\n" + elements(message...))
	}
}
func PrintLinesOfCode() {
	var directory, _ = os.Getwd()
	var cmd = exec.Command("bash", "-c", fmt.Sprintf(`find "%s" -name "*.go" -type f -exec wc -l {} +`, directory))
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdOut
	if err := cmd.Run(); err != nil {
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
		var count, _ = strconv.ParseInt(parts[0], 10, 32)
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
	fmt.Fprintf(&out, "%s\n", "Lines of code in:")

	var printTree func(path, prefix string, isLast bool)
	printTree = func(path, prefix string, isLast bool) {
		connector := "├"
		if isLast {
			connector = "└"
		}

		var name = filepath.Base(path)
		var displayCount = ""
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
			var newPrefix = prefix
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
	var cmd = exec.Command("go", "list", "-f", "{{.ImportPath}} -> {{.Imports}}", "./...")
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Run()

	var lines = strings.Split(strings.TrimSpace(cmdOut.String()), "\n")
	var deps = make(map[string][]string)

	for _, line := range lines {
		var parts = strings.Split(line, "->")
		if len(parts) != 2 {
			continue
		}
		var pkg = strings.TrimSpace(parts[0])
		var imports = strings.Fields(strings.TrimSpace(parts[1]))
		deps[pkg] = imports
	}

	var pkgs []string
	for k := range deps {
		pkgs = append(pkgs, k)
	}
	sort.Strings(pkgs)

	for _, pkg := range pkgs {
		var imports = deps[pkg]
		fmt.Fprintf(&out, "%s\n", pkg)
		sort.Strings(imports)
		for _, imp := range imports {
			imp = strings.ReplaceAll(imp, "[", "")
			imp = strings.ReplaceAll(imp, "]", "")
			fmt.Fprintf(&out, "\t%s\n", imp)
		}
		fmt.Fprintln(&out)
	}

	fmt.Print(out.String())
}
func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Basic memory usage
	fmt.Printf("\nMemory:\n")
	fmt.Printf("UsedNow = %v (current heap in use)\n", byteSize(int(m.Alloc)))
	fmt.Printf("UsedTotal = %v (total allocated since start)\n", byteSize(int(m.TotalAlloc)))
	fmt.Printf("FromOS = %v (memory reserved from OS)\n", byteSize(int(m.Sys)))

	// Heap breakdown
	fmt.Printf("\nHeap:\n")
	fmt.Printf("Used = %v \n", byteSize(int(m.HeapAlloc)))
	fmt.Printf("Reserved = %v \n", byteSize(int(m.HeapSys)))
	fmt.Printf("Idle = %v (not used but still reserved)\n", byteSize(int(m.HeapIdle)))
	fmt.Printf("Active = %v (actively in use)\n", byteSize(int(m.HeapInuse)))
	fmt.Printf("Released = %v (given back to OS)\n", byteSize(int(m.HeapReleased)))

	// Object allocations
	fmt.Printf("\nObject:\n")
	fmt.Printf("Allocs = %v (objects allocated)\n", m.Mallocs)
	fmt.Printf("Frees = %v (objects freed)\n", m.Frees)
	fmt.Printf("Live = %v (currently alive)\n", m.HeapObjects)

	// Garbage collection
	fmt.Printf("\nGarbage Collection:\n")
	fmt.Printf("Total = %v (total collections)\n", m.NumGC)
	fmt.Printf("Forced = %v (manual triggers)\n", m.NumForcedGC)
	fmt.Printf("Next = %v (target heap size of the next GC)\n", byteSize(int(m.NextGC)))
	fmt.Printf("PauseTotal = %.2f s (total time spent in GC)\n", float64(m.PauseTotalNs)/1e9)

	if m.LastGC == 0 {
		fmt.Printf("SinceLast = never\n")
	} else {
		fmt.Printf("SinceLast = %.2f s\n", time.Since(time.Unix(0, int64(m.LastGC))).Seconds())
	}

	// Stacks and other
	fmt.Printf("\nStack:\n")
	fmt.Printf("Used = %v\n", byteSize(int(m.StackInuse)))
	fmt.Printf("Reserved = %v\n", byteSize(int(m.StackSys)))
	fmt.Printf("Other = %v (misc runtime overhead)\n", byteSize(int(m.OtherSys)))
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

		exec.Command("xdg-open", svgFile).Start()
	}()
}

//=================================================================
// private

var logFile = ""

func callInfo(message string) string {
	const maxDepth = 32
	var pcs = make([]uintptr, maxDepth)
	var n = runtime.Callers(3, pcs)
	var frames = runtime.CallersFrames(pcs[:n])
	var sb strings.Builder
	sb.WriteString(message)

	for {
		var frame, more = frames.Next()
		var fileName = filepath.Base(frame.File)
		var funcName = strings.Split(frame.Function, ".")[1]

		sb.WriteString(fmt.Sprintf("\n\tat [%s] %s() { %d }", fileName, funcName, frame.Line))

		if !more || fileName == "main.go" && funcName == "main" {
			break
		}
	}

	return sb.String()
}
func appendFile(content string) {
	if logFile == "" {
		os.MkdirAll("logs", 0755)
		logFile = filepath.Join("logs", time.Now().Format("01-02-2006_15-04-05")+".txt")
		os.WriteFile(logFile, []byte(content), 0644)
		return
	}

	var file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(content)
}

// copied from utility/text
func byteSize(byteSize int) string {
	const unit = 1024
	if byteSize < unit {
		return fmt.Sprintf("%d B", byteSize)
	}
	var div, exp = int(unit), 0
	for n := byteSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.3f %cB", float32(byteSize)/float32(div), "KMGTPE"[exp])
}

// copied from utility/text.New()
func elements(elements ...any) string {
	var result = ""
	for _, e := range elements {
		switch v := e.(type) {
		case string:
			result += v
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case float32:
			result += strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			result += strconv.FormatFloat(v, 'f', -1, 64)
		case fmt.Stringer:
			result += v.String()
		default:
			var value = reflect.ValueOf(e)
			var valueType = value.Type()

			if valueType.Kind() == reflect.Struct {
				result += fmt.Sprintf("%+v", e) // struct
				continue
			}

			if valueType.Kind() == reflect.Ptr && valueType.Elem().Kind() == reflect.Struct {
				result += fmt.Sprintf("%+v", value.Elem().Interface()) // pointer to struct
				continue
			}

			result += fmt.Sprint(e) // fallback
		}
	}
	return result
}
