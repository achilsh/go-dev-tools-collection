* 如果git add 或者 git commit 后，需要撤销这些修改； 
* 常见使用：git reset --mixed HEAD^ 是不删除工作空间的修改代码，撤销 commit, 也 撤销 git add 的提交。
*  git reset --soft HEAD^ 是 不删除工作空间的修改代码，撤销 commit， 但是不撤销 git add 的提交.
*  git reset --hard HEAD^ 是 删除工作空间改动代码，撤销commit，撤销git add . 注意完成这个操作后，就恢复到了上一次的commit状态。
*  

* 查看文件修改或者分支修改： git log -p 文件名 或者 分支名


* 合并多次已经提交的commit：（因为前期随意commit，现在需要做规范处理。） 使用 git rebase -i  before-begin-to-rebase-commitid 其中 before-begin-to-rebase-commitid 是要开始 rebase commitid的前一个id。再修改合并的 commit：把要被合并的commitid前面的pick改成squash， 这样这个commitid就会被合并到 前一个pick commitid上，只要想合并的commitid，都可以在前面的Pick改成 squash.

* 以行为单位查看提交commit： git log --oneline 
* git log --pretty=reference or git log --pretty=oneline

* 查看git的所有历史commit，包括reset的。命令: git reflog 
* 

* 修改 commit 的 msg, 存在场景：有的时候着急提交，或者暂时没有想到好的 commit msg，就匆匆写了个 "fix", 然后就更新上去了；存在场景：
* 1） 如果在这个 commit 前没有别的修改，只是向修改下当前的提交消息注释： 可以直接： git commit --amend -m "to message."
* 2） 如果本地的修改合并到上一个commit（不想对本次的修改添加消息）：可以直接： git commit --amend --no-edit
* 3） 如果 是修改 很早之前的 commit message,该 commit的前后都有其他的提交。git rebase -i commitId 其中commitid是 需要修改commit msg的前一个commitid， 编辑文件把 目标的 commit 前面的 pick 改成 reword。

* 从特性分支或者commitid 切出分支： git checkout -b new_branch_name      CommitId
