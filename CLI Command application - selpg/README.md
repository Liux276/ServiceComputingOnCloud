## 服务计算 - 3 | CLI 命令行实用程序开发基础
---
### 概述

CLI（Command Line Interface）实用程序是Linux下应用开发的基础。正确的编写命令行程序让应用与操作系统融为一体，通过shell或script使得应用获得最大的灵活性与开发效率。Linux提供了cat、ls、copy等命令与操作系统交互；go语言提供一组实用程序完成从编码、编译、库管理、产品发布全过程支持；容器服务如docker、k8s提供了大量实用程序支撑云服务的开发、部署、监控、访问等管理任务；git、npm等都是大家比较熟悉的工具。尽管操作系统与应用系统服务可视化、图形化，但在开发领域，CLI在编程、调试、运维、管理中提供了图形化程序不可替代的灵活性与效率。

---
### 基本要求

参阅[selpg命令行程序设计逻辑](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)，实现一个selpg页选择程序，满足selpg设计要求。

---
### 程序实现

#### 流程分析
   ![CLI流程图](https://segmentfault.com/img/bVbh06Q?w=565&h=590)

   程序按照读取参数、判断参数是否合规、读取文件、确定输出位置并输出顺序执行。当发现错误时抛出错误并终止流程。

#### 代码实现
1. selpg所需参数有必须的开始页码-s以及结束页码-e，可选的输入文件名、自定页长-l、遇换页符换页-f和输出地址。其中自定页长和遇换页符换页两个选项是互斥的，不能同时使用。
    * 定义保存参数数据的结构体
        ```go
        type selpgArgs struct {
            startPage  int
            endPage    int
            inFileName string
            pageLen    int
            pageType   bool
            printDest  string
        }
        ```
    * 输入参数使用 github.com/spf13/pflag 包提供的pflag进行处理，pflag包于flag用法类似，但pflag相对于flag能够更好地满足 Unix 命令行规范。参考：[Golang pflag](https://godoc.org/github.com/spf13/pflag)
        ```go
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
        ```
        * pflag包中的函数XXXVarP（XXX为Int、String、Bool等可选类型）可以取出命令行参数名称shorthand的参数的值，value指定*p的默认值，name为自定的名称，usage为自定的该参数的描述。该函数无返回值。
        > func XXXVarP(p *XXX, name, shorthand string, value XXX, usage string)
        
        * 获得flag参数后，要用pflag.Parse()函数才能把参数解析出来。这里还有一个东西，解析完指定的参数后，可以通过调用argLeft := pflag.Args()来获得未定义但输入了的参数如文件名。

2. 命令行参数获取之后，首先要进行参数检查以尽量避免参数谬误。出现错误时输出问题并正常结束程序。参数正确则将各个参数值输出到屏幕上。
    ```go
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
            fmt.Fprintf(os.Stderr, "\n[Error]The pageLen can't be less than 1 !\n")
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
    ```
    * 在这个函数中，首先检查了开始页args.startPage和结束页args.endPage是否被赋值，然后检查开始页args.startPage和结束页args.endPage是否为正数，接下来检查开始页args.startPage是否大于结束页args.endPage，然后检查自定页长-l和遇换页符换页-f是否同时出现，最后判断当自定页长-l出现时args.pageLen是否小于1。遇到不合规的地方正常结束程序，全部合规则输出得到的参数。

3. 参数检查结束之后，程序开始调用excuteCMD函数执行命令。
    ```go
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

        if len(args.printDest) == 0 {
            output2Des(os.Stdout, fin, args.startPage, args.endPage, args.pageLen, args.pageType)
        } else {
            output2Des(cmdExec(args.printDest), fin, args.startPage, args.endPage, args.pageLen, args.pageType)
        }
    }

    func checkFileAccess(filename string) {
        _, errFileExits := os.Stat(filename)
        if os.IsNotExist(errFileExits) {
            fmt.Fprintf(os.Stderr, "\n[Error]: input file \"%s\" does not exist\n", filename)
            os.Exit(7)
        }
    }
    ```
    * 第一步检查输入。如果没有给定文件名，则从标准输入中获取；如果给出读取的文件名，则调用函数checkFileAccess检查文件是否存在。
    * 第二步是打开文件，使用函数checkError检查是否出现错误。如果打开出错则输出错误并抛出恐慌。
    * 第三步判断是否有-d参数。如果没有-d参数，选择的页直接从os.Stdout标准输出中输出。如果-d存在，则从指定的打印通道中输出。

4. 在-d参数存在时，涉及到了os/exec包的使用，这里可以参考[golang中os/exec包用法](https://blog.csdn.net/chenbaoke/article/details/42556949)。
    ```go
    func cmdExec(printDest string) io.WriteCloser {
        cmd := exec.Command("lp", "-d"+printDest)
        fout, err := cmd.StdinPipe()
        checkError(err, "StdinPipe")
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        errStart := cmd.Run()
        checkError(errStart, "CMD run")
        return fout
    }
    ```
    * 其中printDest就是获取的打印地址，将命令行的输入管道cmd.StdinPipe()获取的指针赋值给fout，然后再将fout返回给output2Des函数中的作为输出位置参数的输入，最后在output2Des中将需要的页输出到fout。

5. 输出函数output2Des将输入的文件，按页码要求读取并输出到fout中。
    ```go
        func output2Des(fout interface{}, fin *os.File, pageStart int, pageEnd int, pageLen int, pageType bool) {

        lineCount := 0
        pageCount := 1
        buf := bufio.NewReader(fin)
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
            checkError(err, "file read in")

            if (pageCount >= pageStart) && (pageCount <= pageEnd) {
                var outputErr error
                if stdOutput, ok := fout.(*os.File); ok {
                    _, outputErr = fmt.Fprintf(stdOutput, "%s", line)
                } else if pipeOutput, ok := fout.(io.WriteCloser); ok {
                    _, outputErr = pipeOutput.Write([]byte(line))
                } else {
                    fmt.Fprintf(os.Stderr, "\n[Error]:fout type error. ")
                    os.Exit(8)
                }
                checkError(outputErr, "Error happend when output the pages.")
            }
        }
        if pageCount < pageStart {
            fmt.Fprintf(os.Stderr, "\n[Error]: startPage (%d) greater than total pages (%d), no output written\n", pageStart, pageCount)
            os.Exit(9)
        } else if pageCount < pageEnd {
            fmt.Fprintf(os.Stderr, "\n[Error]: endPage (%d) greater than total pages (%d), less output than expected\n", pageEnd, pageCount)
            os.Exit(10)
        }
    }
    ```
    * bufio包实现了带缓存的 I/O 操作，在文件读取中十分方便。具体使用参见[Golang学习 - bufio 包](https://www.cnblogs.com/maxiaoyun/p/7007755.html)。这里使用buf.ReadString(symbol),每次读取字符串直到遇到字符symbol为止。
    * 由于fout存在两种输入-io.Stdout标准输出作为输入、cmd.StdinPipe()管道作为输入。所以使用空接口interface{}作为fout的类型，借助类型断言stdOutput, ok := fout.(*os.File)和pipeOutput, ok := fout.(io.WriteCloser)来判断fout具体类型并调用相应函数。
---
### 程序测试
* 按文档[ 使用 selpg ](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)章节要求测试该程序
* ```cmd测试文档input_file.txt包含两个换页符```
    ![测试文档](https://segmentfault.com/img/bVbh063?w=617&h=461)
1. ```cmd
    ./selpg -s1 -e1 input_file.txt
    ```
    ![pic1](https://segmentfault.com/img/bVbh06T?w=518&h=315)
2. ```cmd 
    ./selpg -s1 -e1 < input_file.txt
    ```
    ![pic2](https://segmentfault.com/img/bVbh06U?w=537&h=316)
3. ```cmd
    ./selpg -s1 -e2 input_file.txt >output_file
    ```
    ![pic3](https://segmentfault.com/img/bVbh06W?w=615&h=100)
    ![output1](https://segmentfault.com/img/bVbh07H?w=441&h=360)
4. ```cmd
    ./selpg -s1 -e4 input_file.txt 2>error_file
    ```
    ![pic4](https://segmentfault.com/img/bVbh07J?w=614&h=332)
    ![output2](https://segmentfault.com/img/bVbh07K?w=625&h=110)
5.  ```cmd
    ./selpg -s1 -e3 input_file.txt >output_file 2>error_file
    ```
    ![pic5](https://segmentfault.com/img/bVbh07M?w=456&h=365)
    ![output3](https://segmentfault.com/img/bVbh07P?w=501&h=118)
6. ```cmd 
    ./selpg -s1 -e2 -f input_file.txt
    ```
    ![pic6](https://segmentfault.com/img/bVbh07R?w=551&h=243)
7. ```cmd 
    ./selpg -s1 -e1 -dlp1 input_file.txt
    ```
    ![pic7](https://segmentfault.com/img/bVbh07R?w=551&h=243)
    ![打印结果](https://segmentfault.com/img/bVbh2kf?w=828&h=598)
---
### [博客地址](https://segmentfault.com/a/1190000016648238)