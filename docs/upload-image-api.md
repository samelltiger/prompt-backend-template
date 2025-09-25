
帮我开发一个将传入的图片的base64编码的文件上传打oss的接口，post请求，返回一个oss短链和可以公开访问的url（2小时过期）。
接口需要定义在 internal/api/admin 文件夹内，需要admin鉴权

参考 internal/api/admin/redemption.go 文件开发接口
参考 pkg/oss/oss.go 定义的方法调用oss组件

