# 这些是该组织内的其它工具，也都是服务于kratos项目的，能够方便你的开发和使用，当然主要是我自己用
init:
	# 这个工具能够让你方便的添加proto文件
	go install github.com/orzkratos/orzkratos/cmd/orzkratos-add-proto@latest
	# 这个工具能够让你的service实现和proto接口保持同步修改，使用前请先做好代码备份或者提交代码
	go install github.com/orzkratos/orzkratos/cmd/orzkratos-srv-proto@latest
	# 这个工具能让项目升级依赖包
	go install github.com/go-mate/depbump/cmd/depbump@latest
	# 这个工具用于项目的静态检查
	go install github.com/go-mate/go-lint/cmd/go-lint@latest

# 当你的kratos项目在proto里新增接口时，通过这个命令能够在对应的服务里也增加函数逻辑，在删除接口时也能把服务代码改为非导出的，以下是使用的样例
orz:
	cd demo1kratos && orzkratos-srv-proto -auto
	cd demo2kratos && orzkratos-srv-proto -auto

# 当你需要同步源项目的最新修改到 fork 项目时，就可以使用这些操作，在IDE里能直接点击执行
merge-step1:
	git remote add upstream git@github.com:orzkratos/demokratos.git

merge-step2:
	git fetch --no-tags upstream main

merge-step3:
	git checkout main

merge-step4:
	# git merge upstream/main
	# 假如在界面里直接点击执行时，使用 "--no-edit" 就不会弹出对话框
	git merge upstream/main --no-edit
	git status

merge-step5:
	depbump
	cd demo1kratos && depbump directs
	cd demo2kratos && depbump directs

merge-step6:
	go clean -testcache
	go test -v ./...
	cd demo1kratos && go test -v ./...
	cd demo2kratos && go test -v ./...

merge-step7:
	go mod tidy -e
	cd demo1kratos && go mod tidy -e
	cd demo2kratos && go mod tidy -e

merge-step8:
	go-lint
