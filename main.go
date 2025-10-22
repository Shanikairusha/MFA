package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	showBanner()

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Separate files by year")
		fmt.Println("2. Compare internal and external lists")
		fmt.Println("3. Exit")
		fmt.Print("Enter choice (1-3): ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Enter path of input file: ")
			var inputFile string
			fmt.Scanln(&inputFile)
			separateByYear(inputFile)
		case 2:
			fmt.Print("Enter path of internal file: ")
			var internal string
			fmt.Scanln(&internal)

			fmt.Print("Enter path of external file: ")
			var external string
			fmt.Scanln(&external)

			compare(internal, external)
		case 3:
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func showBanner() {
	fmt.Println("======================================")
	fmt.Println("   File Utility CLI Tool (GoLang)     ")
	fmt.Println("   - Separate files by year           ")
	fmt.Println("   - Compare internal vs external     ")
	fmt.Println("   - Author: Shanika_s@epiclanka.net  ")
	fmt.Println("======================================")
}

func separateByYear(inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writers := make(map[string]*os.File)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) == 4 {
			path, size, year, name := parts[0], parts[1], parts[2], parts[3]
			line := fmt.Sprintf("%s|%s|%s\n", path, size, name)

			// Create/open year file if not exists
			if _, ok := writers[year]; !ok {
				outFile, err := os.Create(year + ".txt")
				if err != nil {
					panic(err)
				}
				writers[year] = outFile
			}
			writers[year].WriteString(line)
		}
	}

	for _, w := range writers {
		w.Close()
	}

	fmt.Println("✅ Files separated by year successfully!")
}

func compare(internal, external string) {
	files := map[string]string{
		"internalFile": internal,
		"externalFile": external,
	}

	for name, path := range files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("%s does not exist\n", name)
			return
		} else if err != nil {
			fmt.Printf("Error checking %s: %v\n", name, err)
			return
		}
	}

	internalMap := make(map[string]string)
	externalMap := make(map[string]string)

	// Process internal file
	intenal, err := os.Open(internal)
	if err != nil {
		panic(err)
	}
	defer intenal.Close()

	scannerA := bufio.NewScanner(intenal)
	for scannerA.Scan() {
		parts := strings.Split(scannerA.Text(), "|")
		if len(parts) < 3 {
			continue
		}
		key := parts[1] + "|" + parts[2] // a[1], a[2]
		internalMap[key] = scannerA.Text()
	}
	if err := scannerA.Err(); err != nil {
		panic(err)
	}

	// Process external file
	externalFile, err := os.Open(external)
	if err != nil {
		panic(err)
	}
	defer externalFile.Close()

	scannerB := bufio.NewScanner(externalFile)
	for scannerB.Scan() {
		parts := strings.Split(scannerB.Text(), "|")
		if len(parts) < 4 {
			continue
		}
		key := parts[1] + "|" + parts[3] // b[1], b[3]
		externalMap[key] = scannerB.Text()
	}
	if err := scannerB.Err(); err != nil {
		panic(err)
	}

	// Compare and output differences
	diffFile, _ := os.Create("External_has_Internal_Not.txt")
	defer diffFile.Close()

	diffViceFile, _ := os.Create("Internal_has_External_Not.txt")
	defer diffViceFile.Close()

	commonFile, _ := os.Create("Common_files.txt")
	defer commonFile.Close()

	// files in Internal but not in External
	fmt.Fprintln(diffViceFile, "Files present in internal but missing in external:")
	for k, v := range internalMap {
		if _, exists := externalMap[k]; !exists {
			fmt.Fprintln(diffViceFile, v)
		}
	}

	// files in External but not in Internal
	fmt.Fprintln(diffFile, "Files present in External but missing in Internal:")
	for k, v := range externalMap {
		if _, exists := internalMap[k]; !exists {
			fmt.Fprintln(diffFile, v)
		}
	}

	// files in both
	fmt.Fprintln(commonFile, "Files present in both:")
	for k := range internalMap {
		if _, exists := externalMap[k]; exists {
			fmt.Fprintln(commonFile, k)
		}
	}

	fmt.Println("✅ Comparison done! Results saved in output files.")
}
