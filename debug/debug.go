// Unifies logging, printing and error handling consistently.
//
// Profiles the running application and analyzes the state of the project.
package debug

// this package shouldn't have any engine dependencies
// because every other package should be able to use its error logging (avoid circular dependency)
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
	"unsafe"
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

func LinesOfCode() string {
	var directory, _ = os.Getwd()
	var cmd = exec.Command("bash", "-c", fmt.Sprintf(`find "%s" -name "*.go" -type f -exec wc -l {} +`, directory))
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdOut
	if err := cmd.Run(); err != nil {
		return ""
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

	return out.String()
}
func Dependencies() string {
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

	return out.String()
}
func MemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memBuf = memBuf[:0]

	memBuf = append(memBuf, "Memory:\n"...)
	memBuf = append(memBuf, "UsedNow = "...)
	memBuf = appendByteSize(memBuf, int(m.Alloc))
	memBuf = append(memBuf, " (current heap in use)\n"...)
	memBuf = append(memBuf, "UsedTotal = "...)
	memBuf = appendByteSize(memBuf, int(m.TotalAlloc))
	memBuf = append(memBuf, " (total allocated since start)\n"...)
	memBuf = append(memBuf, "FromOS = "...)
	memBuf = appendByteSize(memBuf, int(m.Sys))
	memBuf = append(memBuf, " (memory reserved from OS)\n"...)

	memBuf = append(memBuf, "\nHeap:\n"...)
	memBuf = append(memBuf, "Used = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapAlloc))
	memBuf = append(memBuf, " \n"...)
	memBuf = append(memBuf, "Reserved = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapSys))
	memBuf = append(memBuf, " \n"...)
	memBuf = append(memBuf, "Idle = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapIdle))
	memBuf = append(memBuf, " (not used but still reserved)\n"...)
	memBuf = append(memBuf, "Active = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapInuse))
	memBuf = append(memBuf, " (actively in use)\n"...)
	memBuf = append(memBuf, "Released = "...)
	memBuf = appendByteSize(memBuf, int(m.HeapReleased))
	memBuf = append(memBuf, " (given back to OS)\n"...)

	memBuf = append(memBuf, "\nStack:\n"...)
	memBuf = append(memBuf, "Used = "...)
	memBuf = appendByteSize(memBuf, int(m.StackInuse))
	memBuf = append(memBuf, "\n"...)
	memBuf = append(memBuf, "Reserved = "...)
	memBuf = appendByteSize(memBuf, int(m.StackSys))
	memBuf = append(memBuf, "\n"...)
	memBuf = append(memBuf, "Other = "...)
	memBuf = appendByteSize(memBuf, int(m.OtherSys))
	memBuf = append(memBuf, " (misc runtime overhead)\n"...)

	memBuf = append(memBuf, "\nObjects:\n"...)
	memBuf = append(memBuf, "Allocs = "...)
	memBuf = appendThousands(memBuf, m.Mallocs)
	memBuf = append(memBuf, " (objects allocated)\n"...)
	memBuf = append(memBuf, "Frees = "...)
	memBuf = appendThousands(memBuf, m.Frees)
	memBuf = append(memBuf, " (objects freed)\n"...)
	memBuf = append(memBuf, "Live = "...)
	memBuf = appendThousands(memBuf, m.HeapObjects)
	memBuf = append(memBuf, " (currently alive)\n"...)

	memBuf = append(memBuf, "\nGarbage Collection:\n"...)
	memBuf = append(memBuf, "Total = "...)
	memBuf = appendThousands(memBuf, uint64(m.NumGC))
	memBuf = append(memBuf, " (total collections)\n"...)
	memBuf = append(memBuf, "Forced = "...)
	memBuf = strconv.AppendUint(memBuf, uint64(m.NumForcedGC), 10)
	memBuf = append(memBuf, " (manual triggers)\n"...)
	memBuf = append(memBuf, "Next = "...)
	memBuf = appendByteSize(memBuf, int(m.NextGC))
	memBuf = append(memBuf, " (target heap size of the next GC)\n"...)
	memBuf = append(memBuf, "PauseTotal = "...)
	memBuf = strconv.AppendFloat(memBuf, float64(m.PauseTotalNs)/1e9, 'f', 2, 64)
	memBuf = append(memBuf, " s (total time spent in GC)\n"...)
	if m.LastGC == 0 {
		memBuf = append(memBuf, "SinceLast = never\n"...)
	} else {
		memBuf = append(memBuf, "SinceLast = "...)
		memBuf = strconv.AppendFloat(memBuf, time.Since(time.Unix(0, int64(m.LastGC))).Seconds(), 'f', 2, 64)
		memBuf = append(memBuf, " s\n"...)
	}

	return unsafe.String(unsafe.SliceData(memBuf), len(memBuf))
}
func ProfileAllocations(seconds float32) {
	go func() {
		var ts = time.Now().Format("2006-01-02_15-04-05")
		var profileFile = fmt.Sprintf("allocs_%s.prof", ts)

		log.Printf("Allocation profiling: capturing for %.2f seconds...\n", seconds)

		var duration = time.Duration(float64(seconds) * float64(time.Second))
		time.Sleep(duration)

		runtime.GC() // flush pending frees so the snapshot is accurate

		var f, err = os.Create(profileFile)
		if err != nil {
			log.Println("could not create allocs profile:", err)
			return
		}
		defer f.Close()

		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Println("could not write allocs profile:", err)
			return
		}

		log.Println("Allocation profile saved at", profileFile)
		log.Println("Opening browser at http://localhost:8081 ...")

		exec.Command("go", "tool", "pprof", "-http=:8081", profileFile).Start()
	}()
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
		out, err := cmd.CombinedOutput() // Captures both Stdout and Stderr
		if err != nil {
			log.Printf("failed to generate svg: %v. Output: %s", err, string(out))
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
var memBuf []byte

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
func appendByteSize(buf []byte, n int) []byte {
	const unit = 1024
	if n < unit {
		buf = strconv.AppendInt(buf, int64(n), 10)
		return append(buf, " B"...)
	}
	var div, exp = int(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	buf = strconv.AppendFloat(buf, float64(n)/float64(div), 'f', 3, 64)
	buf = append(buf, ' ')
	buf = append(buf, "KMGTPE"[exp])
	return append(buf, 'B')
}
func appendThousands(buf []byte, n uint64) []byte {
	var tmp [32]byte
	var s = strconv.AppendUint(tmp[:0], n, 10)
	var length = len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, c)
	}
	return buf
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

			if value.IsNil() {
				continue
			}

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
