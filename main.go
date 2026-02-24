package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Handler     func(args []string)
}

func main() {
	printBanner()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("  luna> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		if cmd == "exit" || cmd == "quit" {
			fmt.Println()
			fmt.Println("  Goodbye!")
			break
		}

		handleCommand(cmd, args)
	}
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "help", "?":
		printHelp()
	case "init":
		cmdInit(args)
	case "ingest":
		cmdIngest(args)
	case "compute":
		cmdCompute(args)
	case "reconcile":
		cmdReconcile(args)
	case "generate":
		cmdGenerate(args)
	case "status":
		cmdStatus(args)
	default:
		fmt.Printf("  Unknown command: %s\n", cmd)
		fmt.Println("  Type 'help' for available commands")
		fmt.Println()
	}
}

func cmdInit(args []string) {
	if len(args) < 1 {
		fmt.Println("  Usage: init <client-name>")
		fmt.Println()
		return
	}
	client := args[0]
	fmt.Printf("  Creating workspace for: %s\n", client)
	fmt.Println("  (Workspace creation coming soon)")
	fmt.Println()
}

func cmdIngest(args []string) {
	if len(args) < 1 {
		fmt.Println("  Usage: ingest <file-path>")
		fmt.Println()
		return
	}
	file := args[0]
	fmt.Printf("  Parsing: %s\n", file)
	fmt.Println("  (Document parsing coming soon)")
	fmt.Println()
}

func cmdCompute(args []string) {
	if len(args) < 1 {
		fmt.Println("  Usage: compute <tax|tds|gst>")
		fmt.Println()
		return
	}
	computeType := args[0]
	fmt.Printf("  Computing: %s\n", computeType)
	fmt.Println("  (Computation engine coming soon)")
	fmt.Println()
}

func cmdReconcile(args []string) {
	if len(args) < 1 {
		fmt.Println("  Usage: reconcile <26as|gstr>")
		fmt.Println()
		return
	}
	reconcileType := args[0]
	fmt.Printf("  Reconciling: %s\n", reconcileType)
	fmt.Println("  (Reconciliation coming soon)")
	fmt.Println()
}

func cmdGenerate(args []string) {
	if len(args) < 1 {
		fmt.Println("  Usage: generate <itr|gstr>")
		fmt.Println()
		return
	}
	generateType := args[0]
	fmt.Printf("  Generating: %s\n", generateType)
	fmt.Println("  (JSON generation coming soon)")
	fmt.Println()
}

func cmdStatus(args []string) {
	fmt.Println("  Client Status")
	fmt.Println("  ─────────────────────────────────────")
	fmt.Println("  No client workspace found")
	fmt.Println("  Run 'init <client-name>' first")
	fmt.Println()
}

func printHelp() {
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("  ─────────────────────────────────────")
	fmt.Println("  init <client>     Create new client workspace")
	fmt.Println("  ingest <file>    Parse and ingest document")
	fmt.Println("  compute <type>   Compute tax/TDS/GST")
	fmt.Println("  reconcile <type> Reconcile 26AS/GSTR")
	fmt.Println("  generate <type>  Generate ITR/GSTR JSON")
	fmt.Println("  status           Show client status")
	fmt.Println("  help             Show this help")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("  ─────────────────────────────────────")
	fmt.Println("  luna init rajesh-sharma")
	fmt.Println("  luna ingest form16.pdf")
	fmt.Println("  luna compute tax")
	fmt.Println("  luna reconcile 26as")
	fmt.Println("  luna generate itr")
	fmt.Println()
}

func printBanner() {
	fmt.Println()
	fmt.Println("  ██╗     ██╗   ██╗███╗   ██╗ █████╗ ")
	fmt.Println("  ██║     ██║   ██║████╗  ██║██╔══██╗")
	fmt.Println("  ██║     ██║   ██║██╔██╗ ██║███████║")
	fmt.Println("  ██║     ██║   ██║██║╚██╗██║██╔══██║")
	fmt.Println("  ███████╗╚██████╔╝██║ ╚████║██║  ██║")
	fmt.Println("  ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝")
	fmt.Println()
	fmt.Println("  AI CA Agent - Type 'help' for commands, 'exit' to quit")
	fmt.Println("  ─────────────────────────────────────────────────────")
	fmt.Println()
}
