package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

		if input == "exit" || input == "quit" {
			fmt.Println()
			fmt.Println("  Goodbye!")
			break
		}

		if input == "help" {
			printHelp()
			continue
		}

		fmt.Println()
		fmt.Printf("  You said: %q\n\n", input)
	}
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

func printHelp() {
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("  ─────────────────────────────────────")
	fmt.Println("  init <client>     Create new client workspace")
	fmt.Println("  ingest <file>    Parse and ingest document")
	fmt.Println("  compute tax      Compute tax liability")
	fmt.Println("  reconcile 26as  Match 26AS with books")
	fmt.Println("  generate itr     Generate ITR JSON")
	fmt.Println("  help            Show this help")
	fmt.Println("  exit            Exit Luna")
	fmt.Println()
}
