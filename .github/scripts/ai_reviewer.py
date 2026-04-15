#!/usr/bin/env python3
"""
AI代码审查脚本 - 使用阿里云DashScope API
支持web前后端项目代码审查
"""

import os
import sys
import json
import requests
import subprocess
from pathlib import Path

# 配置参数
DASHSCOPE_API_KEY = os.getenv("DASHSCOPE_API_KEY")
DASHSCOPE_BASE_URL = os.getenv("DASHSCOPE_BASE_URL", "https://coding.dashscope.aliyuncs.com/v1")
MODEL = "glm-5"  # 可选择 qwen-turbo, qwen-plus, qwen-max

# GitHub环境变量
GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")
PR_NUMBER = os.getenv("PR_NUMBER")
GITHUB_REPOSITORY = os.getenv("GITHUB_REPOSITORY")

def get_changed_files():
    """获取PR中变更的文件列表"""
    try:
        # 获取当前分支和目标分支
        result = subprocess.run(
            ["git", "log", "--oneline", "-1"],
            capture_output=True,
            text=True,
            check=True
        )
        
        # 获取变更的文件
        result = subprocess.run(
            ["git", "diff", "--name-only", "HEAD~1"],
            capture_output=True,
            text=True,
            check=True
        )
        
        files = [f.strip() for f in result.stdout.split('\n') if f.strip()]
        return files
    except subprocess.CalledProcessError as e:
        print(f"Error getting changed files: {e}")
        return []

