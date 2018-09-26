## GO - The First Package
### GO入门
#### GOPATH环境变量
GOPATH 环境变量指定了你的工作空间位置。它或许是你在开发Go代码时， 唯一需要设置的环境变量。

首先创建一个工作空间目录，并设置相应的 GOPATH。你的工作空间可以放在任何地方， 在此文档中我们使用 $HOME/work。注意，它绝对不能和你的Go安装目录相同。 （另一种常见的设置是 GOPATH=$HOME。）

```cmd
$ mkdir $HOME/work
$ export GOPATH=$HOME/work
```
将此工作空间的 bin 子目录添加到你的 PATH 中：
```cmd
$ export PATH=$PATH:$GOPATH/bin
```
![配置环境变量](https://segmentfault.com/img/bVbhwfg?w=432&h=58)
#### 包路径
标准库中的包有给定的短路径，比如 "fmt" 和 "net/http"。 对于你自己的包，你必须选择一个基本路径，来保证它不会与将来添加到标准库， 或其它扩展库中的包相冲突。

如果你将你的代码放到了某处的源码库，那就应当使用该源码库的根目录作为你的基本路径。 例如，若你在 GitHub 上有账户 github.com/user 那么它就应该是你的基本路径。

注意，在你能构建这些代码之前，无需将其公布到远程代码库上。只是若你某天会发布它， 这会是个好习惯。在实践中，你可以选择任何路径名，只要它对于标准库和更大的Go生态系统来说， 是唯一的就行。

我们将使用 github.com/user 作为基本路径。在你的工作空间里创建一个目录， 我们将源码存放到其中：
```cmd
$ mkdir -p $GOPATH/src/github.com/user
```
#### 创建第一个包
1. 创建包目录：
```cmd
$ mkdir $GOPATH/src/github.com/user/TheFirstPackage
$ mkdir $GOPATH/src/github.com/user/TheFirstPackage/hello
$ mkdir $GOPATH/src/github.com/user/TheFirstPackage/stringutil
```
![创建包目录](https://segmentfault.com/img/bVbhwfw?w=591&h=74)
2. 添加文件：
```cmd 
$ code $GOPATH/src/github.com/user/TheFirstPackage/hello/hello.go
$ code $GOPATH/src/github.com/user/TheFirstPackage/stringutil/reverse.go
```
3. 添加代码：
> hello.go的代码
```go
package main

import (
	"fmt"

	"github.com/user/stringutil"
)

func main() {
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
}
```
> reverse.go的代码
```go
// stringutil 包含有用于处理字符串的工具函数。
package stringutil

// Reverse 将其实参字符串以符文为单位左右反转。
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```
4. 编译运行：
    * 首先编译stringutil包：
    ```cmd
    $ go install github.com/user/TheFirstPackage/hello
    ```
    * 编译hello.go文件：
    ```cmd
    go install github.com/user/TheFirstPackage/hello
    ```
    * 运行hello程序：
    ```cmd
    $ hello
    ```
![运行结果](https://segmentfault.com/img/bVbhwfz?w=728&h=61)
### Git入门
##### Git日常使用总结
1. 开始对项目进行git管理，只需到此项目所在的目录，输入：
    > git init
2. 从现有仓库克隆：
    > git clone git://github.com/schacon/grit.git
3. 如果当前目录下有几个文件想要纳入版本控制，需要先用 git add 命令告诉 Git 开始对这些文件进行跟踪，然后提交：
    > git add readme.txt(将文本添加到本地git库)
4. 检查当前文件状态：
    > git status(这句话查看添加的状态)
5. 提交项目到GitHub：
    > git commit -m "first"(可以将他理解为注释，当上传项目到github时，作用是标注我上传该项目时想要说的话)
6. 保存下我们当前本地的修改：
    > git stash
7. 当前本地将回到修改前的commit，此时可以继续更新远程分支上的提交：
    > git pull origin
8. 将我们之前所做的修改与当前库上最新的代码进行合并：
    > git stash pop
9.  给文件改名：
    > git mv file_from file_to
10. 查看提交历史:
    > git log
11. 修改最后一次提交即撤销刚刚的操作:
    > git commit --amend
12. 取消对文件的修改：
    > git checkout -- （文件名）
13. 查看当前的远程仓库：
    > git remote (-v 将所有的远程仓库列出)
14. 添加远程仓库：
    > git remote add [shortname] [url]
15. 可以用下面的命令从远程仓库抓取数据到本地:
    > git fetch [remote-name]
16. 查看远程仓库信息:
    > git remote show [remote-name]
17. 修改某个远程仓库在本地的简称:
    > git remote rename
18. 移除对应的远端仓库:
    > git remote rm
19. 暂存当前本地工作临时状态:
    * 会将当前的工作状态进行缓存，默认有递增id
    > git stash
    * 在将当前的工作状态进行缓存的同时，给本地缓存进行备注
    > git stash save "description"
    * 查看当前已经缓存列表
    > git stash list
    * 根据在stash list中的缓存列表id，可以指定恢复到某一个工作状态
    > git stash apply stash@{id}
    * 与apply命令不同，执行pop命令之后，stash缓存中即会移除该stash；
    > git stash pop stash@{id}
    * 清除所有stash缓存
    > git stash clear
    * 移除某个stash
    > git stash drop stash@{id}
#### Git参考网站
> [Git官方中文文档](https://git-scm.com/book/zh/v2)
> 
> [Git简单教程博客](https://www.liaoxuefeng.com/wiki/0013739516305929606dd18361248578c67b8067c8c017b000)