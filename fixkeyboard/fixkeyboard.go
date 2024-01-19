package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const keymapFile = "config/adv360.keymap"

func warnOnForbiddenKeycodes(lines []string) []string {
	// I've had this issue where there are some keys that can't be bound (or the symbol is wrong)
	// and using them causes the keyboard to become non-functional.
	return lines
}

func readLines() ([]string, error) {
	keymapFile, err := os.Open(keymapFile)

	if err != nil {
		return nil, err
	}

	defer keymapFile.Close()

	scanner := bufio.NewScanner(keymapFile)
	var srcLines []string

	for scanner.Scan() {
		srcLines = append(srcLines, scanner.Text())
	}

	return srcLines, nil
}

func writeLines(newLines []string) error {
	keymapFile, err := os.Create(keymapFile)

	if err != nil {
		return err
	}

	defer keymapFile.Close()

	for _, line := range newLines {
		if _, err := keymapFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := func() error {
		// the file dumped out from the GUI configurator always obliterates my &hm mods.
		// also, the formatting makes it tough to see what I've changed.

		srcLines, err := readLines()

		if err != nil {
			return err
		}

		// I have a bunch of &hm mods that get rewritten by the
		// Kinesis GUI configurator as &mt's instead, which isn't
		// great for homerow stuff. So we just fix them here.
		homerowModsKeys := []string{
			"LGUI", "LSHIFT", "LALT", "LCTRL",
			"RGUI", "RSHIFT", "RALT", "RCTRL",
		}

		// find the line that says `            bindings = <`
		var destLines []string

		inBindings := false

		for _, line := range srcLines {
			if strings.Contains(line, "bindings = <") {
				// we've found our bindings
				destLines = append(destLines, line)
				inBindings = true
				continue
			}

			if strings.Contains(line, ">;") {
				inBindings = false
				destLines = append(destLines, line)
				continue
			}

			if inBindings && strings.Contains(line, "&") {
				for _, key := range homerowModsKeys {
					line = strings.Replace(line,
						fmt.Sprintf("&mt %s", key),
						fmt.Sprintf("&hm %s", key), -1)
				}

				// now introduce newlines between the keys
				//line = strings.ReplaceAll(line, "&", "\n&")
				destLines = append(destLines, line)
				continue
			}

			destLines = append(destLines, line)
		}

		return writeLines(destLines)
	}()

	if err != nil {
		panic(err)
	}
}
