let modelData = [];

document.addEventListener('DOMContentLoaded', function () {
    // 初始化页面
    updateNavigation();

    // 绑定事件
    document.getElementById('api-key-form').addEventListener('submit', handleApiKeySubmit);
    document.getElementById('logout-btn').addEventListener('click', logout);
    document.getElementById('model-search').addEventListener('input', filterModels);

    // 加载价格数据
    loadPricingData();
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

// 加载价格数据
async function loadPricingData() {
    try {
        // 显示加载状态
        document.getElementById('pricing-loading').classList.remove('hidden');
        document.getElementById('pricing-content').classList.add('hidden');

        const response = await getPricing();
        const data = response.data;

        if (data && data.data) {
            modelData = data.data;
            renderPricingTable(modelData);
        } else {
            throw new Error('获取价格数据失败');
        }

        // 隐藏加载状态
        document.getElementById('pricing-loading').classList.add('hidden');
        document.getElementById('pricing-content').classList.remove('hidden');

    } catch (error) {
        console.error('加载价格数据失败:', error);
        document.getElementById('pricing-loading').classList.add('hidden');

        // 显示错误信息
        const tableBody = document.getElementById('pricing-table-body');
        tableBody.innerHTML = `<tr><td colspan="5" style="text-align: center; color: red;">加载价格数据失败: ${error.message}</td></tr>`;
        document.getElementById('pricing-content').classList.remove('hidden');
    }
}

// 渲染价格表格
function renderPricingTable(models) {
    const tableBody = document.getElementById('pricing-table-body');
    tableBody.innerHTML = '';

    if (models && models.length > 0) {
        models.forEach(model => {
            const row = document.createElement('tr');

            // 可用性
            const availabilityCell = document.createElement('td');
            availabilityCell.innerHTML = model.enable_groups && model.enable_groups.includes('default')
                ? '<span class="badge badge-success">✅</span>'
                : '<span class="badge badge-danger">❌</span>';

            // 模型名称
            const nameCell = document.createElement('td');
            nameCell.textContent = model.model_name;

            // 计费类型
            const billingTypeCell = document.createElement('td');
            billingTypeCell.textContent = model.quota_type === 0 ? '按量计费' : '包月计费';

            // 倍率
            const ratioCell = document.createElement('td');
            ratioCell.innerHTML = `模型倍率：${model.model_ratio}<br>补全倍率：${model.completion_ratio}`;

            // 模型价格
            const priceCell = document.createElement('td');
            // 计算价格 (基础价格为 $2 / 1M tokens)
            const promptPrice = (model.model_ratio * 2).toFixed(6);
            const completionPrice = (model.model_ratio * model.completion_ratio * 2).toFixed(6);
            priceCell.innerHTML = `提示 $${promptPrice} / 1M tokens<br>补全 $${completionPrice} / 1M tokens`;

            row.appendChild(availabilityCell);
            row.appendChild(nameCell);
            row.appendChild(billingTypeCell);
            row.appendChild(ratioCell);
            row.appendChild(priceCell);

            tableBody.appendChild(row);
        });
    } else {
        const row = document.createElement('tr');
        row.innerHTML = `<td colspan="5" style="text-align: center;">暂无模型价格数据</td>`;
        tableBody.appendChild(row);
    }
}

// 过滤模型
function filterModels() {
    const searchTerm = document.getElementById('model-search').value.toLowerCase();

    if (!searchTerm) {
        renderPricingTable(modelData);
        return;
    }

    const filteredModels = modelData.filter(model =>
        model.model_name.toLowerCase().includes(searchTerm)
    );

    renderPricingTable(filteredModels);
}

// 掩码 API Key
function maskApiKey(apiKey) {
    if (!apiKey || apiKey.length < 8) return '***';
    return apiKey.substring(0, 4) + '...' + apiKey.substring(apiKey.length - 4);
}