#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import pandas as pd
import requests
import json
import os
import logging
import yaml
from typing import Any, Dict, List

# 设置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class PromptDataProcessor:
    def __init__(self, api_base_url="http://172.31.61.26:16010/api", token="token"):
        self.api_base_url = api_base_url
        self.token = token
        self.headers = {
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        }
    
    def read_excel_prompts(self, excel_file):
        """读取Excel文件中的提示词数据"""
        try:
            # 读取Excel文件
            df = pd.read_excel(excel_file)
            logger.info(f"读取Excel文件: {excel_file}, 共{len(df)}行数据")
            
            # 检查必需的列是否存在
            required_columns = ['标题', '分类', '提示词', '图片描述', 'oss短链']
            missing_columns = [col for col in required_columns if col not in df.columns]
            if missing_columns:
                logger.error(f"Excel文件中缺少必需的列: {missing_columns}")
                return None
            
            # 转换数据格式
            prompts_data = []
            for index, row in df.iterrows():
                # 跳过空行
                if pd.isna(row['标题']) or pd.isna(row['分类']) or pd.isna(row['提示词']):
                    logger.warning(f"第{index+1}行数据不完整，跳过")
                    continue
                
                # 处理OSS短链，如果是单个值则转换为数组
                oss_links = row['oss短链']
                if pd.isna(oss_links) or not oss_links:
                    logger.warning(f"第{index+1}行OSS短链为空，跳过")
                    continue
                
                # 如果oss_links是字符串，转换为数组
                if isinstance(oss_links, str):
                    # 尝试以逗号分隔
                    oss_links_array = [link.strip() for link in oss_links.split(',') if link.strip()]
                else:
                    oss_links_array = [str(oss_links)]
                
                prompt_info = {
                    "title": str(row['标题']),
                    "category_name": str(row['分类']),
                    "prompt": str(row['提示词']),
                    "image_description": str(row['图片描述']) if not pd.isna(row['图片描述']) else "",
                    "oss_short_links": oss_links_array
                }
                
                prompts_data.append(prompt_info)
                logger.info(f"处理第{index+1}行: {prompt_info['title']} - {prompt_info['category_name']}")
            
            logger.info(f"成功读取{len(prompts_data)}条有效数据")
            return prompts_data
            
        except Exception as e:
            logger.error(f"读取Excel文件失败: {e}")
            return None
    
    def batch_add_prompts(self, prompts_data, batch_size=10):
        """批量添加提示词数据"""
        if not prompts_data:
            logger.error("没有数据需要添加")
            return False
        
        total_success = 0
        total_failure = 0
        
        # 分批处理
        for i in range(0, len(prompts_data), batch_size):
            batch = prompts_data[i:i + batch_size]
            logger.info(f"处理批次 {i//batch_size + 1}: {len(batch)}条记录")
            
            # 准备请求数据
            data = {
                "prompts": batch
            }
            
            try:
                response = requests.post(
                    f"{self.api_base_url}/admin/prompts/batch",
                    headers=self.headers,
                    data=json.dumps(data),
                    timeout=60
                )
                
                if response.status_code == 200:
                    result = response.json()
                    if result.get('code') == 200:
                        batch_result = result.get('data', {})
                        success_count = batch_result.get('success_count', 0)
                        failure_count = batch_result.get('failure_count', 0)
                        errors = batch_result.get('errors', [])
                        
                        total_success += success_count
                        total_failure += failure_count
                        
                        logger.info(f"批次结果: 成功{success_count}, 失败{failure_count}")
                        
                        if errors:
                            for error in errors:
                                logger.error(f"批次错误: {error}")
                    else:
                        logger.error(f"API返回错误: {result}")
                        total_failure += len(batch)
                else:
                    logger.error(f"HTTP错误 {response.status_code}: {response.text}")
                    total_failure += len(batch)
                    
            except Exception as e:
                logger.error(f"批量添加失败: {e}")
                total_failure += len(batch)
        
        logger.info(f"总计结果: 成功{total_success}, 失败{total_failure}")
        return total_failure == 0
    
    def test_list_api(self):
        """测试列表API"""
        try:
            response = requests.get(
                f"{self.api_base_url}/prompts",
                timeout=30
            )
            
            if response.status_code == 200:
                result = response.json()
                if result.get('code') == 200:
                    prompts = result.get('data', [])
                    logger.info(f"列表API测试成功，共{len(prompts)}条记录")
                    
                    # 显示前3条记录作为示例
                    for i, prompt in enumerate(prompts[:3]):
                        logger.info(f"示例 {i+1}: {prompt.get('title')} - {prompt.get('category_name')}")
                    
                    return True
                else:
                    logger.error(f"列表API返回错误: {result}")
                    return False
            else:
                logger.error(f"列表API HTTP错误 {response.status_code}: {response.text}")
                return False
                
        except Exception as e:
            logger.error(f"测试列表API失败: {e}")
            return False

def read_yaml(file_path: str) -> Dict[str, Any]:
    """读取YAML文件并返回字典"""
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
    
    # Excel文件路径（使用上传OSS后的文件）
    excel_file = os.path.join(project_root, "data", "prompts_data_with_oss.xlsx")
    config_file = os.path.join(project_root, "config", "config.yaml")
    
    # 检查文件是否存在
    if not os.path.exists(excel_file):
        print(f"❌ Excel文件不存在: {excel_file}")
        print("请先运行 upload_images_to_oss.py 生成包含OSS短链的Excel文件")
        return
    
    if not os.path.exists(config_file):
        print(f"❌ 配置文件不存在: {config_file}")
        return
    
    # 读取配置
    try:
        cfg = read_yaml(config_file)
        token = cfg['new_api']['admin_key']
    except Exception as e:
        print(f"❌ 读取配置失败: {e}")
        return
    
    # 创建处理器
    processor = PromptDataProcessor(token=token)
    
    # 读取Excel数据
    print("📖 读取Excel数据...")
    prompts_data = processor.read_excel_prompts(excel_file)
    
    if not prompts_data:
        print("❌ 读取Excel数据失败")
        return
    
    print(f"📊 共读取到{len(prompts_data)}条有效数据")
    
    # 批量添加数据
    print("🚀 开始批量添加数据...")
    success = processor.batch_add_prompts(prompts_data, batch_size=5)
    
    if success:
        print("✅ 数据添加完成")
        
        # 测试列表API
        print("🧪 测试列表API...")
        if processor.test_list_api():
            print("✅ 列表API测试成功")
        else:
            print("❌ 列表API测试失败")
    else:
        print("❌ 数据添加失败")

if __name__ == "__main__":
    main()