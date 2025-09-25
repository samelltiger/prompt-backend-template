
参考  scripts/upload_images_to_oss.py 文件最终生成的excel文件： data/prompts_data_with_oss.xlsx

excel内容示例：
标题	分类	提示词	图片URL	图片描述	本地图片路径	oss短链
Q版3D人物求婚场景	Q版3D	将照片里的两个人转换成Q版 3D人物，场景换成求婚，背景换成淡雅五彩花瓣做的拱门，背景换成浪漫颜色，地上散落着玫瑰花瓣。除了人物采用Q版 3D人物风格，其他环境采用真实写实风格。	https://raw.githubusercontent.com/JimmyLv/awesome-nano-banana/main/cases/1/example_proposal_scene_q_realistic.png	Q版3D人物求婚场景	images\example_proposal_scene_q_realistic_ee3b1ae0.png	images/20250906/20250906_162409_4845c6c7-e76d-4d1d-960d-d8c0e5a3d22f.png
拍立得照片突破效果	Q版3D	将场景中的角色转化为3D Q版风格，放在一张拍立得照片上，相纸被一只手拿着，照片中的角色正从拍立得照片中走出，呈现出突破二维相片边框、进入二维现实空间的视觉效果。	https://raw.githubusercontent.com/JimmyLv/awesome-nano-banana/main/cases/2/example_polaroid_breakout.png	拍立得照片突破效果	images\example_polaroid_breakout_041a7172.png	images/20250906/20250906_162409_edb6a202-6aab-423f-b4e0-f64b8ce64c5e.png
复古宣传海报	海报设计	复古宣传海报风格，突出中文文字，背景为红黄放射状图案。画面中心位置有一位美丽的年轻女性，以精致复古风格绘制，面带微笑，气质优雅，具有亲和力。主题是GPT最新AI绘画服务的广告促销，强调‘惊爆价9.9/张’、‘适用各种场景、图像融合、局部重绘’、‘每张提交3次修改’、‘AI直出效果，无需修改’，底部醒目标注‘有意向点右下“我想要”’，右下角绘制一个手指点击按钮动作，左下角展示OpenAI标志。	https://raw.githubusercontent.com/JimmyLv/awesome-nano-banana/main/cases/3/example_vintage_poster.png	复古宣传海报	images\example_vintage_poster_40bc8b7e.png	images/20250906/20250906_162410_000febc7-bcc2-4097-9f43-0e20230cbf0d.png

需求1：
你需要先帮我设计用于保存上述数据的数据表（ddl保存到一个sql文件中）
包括分类表和提示词表
注意提示词表中的保存oss短链的字段设计成json类型，保存json数组字符串（未来会有多张图片的情况）
提示词表必须包含的字段：
标题	分类id	提示词	图片描述	oss短链

需求2：
在设计的数据表的基础上帮我提供一个接口，用于我创建提示词记录（支持批量添加）
标题	分类名称	提示词	图片描述	oss短链（数组）

这个接口会判断分类名称是否在分类表中有，如果没有，则会先创建

接口可以参考 internal/api/admin/upload.go 文件
还需要一个list接口（这个接口不需要验证授权就能直接访问），并且要将oss短链进行签名（面向用户端）


需求3：
参考  scripts/upload_images_to_oss.py 文件上传oss的方法，帮我开发一个脚本用于读取excel记录，并调用需求2开发的接口，完成数据添加
