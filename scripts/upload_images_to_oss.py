#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import pandas as pd
import base64
import requests
import json
import os
from pathlib import Path
import logging
import yaml
from typing import Any, Dict


# 设置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class ImageUploader:
    def __init__(self, api_url="http://172.31.61.26:16010/api/admin/upload/image", token="token"):
        self.api_url = api_url
        self.token = token
        self.headers = {
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        }
    
    def encode_image_to_base64(self, image_path):
        """将图片文件编码为base64格式"""
        try:
            with open(image_path, 'rb') as image_file:
                image_data = image_file.read()
                # 获取文件扩展名来确定MIME类型
                ext = Path(image_path).suffix.lower()
                mime_type = {
                    '.jpg': 'jpeg',
                    '.jpeg': 'jpeg', 
                    '.png': 'png',
                    '.gif': 'gif',
                    '.bmp': 'bmp',
                    '.webp': 'webp'
                }.get(ext, 'jpeg')
                
                base64_encoded = base64.b64encode(image_data).decode('utf-8')
                return f"data:image/{mime_type};base64,{base64_encoded}"
        except Exception as e:
            logger.error(f"编码图片失败 {image_path}: {e}")
            return None
    
    def upload_image(self, image_path):
        """上传单个图片到OSS"""
        if not os.path.exists(image_path):
            logger.error(f"图片文件不存在: {image_path}")
            return None
            
        # 编码图片
        base64_image = self.encode_image_to_base64(image_path)
        if not base64_image:
            return None
            
        # 准备请求数据
        data = {
            "image_data": base64_image
        }
        
        try:
            response = requests.post(
                self.api_url,
                headers=self.headers,
                data=json.dumps(data),
                timeout=30
            )
            
            if response.status_code == 200:
                result = response.json()
                if result.get('code') == 200:
                    file_name = result.get('data', {}).get('file_name')
                    logger.info(f"上传成功 {image_path} -> {file_name}")
                    return file_name
                else:
                    logger.error(f"API返回错误: {result}")
                    return None
            else:
                logger.error(f"HTTP错误 {response.status_code}: {response.text}")
                return None
                
        except Exception as e:
            logger.error(f"上传图片失败 {image_path}: {e}")
            return None
    
    def process_excel(self, input_file, output_file):
        """处理Excel文件，上传图片并添加OSS短链列"""
        try:
            # 读取Excel文件
            df = pd.read_excel(input_file)
            logger.info(f"读取Excel文件: {input_file}, 共{len(df)}行数据")
            
            # 检查是否存在"本地图片路径"列
            if '本地图片路径' not in df.columns:
                logger.error("Excel文件中未找到'本地图片路径'列")
                return False
                
            # 添加新列
            df['oss短链'] = ''
            
            # 获取脚本所在目录的父目录作为项目根目录
            script_dir = os.path.dirname(os.path.abspath(__file__))
            project_root = os.path.dirname(script_dir)
            
            # 处理每一行
            for index, row in df.iterrows():
                image_path = row['本地图片路径']
                if pd.isna(image_path) or not image_path:
                    logger.warning(f"第{index+1}行图片路径为空")
                    continue
                
                # 如果是相对路径，则相对于项目根目录解析
                if not os.path.isabs(image_path):
                    full_image_path = os.path.join(project_root, image_path)
                else:
                    full_image_path = image_path
                    
                logger.info(f"处理第{index+1}行: {image_path} -> {full_image_path}")
                
                # 上传图片
                oss_link = self.upload_image(full_image_path)
                if oss_link:
                    df.at[index, 'oss短链'] = oss_link
                else:
                    logger.warning(f"第{index+1}行图片上传失败")
            
            # 保存新的Excel文件
            df.to_excel(output_file, index=False)
            logger.info(f"处理完成，保存到: {output_file}")
            return True
            
        except Exception as e:
            logger.error(f"处理Excel文件失败: {e}")
            return False


def read_yaml(file_path: str) -> Dict[str, Any]:
    """
    读取YAML文件并返回字典
    
    Args:
        file_path (str): YAML文件路径
        
    Returns:
        Dict[str, Any]: 解析后的字典数据
        
    Raises:
        FileNotFoundError: 当文件不存在时
        yaml.YAMLError: 当YAML格式错误时
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            data = yaml.safe_load(file)
            return data if data is not None else {}
    except FileNotFoundError:
        raise FileNotFoundError(f"文件不存在: {file_path}")
    except yaml.YAMLError as e:
        raise yaml.YAMLError(f"YAML解析错误: {e}")

def main():
    """主函数"""
    # 获取脚本所在目录的父目录作为项目根目录
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    
    input_file = os.path.join(project_root, "data", "prompts_data.xlsx")
    output_file = os.path.join(project_root, "data", "prompts_data_with_oss.xlsx")
    config_file = os.path.join(project_root, "config", "config.yaml")

    cfg = read_yaml(config_file)
    token = cfg['new_api']['admin_key']
    
    uploader = ImageUploader(token=token)
    success = uploader.process_excel(input_file, output_file)
    
    if success:
        print(f"✅ 图片上传完成，结果保存到: {output_file}")
    else:
        print("❌ 处理失败")

if __name__ == "__main__":
    main()