
import oss2
from oss2.credentials import EnvironmentVariableCredentialsGetter

# 初始化（请替换为您的配置）
auth = oss2.ProviderAuth(EnvironmentVariableCredentialsGetter())
bucket = oss2.Bucket(auth, 'https://oss-cn-shenzhen.aliyuncs.com', 'xxxxx')

# 定义原始对象键
object_key = 'images/20250906/20250906_162422_3475ce20-1887-4401-95a0-dd1fca8c2d32.png'

# 定义要签名的查询参数，包含图片处理指令
params = {
    'x-oss-process': 'image/resize,m_fill,w_330,h_240'
}

# 生成一个7200秒后过期的签名URL，并包含我们的图片处理参数
signed_url = bucket.sign_url('GET', object_key, 3600, params=params)

print(signed_url)
# 输出：https://xxxxx.oss-cn-shenzhen.aliyuncs.com/images/20250906/20250906_...png?x-oss-process=image%2Fresize%2Cm_fill%2Cw_330%2Ch_240&OSSAccessKeyId=...&Expires=...&Signature=...



在 internal\api\prompts.go 文件中帮我添加一个字段 thumbs 这个字段表示原始图片的缩略图列表