let modelData = [];

document.addEventListener('DOMContentLoaded', function() {
  // 初始化页面
  updateNavigation();

  // 绑定事件
  document.getElementById('api-key-form').addEventListener('submit', handleApiKeySubmit);
  document.getElementById('logout-btn').addEventListener('click', logout);
  document.getElementById('model-search').addEventListener('input', filterModels);

  // 加载状态数据
  loadStatusData();

  // 加载模型数据
  loadModelsData();
});

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
    document.getElementById('user-api-key').textContent = maskApiKey(getApiKey());
  } catch (error) {
    clearApiKey();
    alert(`API Key 无效: ${error.message}`);
  }
}

// 加载状态数据
async function loadStatusData() {
  try {
    // 显示加载状态
    document.getElementById('status-loading').classList.remove('hidden');
    document.getElementById('status-content').classList.add('hidden');

    const startTime = performance.now();
    const response = await getStatus();
    const endTime = performance.now();
    const responseTime = Math.round(endTime - startTime);

    // 更新状态信息
    document.getElementById('api-status').innerHTML = '<span class="badge badge-success">正常</span>';
    document.getElementById('response-time').textContent = `${responseTime} ms`;

    // 隐藏加载状态
    document.getElementById('status-loading').classList.add('hidden');
    document.getElementById('status-content').classList.remove('hidden');

  } catch (error) {
    console.error('加载状态数据失败:', error);
    document.getElementById('status-loading').classList.add('hidden');

    // 显示错误信息
    document.getElementById('api-status').innerHTML = '<span class="badge badge-danger">异常</span>';
    document.getElementById('response-time').textContent = `-- ms`;
    document.getElementById('status-content').classList.remove('hidden');
  }
}

// 加载模型数据
async function loadModelsData() {
  try {
    // 显示加载状态
    document.getElementById('models-loading').classList.remove('hidden');
    document.getElementById('models-content').classList.add('hidden');

    const response = await getPricing();
    const data = response.data;

    if (data && data.data) {
      modelData = data.data;
      renderModelsTable(modelData);
    } else {
      throw new Error('获取模型数据失败');
    }

    // 隐藏加载状态
    document.getElementById('models-loading').classList.add('hidden');
    document.getElementById('models-content').classList.remove('hidden');

  } catch (error) {
    console.error('加载模型数据失败:', error);
    document.getElementById('models-loading').classList.add('hidden');

    // 显示错误信息
    const tableBody = document.getElementById('models-table-body');
    tableBody.innerHTML = `<tr><td colspan="4" style="text-align: center; color: red;">加载模型数据失败: ${error.message}</td></tr>`;
    document.getElementById('models-content').classList.remove('hidden');
  }
}

// 渲染模型表格
function renderModelsTable(models) {
  const tableBody = document.getElementById('models-table-body');
  tableBody.innerHTML = '';

  if (models && models.length > 0) {
    models.forEach(model => {
      const row = document.createElement('tr');

      // 状态
      const statusCell = document.createElement('td');
      statusCell.innerHTML = model.enable_groups && model.enable_groups.includes('default') 
        ? '<span class="badge badge-success">可用</span>' 
        : '<span class="badge badge-danger">不可用</span>';

      // 模型名称
      const nameCell = document.createElement('td');
      nameCell.textContent = model.model_name;

      // 分组
      const groupCell = document.createElement('td');
      if (model.enable_groups && model.enable_groups.length > 0) {
        groupCell.textContent = model.enable_groups.join(', ');
      } else {
        groupCell.textContent = '无可用分组';
      }

      // 详情
      const detailsCell = document.createElement('td');
      detailsCell.textContent = model.owner_by || '系统模型';

      row.appendChild(statusCell);
      row.appendChild(nameCell);
      row.appendChild(groupCell);
      row.appendChild(detailsCell);

      tableBody.appendChild(row);
    });
  } else {
    const row = document.createElement('tr');
    row.innerHTML = `<td colspan="4" style="text-align: center;">暂无模型数据</td>`;
    tableBody.appendChild(row);
  }
}

// 过滤模型
function filterModels() {
  const searchTerm = document.getElementById('model-search').value.toLowerCase();

  if (!searchTerm) {
    renderModelsTable(modelData);
    return;
  }

  const filteredModels = modelData.filter(model => 
    model.model_name.toLowerCase().includes(searchTerm)
  );

  renderModelsTable(filteredModels);
}

// 掩码 API Key
function maskApiKey(apiKey) {
  if (!apiKey || apiKey.length < 8) return '***';
  return apiKey.substring(0, 4) + '...' + apiKey.substring(apiKey.length - 4);
}