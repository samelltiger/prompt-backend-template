document.addEventListener('DOMContentLoaded', function() {
    // 初始化页面
    updateNavigation();
    checkLoginStatus();
  
    // 绑定事件
    document.getElementById('api-key-form').addEventListener('submit', handleApiKeySubmit);
    document.getElementById('logout-btn').addEventListener('click', logout);
  });
  
  // 检查登录状态并加载数据
  function checkLoginStatus() {
    if (isLoggedIn()) {
      document.getElementById('not-logged-in').classList.add('hidden');
      document.getElementById('logged-in').classList.remove('hidden');
      document.getElementById('user-api-key').textContent = maskApiKey(getApiKey());
  
      // 加载数据
      loadSubscriptionData();
      loadUsageData();
      loadLogs(1);
    } else {
      document.getElementById('not-logged-in').classList.remove('hidden');
      document.getElementById('logged-in').classList.add('hidden');
    }
  }
  
  // 处理 API Key 提交
  async function handleApiKeySubmit(e) {
    e.preventDefault();
    const apiKeyInput = document.getElementById('api-key-input');
    const apiKey = apiKeyInput.value.trim();
  
    if (!apiKey) {
      alert('请输入有效的 API Key');
      return;
    }
  
    try {
      // 临时存储 API Key
      setApiKey(apiKey);
  
      // 尝试获取订阅信息验证 API Key 有效性
      await getSubscription();
  
      // 更新页面状态
      updateNavigation();
      checkLoginStatus();
    } catch (error) {
      clearApiKey();
      alert(`API Key 无效: ${error.message}`);
    }
  }
  
  // 加载订阅数据
  async function loadSubscriptionData() {
    try {
      const response = await getSubscription();
      const data = response.data;
  
      // 填充数据
      document.getElementById('total-quota').textContent = `$${data.hard_limit_usd.toFixed(3)}`;
  
      // 处理过期时间
      const expiryDate = data.access_until === 0 
        ? '永不过期' 
        : formatDate(data.access_until);
      document.getElementById('expiry-date').textContent = expiryDate;
  
    } catch (error) {
      console.error('加载订阅数据失败:', error);
    }
  }
  
  // 加载使用情况数据
  async function loadUsageData() {
    try {
      const response = await getUsage();
      const data = response.data;
  
      // 计算已使用额度 (total_usage 是美分)
      const usedAmount = (data.total_usage / 100).toFixed(3);
      document.getElementById('used-quota').textContent = `$${usedAmount}`;
  
      // 计算剩余额度 (从订阅中获取总额)
      const subscriptionResponse = await getSubscription();
      const totalAmount = subscriptionResponse.data.hard_limit_usd;
      const remainingAmount = (totalAmount - (data.total_usage / 100)).toFixed(3);
      document.getElementById('remaining-quota').textContent = `$${remainingAmount}`;
  
    } catch (error) {
      console.error('加载使用情况数据失败:', error);
    }
  }
  
  // 加载调用日志
  async function loadLogs(page = 1, pageSize = 20) {
    try {
      // 显示加载状态
      document.getElementById('logs-loading').classList.remove('hidden');
      document.getElementById('logs-content').classList.add('hidden');
  
      const response = await getLogs(page, pageSize);
      const { data, meta } = response.data;
  
      // 填充日志表格
      const tableBody = document.getElementById('logs-table-body');
      tableBody.innerHTML = '';
  
      if (data && data.length > 0) {
        data.forEach(log => {
          // 计算花费 (tokens * price)
          const totalTokens = log.prompt_tokens + log.completion_tokens;
          // 假设 quota 是以 1/100 美分计算的
          const cost = log.quota / 500000;
  
          const row = document.createElement('tr');
          row.innerHTML = `
            <td>${formatDate(log.created_at)}</td>
            <td>${log.model_name}</td>
            <td>${log.use_time}秒</td>
            <td>${log.prompt_tokens}</td>
            <td>${log.completion_tokens}</td>
            <td>${formatAmount(cost)}</td>
            <td>${log.content}</td>
          `;
          tableBody.appendChild(row);
        });
      } else {
        const row = document.createElement('tr');
        row.innerHTML = `<td colspan="7" style="text-align: center;">暂无日志数据</td>`;
        tableBody.appendChild(row);
      }
  
      // 生成分页
      generatePagination(meta.current_page, meta.total_pages);
  
      // 隐藏加载状态
      document.getElementById('logs-loading').classList.add('hidden');
      document.getElementById('logs-content').classList.remove('hidden');
  
    } catch (error) {
      console.error('加载日志失败:', error);
      document.getElementById('logs-loading').classList.add('hidden');
  
      // 显示错误信息
      const tableBody = document.getElementById('logs-table-body');
      tableBody.innerHTML = `<tr><td colspan="7" style="text-align: center; color: red;">加载日志失败: ${error.message}</td></tr>`;
      document.getElementById('logs-content').classList.remove('hidden');
    }
  }
  
  // 生成分页
  function generatePagination(currentPage, totalPages) {
    const paginationElement = document.getElementById('logs-pagination');
    paginationElement.innerHTML = '';
  
    if (totalPages <= 1) {
      return;
    }
  
    // 计算显示的页码范围
    let startPage = Math.max(1, currentPage - 2);
    let endPage = Math.min(totalPages, startPage + 4);
  
    if (endPage - startPage < 4) {
      startPage = Math.max(1, endPage - 4);
    }
  
    // 首页
    if (startPage > 1) {
      addPageLink(paginationElement, 1, '首页');
    }
  
    // 上一页
    if (currentPage > 1) {
      addPageLink(paginationElement, currentPage - 1, '上一页');
    }
  
    // 页码
    for (let i = startPage; i <= endPage; i++) {
      addPageLink(paginationElement, i, i.toString(), i === currentPage);
    }
  
    // 下一页
    if (currentPage < totalPages) {
      addPageLink(paginationElement, currentPage + 1, '下一页');
    }
  
    // 末页
    if (endPage < totalPages) {
      addPageLink(paginationElement, totalPages, '末页');
    }
  }
  
  // 添加分页链接
  function addPageLink(container, page, text, isActive = false) {
    const li = document.createElement('li');
    const a = document.createElement('a');
    a.href = 'javascript:void(0)';
    a.textContent = text;
  
    if (isActive) {
      a.classList.add('active');
    }
  
    a.addEventListener('click', () => loadLogs(page));
  
    li.appendChild(a);
    container.appendChild(li);
  }
  
  // 掩码 API Key
  function maskApiKey(apiKey) {
    if (!apiKey || apiKey.length < 8) return '***';
    return apiKey.substring(0, 4) + '...' + apiKey.substring(apiKey.length - 4);
  }