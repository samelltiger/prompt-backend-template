我想要添加一个资料库页面，会放一些收集到的资料，如百度网盘地址、网站url、b站视频链接，文章链接等等。帮我设计几个页面原型出来，供我挑选。
支持分类：涛哥资料、AI副业大全、AI编程、AI提示词、办公软件  （这个分类就不要设计新的表了，直接写在代码里）

资源库表：
需要存放收集的资料数据，如百度网盘地址、网站url、b站视频链接，文章链接等等。支持自定义排序和置顶操作，如果排序字段没有设置，那么按时间的倒叙排列。
支持点赞，用户每次点击某个资源，需要记录点击量。
部分资源需要支持用户添加微信领取，用户点击查看后，弹出对应的微信二维码，告知用户添加微信领取。
搜索功能支持搜索标题和标签。

你需要先设计出响应的数据表的sql，在编写go代码。

后端接口参考：internal\api\prompts.go   internal\api\admin\redemption.go 等文件
模型定义参考： internal\model\gateway_models.go

前端展示参考，需根据上面的需求做对应修改。
<div class="container">
        <header>
            <h1>我的资料库</h1>
            <p class="subtitle">收集整理各类学习资源与参考资料</p>
        </header>
        
        <div class="layout-selector">
            <button class="layout-btn active" data-layout="card">卡片式布局</button>
            <button class="layout-btn" data-layout="list">列表式布局</button>
            <button class="layout-btn" data-layout="grid">网格式布局</button>
        </div>
        
        <div class="search-bar">
            <input type="text" placeholder="搜索资料...">
            <button>搜索</button>
        </div>
        
        <div class="category-filter">
            <button class="category-btn active">全部</button>
            <button class="category-btn">涛哥资料</button>
            <button class="category-btn">AI副业大全</button>
            <button class="category-btn">AI编程</button>
            <button class="category-btn">办公软件</button>
        </div>
        
        <!-- 卡片式布局 -->
        <div class="card-layout" id="cardLayout">
            <!-- 示例卡片1 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon baidu">度</div>
                    <div class="card-title">前端开发学习资料包</div>
                </div>
                <div class="card-body">
                    <p class="card-description">包含HTML、CSS、JavaScript进阶教程、框架学习资料和实战项目源码。</p>
                    <div class="card-tags">
                        <span class="tag">前端</span>
                        <span class="tag">JavaScript</span>
                        <span class="tag">教程</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">查看资源</a>
                    <span class="card-date">2023-10-15</span>
                </div>
            </div>
            
            <!-- 示例卡片2 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon website">网</div>
                    <div class="card-title">MDN Web文档</div>
                </div>
                <div class="card-body">
                    <p class="card-description">最权威的Web技术文档，包含HTML、CSS、JavaScript等前端技术的详细参考。</p>
                    <div class="card-tags">
                        <span class="tag">文档</span>
                        <span class="tag">参考</span>
                        <span class="tag">前端</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">访问网站</a>
                    <span class="card-date">2023-09-28</span>
                </div>
            </div>
            
            <!-- 示例卡片3 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon bilibili">B</div>
                    <div class="card-title">CSS布局实战教程</div>
                </div>
                <div class="card-body">
                    <p class="card-description">B站优质UP主讲解的CSS Flexbox和Grid布局实战教程，包含大量实例。</p>
                    <div class="card-tags">
                        <span class="tag">CSS</span>
                        <span class="tag">布局</span>
                        <span class="tag">视频</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">观看视频</a>
                    <span class="card-date">2023-11-05</span>
                </div>
            </div>
            
            <!-- 示例卡片4 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon article">文</div>
                    <div class="card-title">响应式设计最佳实践</div>
                </div>
                <div class="card-body">
                    <p class="card-description">深入探讨响应式网页设计的原则、技巧和常见问题解决方案。</p>
                    <div class="card-tags">
                        <span class="tag">响应式</span>
                        <span class="tag">设计</span>
                        <span class="tag">最佳实践</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">阅读文章</a>
                    <span class="card-date">2023-10-22</span>
                </div>
            </div>
            
            <!-- 示例卡片5 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon book">书</div>
                    <div class="card-title">JavaScript高级程序设计</div>
                </div>
                <div class="card-body">
                    <p class="card-description">前端开发者必读的JavaScript经典书籍，涵盖ES6+新特性和高级概念。</p>
                    <div class="card-tags">
                        <span class="tag">JavaScript</span>
                        <span class="tag">书籍</span>
                        <span class="tag">进阶</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">下载电子书</a>
                    <span class="card-date">2023-09-15</span>
                </div>
            </div>
            
            <!-- 示例卡片6 -->
            <div class="resource-card">
                <div class="card-header">
                    <div class="card-icon website">网</div>
                    <div class="card-title">CSS-Tricks技巧库</div>
                </div>
                <div class="card-body">
                    <p class="card-description">收集了大量CSS实用技巧和解决方案的网站，适合前端开发者日常参考。</p>
                    <div class="card-tags">
                        <span class="tag">CSS</span>
                        <span class="tag">技巧</span>
                        <span class="tag">参考</span>
                    </div>
                </div>
                <div class="card-footer">
                    <a href="#" class="card-link">访问网站</a>
                    <span class="card-date">2023-11-10</span>
                </div>
            </div>
        </div>
        
        <!-- 列表式布局 -->
        <div class="list-layout" id="listLayout">
            <div class="resource-list">
                <!-- 示例列表项1 -->
                <div class="list-item">
                    <div class="list-icon baidu">度</div>
                    <div class="list-content">
                        <div class="list-title">前端开发学习资料包</div>
                        <p class="list-description">包含HTML、CSS、JavaScript进阶教程、框架学习资料和实战项目源码。</p>
                        <div class="list-tags">
                            <span class="tag">前端</span>
                            <span class="tag">JavaScript</span>
                            <span class="tag">教程</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">查看资源</a>
                            <span class="card-date">2023-10-15</span>
                        </div>
                    </div>
                </div>
                
                <!-- 示例列表项2 -->
                <div class="list-item">
                    <div class="list-icon website">网</div>
                    <div class="list-content">
                        <div class="list-title">MDN Web文档</div>
                        <p class="list-description">最权威的Web技术文档，包含HTML、CSS、JavaScript等前端技术的详细参考。</p>
                        <div class="list-tags">
                            <span class="tag">文档</span>
                            <span class="tag">参考</span>
                            <span class="tag">前端</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">访问网站</a>
                            <span class="card-date">2023-09-28</span>
                        </div>
                    </div>
                </div>
                
                <!-- 示例列表项3 -->
                <div class="list-item">
                    <div class="list-icon bilibili">B</div>
                    <div class="list-content">
                        <div class="list-title">CSS布局实战教程</div>
                        <p class="list-description">B站优质UP主讲解的CSS Flexbox和Grid布局实战教程，包含大量实例。</p>
                        <div class="list-tags">
                            <span class="tag">CSS</span>
                            <span class="tag">布局</span>
                            <span class="tag">视频</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">观看视频</a>
                            <span class="card-date">2023-11-05</span>
                        </div>
                    </div>
                </div>
                
                <!-- 示例列表项4 -->
                <div class="list-item">
                    <div class="list-icon article">文</div>
                    <div class="list-content">
                        <div class="list-title">响应式设计最佳实践</div>
                        <p class="list-description">深入探讨响应式网页设计的原则、技巧和常见问题解决方案。</p>
                        <div class="list-tags">
                            <span class="tag">响应式</span>
                            <span class="tag">设计</span>
                            <span class="tag">最佳实践</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">阅读文章</a>
                            <span class="card-date">2023-10-22</span>
                        </div>
                    </div>
                </div>
                
                <!-- 示例列表项5 -->
                <div class="list-item">
                    <div class="list-icon book">书</div>
                    <div class="list-content">
                        <div class="list-title">JavaScript高级程序设计</div>
                        <p class="list-description">前端开发者必读的JavaScript经典书籍，涵盖ES6+新特性和高级概念。</p>
                        <div class="list-tags">
                            <span class="tag">JavaScript</span>
                            <span class="tag">书籍</span>
                            <span class="tag">进阶</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">下载电子书</a>
                            <span class="card-date">2023-09-15</span>
                        </div>
                    </div>
                </div>
                
                <!-- 示例列表项6 -->
                <div class="list-item">
                    <div class="list-icon website">网</div>
                    <div class="list-content">
                        <div class="list-title">CSS-Tricks技巧库</div>
                        <p class="list-description">收集了大量CSS实用技巧和解决方案的网站，适合前端开发者日常参考。</p>
                        <div class="list-tags">
                            <span class="tag">CSS</span>
                            <span class="tag">技巧</span>
                            <span class="tag">参考</span>
                        </div>
                        <div class="list-footer">
                            <a href="#" class="card-link">访问网站</a>
                            <span class="card-date">2023-11-10</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        
        <!-- 网格式布局 -->
        <div class="grid-layout" id="gridLayout">
            <!-- 示例网格项1 -->
            <div class="resource-grid">
                <div class="grid-icon baidu">度</div>
                <div class="grid-content">
                    <div class="grid-title">前端开发学习资料包</div>
                    <p class="grid-description">包含HTML、CSS、JavaScript进阶教程、框架学习资料和实战项目源码。</p>
                    <div class="grid-tags">
                        <span class="tag">前端</span>
                        <span class="tag">JavaScript</span>
                    </div>
                    <div class="grid-footer">
                        <a href="#" class="card-link">查看资源</a>
                        <span class="card-date">2023-10-15</span>
                    </div>
                </div>
            </div>
            
            <!-- 示例网格项2 -->
            <div class="resource-grid">
                <div class="grid-icon website">网</div>
                <div class="grid-content">
                    <div class="grid-title">MDN Web文档</div>
                    <p class="grid-description">最权威的Web技术文档，包含HTML、CSS、JavaScript等前端技术的详细参考。</p>
                    <div class="grid-tags">
                        <span class="tag">文档</span>
                        <span class="tag">参考</span>
                    </div>
                    <div class="grid-footer">
                        <a href="#" class="card-link">访问网站</a>
                        <span class="card-date">2023-09-28</span>
                    </div>
                </div>
            </div>
            
            <!-- 示例网格项3 -->
            <div class="resource-grid">
                <div class="grid-icon bilibili">B</div>
                <div class="grid-content">
                    <div class="grid-title">CSS布局实战教程</div>
                    <p class="grid-description">B站优质UP主讲解的CSS Flexbox和Grid布局实战教程，包含大量实例。</p>
                    <div class="grid-tags">
                        <span class="tag">CSS</span>
                        <span class="tag">布局</span>
                    </div>
                    <div class="grid-footer">
                        <a href="#" class="card-link">观看视频</a>
                        <span class="card-date">2023-11-05</span>
                    </div>
                </div>
            </div>
            
            <!-- 示例网格项4 -->
            <div class="resource-grid">
                <div class="grid-icon article">文</div>
                <div class="grid-content">
                    <div class="grid-title">响应式设计最佳实践</div>
                    <p class="grid-description">深入探讨响应式网页设计的原则、技巧和常见问题解决方案。</p>
                    <div class="grid-tags">
                        <span class="tag">响应式</span>
                        <span class="tag">设计</span>
                    </div>
                    <div class="grid-footer">
                        <a href="#" class="card-link">阅读文章</a>
                        <span class="card-date">2023-10-22</span>
                    </div>
                </div>
            </div>
            
        </div>
    </div>