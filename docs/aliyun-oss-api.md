import oss2
import time
import requests


def getAuthedBucket(config):
    # 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录RAM控制台创建RAM账号。
    auth = oss2.Auth(config['SZ_ALIYUN_ACCESS_ID'], config['SZ_ALIYUN_ACCESS_KEY'])
    # Endpoint以杭州为例，其它Region请按实际情况填写。
    bucket = oss2.Bucket(auth, 'http://%s' % config['SZ_ALIYUN_ENDPOINT'], config['SZ_ALIYUN_BUCKET'])
    
    return bucket


class OSSClient:
    def __init__(self, config, upload_end_point, download_end_point, key_id, key_secret, bucket_name,
                 download_domain=None):
        self.config = config
        self.upload_end_point = upload_end_point
        self.download_end_point = download_end_point
        self.key_id = key_id
        self.key_secret = key_secret
        self.bucket_name = bucket_name
        self.download_domain = download_domain
    
    # 将文件上传至阿里云 OSS
    def upload_file(self, file, file_name, download_filename=None, private=True, retry_times=1, retry_interval=10):
        # 开始上传文件
        for retry_time in range(retry_times):
            try:
                auth = oss2.Auth(self.key_id, self.key_secret)
                bucket = oss2.Bucket(auth, self.upload_end_point, self.bucket_name)
                
                headers = {
                    'x-oss-object-acl': 'private' if private else 'public-read'
                }
                if download_filename is not None:
                    """
                    将测试.txt从OSS下载到本地后，需要保留文件名为测试.txt，
                    需按照
                    "attachment;filename="+URLEncoder.encode("测试","UTF-8")+".txt;filename*=UTF-8''"+URLEncoder.encode("测试","UTF-8")+".txt")
                    的格式设置Content-Disposition，
                    即attachment;filename=%E6%B5%8B%E8%AF%95.txt;filename*=%E6%B5%8B%E8%AF%95.txt
                    """
                    encode_filename = requests.utils.quote(download_filename)
                    headers['Content-Disposition'] = 'attachment;filename=%s;filename*=%s' % (
                        encode_filename, encode_filename)
                result = bucket.put_object(file_name, file, headers=headers)
                # 如果返回值不等于 2xx 则表示失败
                if result.status < 200 or result.status > 299:
                    time.sleep(retry_interval)
                    continue
                return True
            except Exception as e:
                print(
                    'oss_client: upload file error. retry: {0}/{1} {2}'.format(str(retry_time + 1), str(retry_times),
                                                                               e))
                time.sleep(retry_interval)
                continue
        
        # 到达重试次数还没有上传成功
        return False
    
    # 获取阿里云 OSS 的临时下载地址
    def get_file_url(self, file_name, expires):
        # 获取文件的临时下载地址
        try:
            auth = oss2.Auth(self.key_id, self.key_secret)
            bucket = oss2.Bucket(auth, self.download_end_point, self.bucket_name)
            acl = bucket.get_object_acl(file_name).acl
            if acl == 'public-read':
                file_url = self.get_public_download_url(file_name)
            else:
                file_url = bucket.sign_url('GET', file_name, expires)
                # file_url = self.replace_private_download_domain(file_url)
        except Exception as e:
            print('oss_client: get file url error.', e)
            return None
        
        return file_url
    
    # 从阿里云 OSS 中下载文件
    def download_file(self, file_name):
        # 获取文件的临时下载地址
        try:
            auth = oss2.Auth(self.key_id, self.key_secret)
            bucket = oss2.Bucket(auth, self.upload_end_point, self.bucket_name)
            remote_stream = bucket.get_object(file_name)
        except Exception as e:
            print('oss_client: download file error: \n{0}'.format(e))
            return None
        
        return remote_stream
    
    def delete_file(self, file_name):
        # 删除文件
        try:
            auth = oss2.Auth(self.key_id, self.key_secret)
            bucket = oss2.Bucket(auth, self.upload_end_point, self.bucket_name)
            result = bucket.delete_object(file_name)
            # 如果返回值不等于 2xx 则表示失败
            if result.status < 200 or result.status > 299:
                return False
        except Exception as e:
            print('oss_client: delete file error: \n{0}'.format(e))
            return False
        return True
    
    def get_public_download_url(self, file_name):
        if file_name is None or len(file_name) == 0:
            return ''
        
        file_url = ''
        if self.download_domain is not None and len(self.download_domain) > 0:
            file_url = self.download_domain
            if file_url[:-1] != '/':
                file_url += '/'
            file_url += file_name
        else:
            if self.download_end_point[:7].lower() != 'http://' and self.download_end_point[:8].lower() != 'https://':
                file_url += 'https://' + self.bucket_name + '.'
            file_url += self.download_end_point + '/' + file_name
        return file_url


if __name__ == '__main__':
    import os
    import sys
    import uuid
    
    # # 导入scripts下的utils工具库
    BASE_PATH = os.path.dirname(os.path.dirname(os.path.dirname(__file__)))
    sys.path.append(BASE_PATH)
    from config import Config
    from core.utils import timex
    #
    config = Config()
    oss_conf = config['oss']['aliyun']
    print(oss_conf)
    oss_cli = OSSClient(oss_conf, oss_conf['SZ_ALIYUN_ENDPOINT'], oss_conf['SZ_ALIYUN_ENDPOINT']
                        , oss_conf['SZ_ALIYUN_ACCESS_ID'], oss_conf['SZ_ALIYUN_ACCESS_KEY']
                        , oss_conf['SZ_ALIYUN_BUCKET']
                        )
    # #
    # # file key:  resources/f0/2022-04-30/8ef1cda30a6a437da147d9eac1f3fa03.mobi
    # # upload_file()
    #
    # get_temp_download_url("resources/f0/2022-04-30/8ef1cda30a6a437da147d9eac1f3fa03.mobi")
    
    # 上传static静态文件夹
    # g = os.walk(r"H:\tools\talebook_data\books\static")
    #
    # count = 0
    # for path, dir_list, file_list in g:
    #     for file_name in file_list:
    #         count += 1
    #         full_path = os.path.join(path, file_name)
    #         oss_key = 'website/'+full_path[full_path.index('static'):].replace('\\', '/')
    #         # print(full_path, '\t', oss_key)
    #         with open(full_path, 'rb') as fp:
    #             if oss_cli.upload_file(fp, oss_key, None, False):
    #                 print("上传成功")
    #
    #                 print("file key: ", oss_key)
    #
    # print("共%d个文件" % count)

        

参考上面的python代码帮我在go中实现一个oss文件上传的通用工具组件，这个公共组件添加到  pkg 文件夹下

配置请参考 config/config.go 文件
pkg下面的实现可以参考 pkg/database/mysql.go 文件
提供一个单元测试文件

（这个已经实现了，代码文件在 pkg/oss/oss.go）