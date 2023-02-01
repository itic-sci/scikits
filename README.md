# 发布自己的golang库

- 在github上创建一个public仓库，仓库名与go库名一致，然后将该仓库clone到本地
git clone https://github.com/gl1n/gowave.git

- 进入项目文件夹，初始化go mod
go mod init github.com/xulehexuwei/scikits

- 在项目文件夹中编写代码(可以添加子文件夹)，然后push到github
git add * git commit -m "第一次提交" git push

- 上传成功后，还需发布代码。进入GitHub仓库，点击release发布，版本号采用v0.0.0形式

- 发布成功后，测试代码能够被他人使用
go get github.com/gl1n/gowave

- 需要在Goland中勾选Settings>Go>Go Modules的Enable Go modules integration以正确识别外部导入的包


## 功能

- jwt
- minio
- gin response
- time_fun
- viper
- zap_logger
- randomStr 得到指定长度的随机字符串