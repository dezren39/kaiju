package profiler

import (
	"bufio"
	"fmt"
	"kaiju/contexts"
	"kaiju/engine"
	"kaiju/klib"
	"kaiju/systems/console"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"
)

const (
	pprofCPUFile   = "cpu.prof"
	pprofHeapFile  = "heap.prof"
	pprofMergeFile = "default.pgo"
	pprofWebPort   = "9382"

	ctxDataKey        = "pprofWebCtx"
	pprofFileKey      = "pprofFile"
	pprofWebOpenedKey = "pprofWebOpened"
)

func consoleTop(host *engine.Host) string {
	cmd := exec.Command("go", "tool", "pprof", "-top", pprofCPUFile)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	out := klib.MustReturn(cmd.StdoutPipe())
	scanner := bufio.NewScanner(out)
	err := cmd.Start()
	if err != nil {
		return err.Error()
	}
	sb := strings.Builder{}
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	return sb.String()
}

func consoleMerge(host *engine.Host, argStr string) string {
	// First arg in split will be "merge" and can be skipped
	args := strings.Split(argStr, " ")[1:]
	cmdArgs := make([]string, 0, len(args)+5)
	cmdArgs = append(cmdArgs, "tool", "pprof", "-proto")
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, ">", pprofMergeFile)
	cmd := exec.Command("go", cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	err := cmd.Start()
	if err != nil {
		return err.Error()
	}
	cmd.Wait()
	return "Files merged into " + pprofMergeFile
}

func launchWeb(c *console.Console, webType string) (*contexts.Cancellable, error) {
	ctx := contexts.NewCancellable()
	targetFile := pprofCPUFile
	if webType == "mem" {
		targetFile = pprofHeapFile
	}
	cmd := exec.CommandContext(ctx, "go", "tool", "pprof", "-http=:"+pprofWebPort, targetFile)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		c.Write("Starting server on localhost:" + pprofWebPort)
		<-ctx.Done()
		cmd.Process.Kill()
		c.DeleteData(ctxDataKey)
		// Go tool pprof spawns child process pprof.exe which is not killed by the above command
		// So we need to kill it manually
		if runtime.GOOS == "windows" {
			killCmd := exec.Command("taskkill", "/F", "/IM", "pprof.exe")
			killCmd.Run()
		} else {
			// On mac, the child process is named pprof
			killCmd := exec.Command("pkill", "pprof")
			killCmd.Run()
		}
		if ctx.Err() == nil {
			c.Write("Failed to start web server, make sure you have go and graphviz installed.")
			ctx.Cancel()
		}
	}()
	return ctx, err
}

func pprofStart(c *console.Console, arg string) string {
	pprofFile := klib.MustReturn(os.Create(pprofCPUFile))
	pprof.StartCPUProfile(pprofFile)
	c.SetData(pprofFileKey, pprofFile)
	return "CPU profile started"
}

func pprofStop(c *console.Console, arg string) string {
	pprofFile, ok := c.Data(pprofFileKey).(*os.File)
	if !ok || pprofFile == nil {
		return "CPU profile not yet started"
	}
	pprof.StopCPUProfile()
	pprofFile.Close()
	return "CPU profile written to " + pprofCPUFile
}

func pprofHeap() string {
	hp := klib.MustReturn(os.Create(pprofHeapFile))
	pprof.WriteHeapProfile(hp)
	hp.Close()
	return "Heap profile written to " + pprofHeapFile
}

func pprofWebStart(c *console.Console, webType string) string {
	ctx, ok := c.Data(ctxDataKey).(*contexts.Cancellable)
	if ok && ctx != nil {
		ctx.Cancel()
		c.DeleteData(ctxDataKey)
	}
	if ctx, err := launchWeb(c, webType); err != nil {
		return err.Error()
	} else {
		if !c.HasData(ctxDataKey) {
			c.Host().OnClose.Add(func() {
				if c.HasData(ctxDataKey) {
					c.Data(ctxDataKey).(*contexts.Cancellable).Cancel()
				}
			})
			c.SetData(ctxDataKey, ctx)
		}
		return "Starting up web server..."
	}
}

func pprofWebStop(c *console.Console) string {
	ctx, ok := c.Data(ctxDataKey).(*contexts.Cancellable)
	if ok && ctx != nil {
		ctx.Cancel()
		ctx = nil
		return "Web server stopped"
	} else {
		return "Web server not running"
	}
}

func pprofWeb(c *console.Console, args []string) string {
	if len(args) < 1 {
		return `Expected "start" or "stop"`
	}
	switch args[0] {
	case "mem":
		fallthrough
	case "cpu":
		return pprofWebStart(c, args[0])
	case "stop":
		return pprofWebStop(c)
	default:
		return `Expected "cpu" or "mem" or "stop"`
	}
}

func pprofCommands(host *engine.Host, arg string) string {
	c := console.For(host)
	arg = klib.ReplaceStringRecursive(arg, "  ", " ")
	args := strings.Split(arg, " ")
	if arg == "start" {
		return pprofStart(c, arg)
	} else if arg == "stop" {
		return pprofStop(c, arg)
	} else if arg == "mem" {
		return pprofHeap()
	} else if arg == "top" {
		return consoleTop(host)
	} else if args[0] == "web" {
		return pprofWeb(c, args[1:])
	} else if strings.HasPrefix(arg, "merge") {
		return consoleMerge(host, arg)
	} else {
		return ""
	}
}

func gc(host *engine.Host, arg string) string {
	runtime.GC()
	return "Garbage collection done"
}

func memStats(host *engine.Host, arg string) string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return fmt.Sprintf("Alloc: %d, TotalAlloc: %d, Sys: %d, NumGC: %d", mem.Alloc, mem.TotalAlloc, mem.Sys, mem.NumGC)
}

func SetupConsole(host *engine.Host) {
	c := console.For(host)
	c.AddCommand("pprof", pprofCommands)
	console.For(host).AddCommand("GC", gc)
	console.For(host).AddCommand("MemStats", memStats)
}
