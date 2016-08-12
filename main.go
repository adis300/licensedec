package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultExts = "m,h,js,swift,go,cpp,mm,hh,hpp,java"

var templatePath, rootPath string
var extensions = map[string]bool{}
var isRecursive, shouldCleanBackup bool
var fileProcessRecord = map[string]bool{}

var licenseContent []byte

func main() {
	// Define an parse arguments
	templatePathPtr := flag.String("template", "LICENSE", "Path to the license template")
	rootPathPtr := flag.String("root", "./", "Absolute path to find all file to modify")
	extensionsPtr := flag.String("exts", defaultExts, "All extensions to modify, comma separated values")
	isRecursivePtr := flag.Bool("recursive", false, "Marks if the program should search subfolders")
	shouldCleanBackupPtr := flag.Bool("clean", true, "Marks if the program should clean back up files")

	flag.Parse()

	// Populate arguments
	templatePath = *templatePathPtr
	rootPath = *rootPathPtr
	isRecursive = *isRecursivePtr
	shouldCleanBackup = *shouldCleanBackupPtr
	extSlice := strings.Split(strings.Replace(*extensionsPtr, " ", "", -1), ",")
	for _, val := range extSlice {
		extensions[val] = true
	}
	if !(*rootPathPtr == "./" || *rootPathPtr == "." || *rootPathPtr == "") {
		relativePath = *rootPathPtr
	}

	// Print argument results
	fmt.Println("[+] Using license template:", *templatePathPtr)
	fmt.Println("[+] Searching files in:", *rootPathPtr)
	fmt.Println("[+] Formattign extensions:", *extensionsPtr)
	fmt.Println("[+] Format is recursive?:", *isRecursivePtr)
	fmt.Println("[+] Should clean back up files?:", *shouldCleanBackupPtr)
	fmt.Println("[+] Path to scan:", relativePath)
	fmt.Println("tail:", flag.Args())

	// Load license content
	licenseContent = LoadFileRelative(templatePath)

	// Walk through the directory
	_ = filepath.Walk(relativePath, visit)
}

func visit(path string, f os.FileInfo, err error) error {

	relPath := strings.Replace(path, relativePath, "", -1)
	fmt.Println("[-] Visiting: ", relPath)

	pathLen := len(strings.Split(relPath, "/"))
	// Ignore root
	if pathLen < 2 {
		fmt.Println("[-] Ignoring root folder")
		return nil
	}
	// Ignore subfolders
	if !isRecursive {
		if pathLen > 2 {
			fmt.Println("[-] Ignoring sub folder...")
			return nil
		}
	}
	// Ignore folders
	if f.IsDir() {
		fmt.Println("[-] Ignoring directory...")
		return nil
	}

	// Check file extension
	pathComponents := strings.Split(relPath, ".")
	componentLen := len(pathComponents)
	if componentLen < 1 {
		return nil
	}
	fileExt := pathComponents[componentLen-1]

	ok := extensions[fileExt]
	if !ok {
		fmt.Println("[-] Ignoring a non-code file...")
		return nil
	}

	// Open a new file
	fileProcessRecord[path] = true

	newFilePath := path + "_tmp"
	oldFilePath := path + "_bkup"

	if _, err := os.Stat(newFilePath); err == nil {
		// File already exists
		err = os.Remove(newFilePath)
		checkErr(err)
	}
	if _, err := os.Stat(oldFilePath); err == nil {
		// File already exists
		err = os.Remove(oldFilePath)
		checkErr(err)
	}

	newFile, err := os.Create(newFilePath)
	checkErr(err)

	_, err = newFile.Write(licenseContent)
	checkErr(err)
	// Adding white spaces between license and code
	newFile.WriteString("\n")

	currentFile, err := os.Open(path)
	checkErr(err)
	// Make sure we close the file no matter what.
	defer func() {
		// Close all files in use
		err := currentFile.Close()
		if err != nil {
			fmt.Println("[!] Could not close file:", path)
			fmt.Println("[!] You may need to manually delete it!!!")
		}
		err = newFile.Close()
		if err != nil {
			fmt.Println("[!] Could not close temporary new file:", newFilePath)
			fmt.Println("[!] You may need to manually delete it!!!")
		}
		// Rename files
		err = os.Rename(path, oldFilePath)
		if err != nil {
			fmt.Println("[!] Could not move current file as an back up:", oldFilePath)
			fmt.Println("[!] You may need to restart the program!!!")
		}

		err = os.Rename(newFilePath, path)
		if err != nil {
			fmt.Println("[!] Could not move new tmp file to code:", newFilePath)
			fmt.Println("[!] You may need to recover this file from:", oldFilePath)
		}

		if shouldCleanBackup && err == nil {
			err = os.Remove(oldFilePath)
			if err != nil {
				fmt.Println("[!] Could not delete back up file:", oldFilePath)
				fmt.Println("[!] You may need to manually remove this file:", oldFilePath)
			}
		}

	}()

	scanner := bufio.NewScanner(currentFile)

	var startedOldLicense, finishedOldLicense bool // default to false

	for scanner.Scan() {
		testText := strings.TrimSpace(scanner.Text())

		// Skip empty lines at top
		if testText == "" && !finishedOldLicense {
			fmt.Println("[...]Skiping top empty lines")
			continue
		}
		if !startedOldLicense {
			if strings.HasPrefix(testText, "//") {
				// Skip all the way to the first line of code
				for scanner.Scan() {
					nextTestText := strings.TrimSpace(scanner.Text())
					if !(nextTestText == "" || strings.HasPrefix(nextTestText, "//")) {
						_, err = newFile.Write(scanner.Bytes())
						checkErr(err)
						newFile.WriteString("\n")
						break
					}
					fmt.Println("[...]Skiping top comment and empty lines in between")
				}
			} else if strings.HasPrefix(testText, "/*") {
				// Skip all the way the the first comment end
				if !strings.Contains(testText, "*/") {
					for scanner.Scan() {
						if strings.Contains(scanner.Text(), "*/") {
							break
						}
					}
				}
			} else { // First line of code, let's go ahead and write down
				_, err = newFile.Write(scanner.Bytes())
				checkErr(err)
				newFile.WriteString("\n")
			}
			startedOldLicense = true
			finishedOldLicense = true
		} else {
			_, err = newFile.Write(scanner.Bytes())
			checkErr(err)
			newFile.WriteString("\n")
		}
	}

	return nil

}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
