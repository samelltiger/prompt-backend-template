"""
Extract prompt data from coze-prompt.html and save to Excel file
"""

import pandas as pd
from bs4 import BeautifulSoup
import re
import os
import requests
from urllib.parse import urlparse
import time

def download_image(url, save_dir='images', title=''):
    """Download image from URL and return local path with unique filename"""
    if not url:
        return ''
    
    try:
        # Create images directory if it doesn't exist
        os.makedirs(save_dir, exist_ok=True)
        
        # Get filename from URL
        parsed_url = urlparse(url)
        filename = os.path.basename(parsed_url.path)
        
        # If no filename or invalid, generate one from URL hash
        if not filename or '.' not in filename:
            url_hash = hashlib.md5(url.encode()).hexdigest()[:8]
            filename = f"image_{url_hash}.jpg"
        
        # Create safe filename
        safe_filename = re.sub(r'[^\w\-\.]', '_', filename)
        
        # Generate unique filename to avoid conflicts
        import uuid
        import hashlib
        
        # Get file extension
        name_part, ext_part = os.path.splitext(safe_filename)
        if not ext_part:
            ext_part = '.jpg'
        
        # Create unique identifier using title and timestamp
        unique_id = hashlib.md5(f"{title}_{url}_{time.time()}".encode()).hexdigest()[:8]
        
        # Generate final filename with unique ID
        final_filename = f"{name_part}_{unique_id}{ext_part}"
        local_path = os.path.join(save_dir, final_filename)
        
        # Check again if file exists (should be rare with unique ID)
        counter = 1
        original_path = local_path
        while os.path.exists(local_path):
            name_part, ext_part = os.path.splitext(original_path)
            local_path = f"{name_part}_{counter}{ext_part}"
            counter += 1
        
        # Download image with timeout and retry
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        
        response = requests.get(url, headers=headers, timeout=30)
        response.raise_for_status()
        
        # Save image
        with open(local_path, 'wb') as f:
            f.write(response.content)
        
        print(f"Downloaded: {os.path.basename(local_path)}")
        return local_path
        
    except Exception as e:
        print(f"Error downloading image {url}: {e}")
        return ''

def extract_prompt_data(html_file_path):
    """Extract prompt data from HTML file"""
    
    with open(html_file_path, 'r', encoding='utf-8') as file:
        html_content = file.read()
    
    soup = BeautifulSoup(html_content, 'html.parser')
    
    # Find all prompt cards
    cards = soup.find_all('div', class_='bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300 flex flex-col h-full')
    
    prompt_data = []
    
    for card in cards:
        try:
            # Extract category/tag
            category_tag = card.find('span', class_='bg-blue-500 text-white text-xs font-bold px-2.5 py-0.5 rounded-full')
            category = category_tag.text.strip() if category_tag else ''
            
            # Extract title
            title_tag = card.find('h3', class_='text-lg font-bold text-gray-800 mb-2 line-clamp-1')
            title = title_tag.text.strip() if title_tag else ''
            
            # Extract prompt text
            prompt_tag = card.find('p', class_='whitespace-pre-wrap line-clamp-3')
            prompt_text = prompt_tag.text.strip() if prompt_tag else ''
            # Remove excessive whitespace (multiple consecutive spaces)
            prompt_text = re.sub(r'\s+', ' ', prompt_text)
            
            # Extract image URL
            img_tag = card.find('img')
            image_url = img_tag['src'] if img_tag and 'src' in img_tag.attrs else ''
            
            # Extract alt text
            alt_text = img_tag['alt'] if img_tag and 'alt' in img_tag.attrs else ''
            
            # Download image if URL exists
            local_image_path = ''
            if image_url:
                print(f"Downloading image for: {title}")
                local_image_path = download_image(image_url, title=title)
                if local_image_path:
                    print(f"Successfully downloaded image to: {local_image_path}")
            
            prompt_data.append({
                '标题': title,
                '分类': category,
                '提示词': prompt_text,
                '图片URL': image_url,
                '图片描述': alt_text,
                '本地图片路径': local_image_path
            })
            
        except Exception as e:
            print(f"Error processing card: {e}")
            continue
    
    return prompt_data

def save_to_excel(data, output_file='prompts_data.xlsx'):
    """Save data to Excel file"""
    df = pd.DataFrame(data)
    
    # Create Excel writer with formatting
    with pd.ExcelWriter(output_file, engine='openpyxl') as writer:
        df.to_excel(writer, sheet_name='提示词数据', index=False)
        
        # Get the workbook and worksheet
        workbook = writer.book
        worksheet = writer.sheets['提示词数据']
        
        # Adjust column widths
        for column in worksheet.columns:
            max_length = 0
            column_letter = column[0].column_letter
            
            for cell in column:
                try:
                    if len(str(cell.value)) > max_length:
                        max_length = len(str(cell.value))
                except:
                    pass
            
            # Set minimum width of 10 and maximum of 50
            adjusted_width = min(max(max_length + 2, 10), 50)
            worksheet.column_dimensions[column_letter].width = adjusted_width
    
    print(f"Data saved to {output_file}")
    print(f"Total prompts extracted: {len(data)}")

def main():
    """Main function"""
    html_file = 'docs/coze-prompt.html'
    
    if not os.path.exists(html_file):
        print(f"Error: {html_file} not found!")
        return
    
    print("Extracting prompt data from HTML...")
    prompt_data = extract_prompt_data(html_file)
    
    if prompt_data:
        save_to_excel(prompt_data)
        
        # Display summary statistics
        df = pd.DataFrame(prompt_data)
        print("\n=== Summary Statistics ===")
        print(f"Total prompts: {len(df)}")
        print(f"Categories: {df['分类'].nunique()}")
        print(f"Images downloaded: {df['本地图片路径'].str.len().gt(0).sum()}")
        print("\nTop categories:")
        print(df['分类'].value_counts().head(10))
    else:
        print("No data extracted!")

if __name__ == "__main__":
    main()