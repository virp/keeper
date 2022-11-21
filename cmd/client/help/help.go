package help

import (
	"fmt"
)

var (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

// Usage show client usage help text.
func Usage() {
	fmt.Println(Yellow + "Keeper usage: " + Cyan + "keeper [options] <command> [command args]" + Reset)
	fmt.Println(Yellow + "Commands:" + Reset)
	fmt.Println("\t" + Green + "register" + Reset + " - register new account by login and password")
	fmt.Println("\t" + Green + "login" + Reset + "    - login via existed login and password")
	fmt.Println("\t" + Green + "ls" + Reset + "       - list items")
	fmt.Println("\t" + Green + "get" + Reset + "      - get item details")
	fmt.Println("\t" + Green + "add" + Reset + "      - create item")
	fmt.Println("\t" + Green + "edit" + Reset + "     - edit item")
	fmt.Println("\t" + Green + "del" + Reset + "      - remove item")
	fmt.Println("\tGet command help: " + Cyan + "keeper <command> help" + Reset)
}

// RegisterUsage show "register" command usage help text.
func RegisterUsage() {
	fmt.Println(Yellow + "Register account: " + Cyan + "keeper [options] register" + Reset)
}

// LoginUsage show "login" command usage help text.
func LoginUsage() {
	fmt.Println(Yellow + "Login: " + Cyan + "keeper [options] login" + Reset)
}

// LsUsage show "ls" (list) command usage help text.
func LsUsage() {
	fmt.Println(Yellow + "List items: " + Cyan + "keeper [options] ls" + Reset)
}

// GetUsage show "get" command usage help text.
func GetUsage() {
	fmt.Println(Yellow + "Get item details: " + Cyan + "keeper [options] get <item name>" + Reset)
}

// AddUsage show "add" command usage help text.
func AddUsage() {
	fmt.Println(Yellow + "Add item: " + Cyan + "keeper [options] add <item name> <item type: password, text, card, binary>" + Reset)
}

// EditUsage show "edit" command usage help text.
func EditUsage() {
	fmt.Println(Yellow + "Edit item: " + Cyan + "keeper [options] edit <item name>" + Reset)
}

// DelUsage show "del" (delete) command usage help text.
func DelUsage() {
	fmt.Println(Yellow + "Remote item: " + Cyan + "keeper [options] del <item name>" + Reset)
}
