package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

type selpgArgs struct {
	startPage  int
	endPage    int
	inFileName string
	pageLen    int
	pageType   bool
	printDest  string
}

func main() {

	var args selpgArgs
	getArgs(&args)
	checkArgs(&args)
	excuteCMD(&args)

}

func getArgs(args *selpgArgs) {

	pflag.IntVarP(&(args.startPage), "startPage", "s", -1, "Define startPage")
	pflag.IntVarP(&(args.endPage), "endPage", "e", -1, "Define endPage")
	pflag.IntVarP(&(args.pageLen), "pageLength", "l", 72, "Define pageLength")
	pflag.StringVarP(&(args.printDest), "printDest", "d", "", "Define printDest")
	pflag.BoolVarP(&(args.pageType), "pageType", "f", false, "Define pageType")
	pflag.Parse()

	argLeft := pflag.Args()
	if len(argLeft) > 0 {
		args.inFileName = string(argLeft[0])
	} else {
		args.inFileName = ""
	}
}

func checkArgs(args *selpgArgs) {

	if (args.startPage == -1) || (args.endPage == -1) {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage and endPage can't be empty! Please check your command!\n")
		os.Exit(2)
	} else if (args.startPage <= 0) || (args.endPage <= 0) {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage and endPage can't be negative! Please check your command!\n")
		os.Exit(3)
	} else if args.startPage > args.endPage {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage can't be bigger than the endPage! Please check your command!\n")
		os.Exit(4)
	} else if (args.pageType == true) && (args.pageLen != 72) {
		fmt.Fprintf(os.Stderr, "\n[Error]The command -l and -f are exclusive, you can't use them together!\n")
		os.Exit(5)
	} else if args.pageLen <= 0 {
		fmt.Fprintf(os.Stderr, "\n[Error]The pageLen can't be less than 1 ! Please check your command!\n")
		os.Exit(6)
	} else {
		pageType := "page length."
		if args.pageType == true {
			pageType = "The end sign /f."
		}
		fmt.Printf("\n[ArgsStart]\n")
		fmt.Printf("startPage: %d\nendPage: %d\ninputFile: %s\npageLength: %d\npageType: %s\nprintDestation: %s\n[ArgsEnd]", args.startPage, args.endPage, args.inFileName, args.pageLen, pageType, args.printDest)
	}

}

func checkError(err error, object string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error]%s:", object)
		panic(err)
	}
}

func excuteCMD(args *selpgArgs) {
	var fin *os.File
	if args.inFileName == "" {
		fin = os.Stdin
	} else {
		checkFileAccess(args.inFileName)
		var err error
		fin, err = os.Open(args.inFileName)
		checkError(err, "File input")
	}

	output2Des(args.printDest, fin, args.startPage, args.endPage, args.pageLen, args.pageType)
}

func checkFileAccess(filename string) {
	_, errFileExits := os.Stat(filename)
	if os.IsNotExist(errFileExits) {
		fmt.Fprintf(os.Stderr, "\n[Error]: input file \"%s\" does not exist\n", filename)
		os.Exit(7)
	}
}

func cmdExec(printDest string) (*exec.Cmd, io.WriteCloser) {
	cmd := exec.Command("lp", "-d"+printDest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fout, err := cmd.StdinPipe()
	checkError(err, "Input pipe open\n")
	return cmd, fout
}

func output2Des(printDest string, fin *os.File, pageStart int, pageEnd int, pageLen int, pageType bool) {

	lineCount := 0
	pageCount := 1
	buf := bufio.NewReader(fin)

	var cmd *exec.Cmd
	var fout io.WriteCloser
	if len(printDest) > 0 {
		cmd, fout = cmdExec(printDest)
	}

	for true {

		var line string
		var err error
		if pageType {
			//If the command argument is -f
			line, err = buf.ReadString('\f')
			pageCount++
		} else {
			//If the command argument is -lnumber
			line, err = buf.ReadString('\n')
			lineCount++
			if lineCount > pageLen {
				pageCount++
				lineCount = 1
			}
		}

		if err == io.EOF {
			break
		}
		checkError(err, "file read in\n")

		if (pageCount >= pageStart) && (pageCount <= pageEnd) {
			var outputErr error
			if len(printDest) == 0 {
				_, outputErr = fmt.Fprintf(os.Stdout, "%s", line)
			} else {
				_, outputErr = fout.Write([]byte(line))
				checkError(outputErr, "pipe input")
			}
			checkError(outputErr, "Error happend when output the pages.")
		}
	}

	if len(printDest) != 0 {
		fout.Close()
		errStart := cmd.Run()
		checkError(errStart, "CMD Run")
	}

	if pageCount < pageStart {
		fmt.Fprintf(os.Stderr, "\n[Error]: startPage (%d) greater than total pages (%d), no output written\n", pageStart, pageCount)
		os.Exit(9)
	} else if pageCount < pageEnd {
		fmt.Fprintf(os.Stderr, "\n[Error]: endPage (%d) greater than total pages (%d), less output than expected\n", pageEnd, pageCount)
		os.Exit(10)
	}
}
