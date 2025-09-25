document.addEventListener('DOMContentLoaded', function () {
    // 初始化页面
    updateNavigation();
    checkLoginStatus();

    // 绑定事件
    document.getElementById('api-key-form').addEventListener('submit', handleApiKeySubmit);
    document.getElementById('logout-btn').addEventListener('click', logout);
    document.getElementById('check-btn').addEventListener('click', checkRedeemCode);
    document.getElementById('redeem-form').addEventListener('submit', handleRedeemCode);
});

// 检查登录状态
function checkLoginStatus() {
    if (isLoggedIn()) {
        document.getElementById('not-logged-in').classList.add('hidden');
        document.getElementById('logged-in').classList.remove('hidden');
        document.getElementById('user-api-key').textContent = maskApiKey(getApiKey());
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

// 查询兑换码
async function checkRedeemCode() {
    const codeInput = document.getElementById('redeem-code');
    const code = codeInput.value.trim();

    if (!code) {
        showResult('请输入有效的兑换码', 'danger');
        return;
    }

    try {
        // 隐藏之前的结果
        hideResult();
        document.getElementById('code-info').classList.add('hidden');

        // 查询兑换码
        const response = await getRedeemInfo(code);
        const data = response.data;

        // 显示兑换码信息
        document.getElementById('code-quota').textContent = data.quota;
        document.getElementById('code-amount').textContent = `$${data.amount}`;
        document.getElementById('code-info').classList.remove('hidden');

        showResult('兑换码有效，可以兑换', 'success');
    } catch (error) {
        showResult(`兑换码查询失败: ${error.message}`, 'danger');
    }
}

// 兑换码兑换
async function handleRedeemCode(e) {
    e.preventDefault();
    const codeInput = document.getElementById('redeem-code');
    const code = codeInput.value.trim();

    if (!code) {
        showResult('请输入有效的兑换码', 'danger');
        return;
    }

    try {
        // 隐藏之前的结果
        hideResult();

        // 兑换码
        const response = await redeemCode(code);
        const data = response.data;

        // 显示兑换结果
        showResult(`兑换成功！已添加 $${data.amount} 额度到您的账户`, 'success');

        // 清空输入框
        codeInput.value = '';

        // 隐藏兑换码信息
        document.getElementById('code-info').classList.add('hidden');
    } catch (error) {
        showResult(`兑换失败: ${error.message}`, 'danger');
    }
}

// 显示结果
function showResult(message, type) {
    const resultElement = document.getElementById('redeem-result');
    resultElement.textContent = message;
    resultElement.className = `alert alert-${type}`;
    resultElement.classList.remove('hidden');
}

// 隐藏结果
function hideResult() {
    const resultElement = document.getElementById('redeem-result');
    resultElement.classList.add('hidden');
}

// 掩码 API Key
function maskApiKey(apiKey) {
    if (!apiKey || apiKey.length < 8) return '***';
    return apiKey.substring(0, 4) + '...' + apiKey.substring(apiKey.length - 4);
}