def read_file_content(file_path):
    """读取文件内容"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()
    except Exception as e:
        print(f"Error reading file {file_path}: {e}")
        return None

def analyze_code_with_dashscope(file_content, file_type):
    """使用DashScope API分析代码"""
    
    # 根据文件类型设置不同的提示词
    if file_type in ['js', 'jsx', 'ts', 'tsx', 'vue']:
        prompt = f"""
        请对以下{file_type.upper()}代码进行专业审查：
        1. 代码质量评估
        2. 潜在bug识别
        3. 性能优化建议
        4. 安全性检查
        5. 代码规范检查
        6. 最佳实践建议
        
        代码内容：
        {file_content[:2000]}  # 限制内容长度
        
        请以JSON格式返回结果，包含以下字段：
        - score: 代码质量评分（1-10）
        - summary: 总体评价
        - issues: 问题列表（包含问题描述、严重程度、修复建议）
        - suggestions: 改进建议
        """
    elif file_type in ['py']:
        prompt = f"""
        请对以下Python代码进行专业审查：
        1. PEP 8规范检查
        2. 代码逻辑分析
        3. 异常处理评估
        4. 性能优化建议
        5. 安全性检查
        
        代码内容：
        {file_content[:2000]}
        
        请以JSON格式返回结果
        """
    elif file_type in ['java', 'kt','go']:
        prompt = f"""
        请对以下Java/Kotlin/Golang代码进行专业审查：
        1. 代码规范检查
        2. OOP原则遵循情况
        3. 异常处理
        4. 性能优化
        5. 线程安全性
        
        代码内容：
        {file_content[:2000]}
        
        请以JSON格式返回结果
        """
    else:
        prompt = f"""
        请对以下代码进行专业审查：
        1. 代码质量评估
        2. 潜在问题识别
        3. 改进建议
        
        代码内容：
        {file_content[:2000]}
        
        请以JSON格式返回结果
        """
    
    headers = {
        "Authorization": f"Bearer {DASHSCOPE_API_KEY}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "model": MODEL,
        "input": {
            "messages": [
                {
                    "role": "system",
                    "content": "你是一位专业的代码审查专家，擅长发现代码中的问题并提供改进建议。"
                },
                {
                    "role": "user",
                    "content": prompt
                }
            ]
        },
        "parameters": {
            "result_format": "text"
        }
    }
    
    try:
        response = requests.post(
            f"{DASHSCOPE_BASE_URL}/chat/completions",
            headers=headers,
            json=payload,
            timeout=30
        )
        
        if response.status_code == 200:
            result = response.json()
            return result.get("output", {}).get("text", "")
        else:
            print(f"API请求失败: {response.status_code}")
            return None
    except Exception as e:
        print(f"调用DashScope API出错: {e}")
        return None

def post_github_comment(review_results):
    """将审查结果发布到GitHub PR评论"""
    if not GITHUB_TOKEN or not PR_NUMBER:
        print("缺少GitHub Token或PR编号")
        return
    
    comment_body = "## AI代码审查报告\n\n"
    
    for file_path, result in review_results.items():
        if result:
            try:
                # 尝试解析JSON结果
                review_data = json.loads(result)
                comment_body += f"### 📁 {file_path}\n"
                comment_body += f"**评分**: {review_data.get('score', 'N/A')}/10\n"
                comment_body += f"**总结**: {review_data.get('summary', '')}\n\n"
                
                issues = review_data.get('issues', [])
                if issues:
                    comment_body += "**发现的问题**:\n"
                    for issue in issues:
                        comment_body += f"- ⚠️ {issue.get('description', '')}\n"
                        comment_body += f"  严重程度: {issue.get('severity', '未知')}\n"
                        comment_body += f"  建议: {issue.get('suggestion', '')}\n\n"
                
                suggestions = review_data.get('suggestions', [])
                if suggestions:
                    comment_body += "**改进建议**:\n"
                    for suggestion in suggestions:
                        comment_body += f"- 💡 {suggestion}\n"
                
                comment_body += "---\n\n"
            except json.JSONDecodeError:
                # 如果不是JSON格式，直接显示文本
                comment_body += f"### 📁 {file_path}\n"
                comment_body += f"{result}\n\n"
                comment_body += "---\n\n"
    
    # 添加总结
    comment_body += "### 📋 审查总结\n"
    comment_body += "AI代码审查已完成。请仔细查看以上建议，并根据需要修改代码。\n\n"
    comment_body += "*Powered by 阿里云DashScope*"
    
    # 发布评论到GitHub
    url = f"https://api.github.com/repos/{GITHUB_REPOSITORY}/issues/{PR_NUMBER}/comments"
    headers = {
        "Authorization": f"token {GITHUB_TOKEN}",
        "Accept": "application/vnd.github.v3+json"
    }
    
    data = {
        "body": comment_body
    }
    
    try:
        response = requests.post(url, headers=headers, json=data)
        if response.status_code == 201:
            print("✅ 审查评论已成功发布到GitHub PR")
        else:
            print(f"❌ 发布评论失败: {response.status_code}")
    except Exception as e:
        print(f"❌ 发布评论时出错: {e}")

def main():
    """主函数"""
    if not DASHSCOPE_API_KEY:
        print("❌ 未设置DASHSCOPE_API_KEY环境变量")
        sys.exit(1)
    
    print("🔍 开始AI代码审查...")
    
    # 获取变更的文件
    changed_files = get_changed_files()
    if not changed_files:
        print("📭 未发现变更的文件")
        return
    
    print(f"📂 发现 {len(changed_files)} 个变更的文件")
    
    review_results = {}
    
    # 分析每个变更的文件
    for file_path in changed_files:
        if not os.path.exists(file_path):
            print(f"⚠️  文件不存在: {file_path}")
            continue
        
        # 检查文件类型（只审查代码文件）
        ext = file_path.split('.')[-1].lower()
        supported_exts = ['js', 'jsx', 'ts', 'tsx', 'vue', 'py', 'java', 'kt', 'go', 'cpp', 'c', 'h', 'php', 'rb', 'rs', 'swift']
        
        if ext not in supported_exts:
            print(f"⏭️  跳过非代码文件: {file_path}")
            continue
        
        print(f"🔎 正在审查: {file_path}")
        
        # 读取文件内容
        content = read_file_content(file_path)
        if not content:
            continue
        
        # 使用DashScope API分析代码
        result = analyze_code_with_dashscope(content, ext)
        if result:
            review_results[file_path] = result
            print(f"✅ 完成审查: {file_path}")
        else:
            print(f"❌ 审查失败: {file_path}")
    
    # 发布审查结果到GitHub
    if review_results:
        post_github_comment(review_results)
        print("🎉 AI代码审查完成！")
    else:
        print("📭 未生成有效的审查结果")

if __name__ == "__main__":
    main()
