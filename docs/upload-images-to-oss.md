将通过 scripts/extract_prompts.py 的脚本生成的excel文件： data/prompts_data.xlsx  里面的记录中对应的文件上传到oss中
图片路径的列名： 本地图片路径
获取的oss短链保存到新列： oss短链

上传图片的curl接口：
curl --location --request POST 'http://172.31.61.26:16010/api/admin/upload/image' \
--header 'Authorization: Bearer token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "image_data": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/7QB8UGhvdG9zaG9wIDMuMAA4QklNBAQAAAAAAGAcAigAWkZ..."
}'

响应结构：
{
    "code": 200,
    "message": "success",
    "data": {
        "oss_short_link": "https://xxx.oss-cn-xx.aliyuncs.com/images/20250906/20250906_144852_035aec46-07b9-4c4f-b5fd-0e262d88cc3f.jpg",
        "public_url": "http://xxxx.oss-cn-xxx.aliyuncs.com/images%2F20250906%2F20250906_144852_035aec46-07b9-4c4f-b5fd-0e262d88cc3f.jpg?Expires=1757148532&OSSAccessKeyId=LTAI5tLycAKTsiHNGWhvpK79&Signature=%2FB8qxTCdgjJAVQ52ougXNMOL%2B2g%3D",
        "file_name": "images/20250906/20250906_144852_035aec46-07b9-4c4f-b5fd-0e262d88cc3f.jpg"
    }
}

保存oss短链为 file_name 的值
Authorization的token值先写死token，我自己加

保存成一个新的excle文件
