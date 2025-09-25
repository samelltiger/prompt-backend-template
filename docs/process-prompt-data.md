
根据  docs/coze-prompt.html 的数据帮我提取出数据，具体如下：

内容示例：
        <div class="bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300 flex flex-col h-full"
            style="opacity: 1; transform: none;">
            <div class="relative overflow-hidden aspect-video bg-gray-100"><img
                    src="https://raw.githubusercontent.com/JimmyLv/awesome-nano-banana/main/cases/1/example_proposal_scene_q_realistic.png"
                    alt="Q版3D人物求婚场景"
                    class="w-full h-full object-cover transition-transform duration-500 hover:scale-105" loading="lazy">
                <div class="absolute top-3 right-3"><span
                        class="bg-blue-500 text-white text-xs font-bold px-2.5 py-0.5 rounded-full">Q版3D</span></div>
            </div>
            <div class="p-5 flex flex-col flex-grow">
                <h3 class="text-lg font-bold text-gray-800 mb-2 line-clamp-1">Q版3D人物求婚场景</h3>
                <div class="mt-3 bg-gray-50 rounded-lg p-4 text-sm text-gray-700 font-mono overflow-x-auto">
                    <p class="whitespace-pre-wrap line-clamp-3">将照片里的两个人转换成Q版
                        3D人物，场景换成求婚，背景换成淡雅五彩花瓣做的拱门，背景换成浪漫颜色，地上散落着玫瑰花瓣。除了人物采用Q版 3D人物风格，其他环境采用真实写实风格。</p>
                </div>
                <div class="mt-4 flex justify-end mt-auto"><button
                        class="flex items-center px-3 py-1.5 rounded-lg text-sm font-medium transition-colors duration-200 bg-blue-50 text-blue-600 hover:bg-blue-100">复制提示词
                        <i class="fa-regular fa-clipboard ml-1.5"></i></button></div>
            </div>
        </div>
        <div class="bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow duration-300 flex flex-col h-full"
            style="opacity: 1; transform: none;">
            <div class="relative overflow-hidden aspect-video bg-gray-100"><img
                    src="https://raw.githubusercontent.com/JimmyLv/awesome-nano-banana/main/cases/2/example_polaroid_breakout.png"
                    alt="拍立得照片突破效果" class="w-full h-full object-cover transition-transform duration-500 hover:scale-105"
                    loading="lazy">
                <div class="absolute top-3 right-3"><span
                        class="bg-blue-500 text-white text-xs font-bold px-2.5 py-0.5 rounded-full">Q版3D</span></div>
            </div>
            <div class="p-5 flex flex-col flex-grow">
                <h3 class="text-lg font-bold text-gray-800 mb-2 line-clamp-1">拍立得照片突破效果</h3>
                <div class="mt-3 bg-gray-50 rounded-lg p-4 text-sm text-gray-700 font-mono overflow-x-auto">
                    <p class="whitespace-pre-wrap line-clamp-3">将场景中的角色转化为3D
                        Q版风格，放在一张拍立得照片上，相纸被一只手拿着，照片中的角色正从拍立得照片中走出，呈现出突破二维相片边框、进入二维现实空间的视觉效果。</p>
                </div>
                <div class="mt-4 flex justify-end mt-auto"><button
                        class="flex items-center px-3 py-1.5 rounded-lg text-sm font-medium transition-colors duration-200 bg-blue-50 text-blue-600 hover:bg-blue-100">复制提示词
                        <i class="fa-regular fa-clipboard ml-1.5"></i></button></div>
            </div>
        </div>


我要做一个提示词网站，请帮我把上面的数据整理出来，提取出来成为一个excle文件（使用pandas），代码逻辑保存到一个python脚本中