// API 基础 URL
const BASE_URL = '';

// 获取存储的 API Key
function getApiKey() {
  return localStorage.getItem('api_key');
}

// 设置 API Key
function setApiKey(apiKey) {
  localStorage.setItem('api_key', apiKey);
}

// 清除 API Key
function clearApiKey() {
  localStorage.removeItem('api_key');
}

// 检查是否已登录
function isLoggedIn() {
  return !!getApiKey();
}

// 通用 API 请求函数
async function apiRequest(endpoint, options = {}) {
  const apiKey = getApiKey();

  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
    }
  };

  if (apiKey) {
    defaultOptions.headers['Authorization'] = `Bearer ${apiKey}`;
  }

  const mergedOptions = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers
    }
  };

  try {
    const response = await fetch(`${BASE_URL}${endpoint}`, mergedOptions);
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || '请求失败');
    }

    return data;
  } catch (error) {
    console.error('API 请求错误:', error);
    throw error;
  }
}

// 获取订阅信息
async function getSubscription() {
  return apiRequest('/v1/dashboard/billing/subscription');
}

// 获取使用情况
async function getUsage() {
  return apiRequest('/v1/dashboard/billing/usage');
}

// 获取调用日志
async function getLogs(page = 1, pageSize = 20) {
  return apiRequest(`/api/logs?page=${page}&page_size=${pageSize}`);
}

// 获取价格信息
async function getPricing() {
  return apiRequest('/api/pricing');
}

// 获取健康状态
async function getStatus() {
  return apiRequest('/api/about');
}

// 兑换码信息查询
async function getRedeemInfo(code) {
  return apiRequest(`/api/redeem?code=${code}`);
}

// 兑换码兑换
async function redeemCode(code) {
  return apiRequest('/api/redeem', {
    method: 'POST',
    body: JSON.stringify({ code })
  });
}

// 更新导航栏状态
function updateNavigation() {
  const isLoggedIn = !!getApiKey();
  const loginSection = document.getElementById('login-section');
  const userSection = document.getElementById('user-section');
  const navLinks = document.querySelectorAll('.nav-links a');

  if (loginSection && userSection) {
    if (isLoggedIn) {
      loginSection.classList.add('hidden');
      userSection.classList.remove('hidden');
    } else {
      loginSection.classList.remove('hidden');
      userSection.classList.add('hidden');
    }
  }

  // 激活当前页面的导航链接
  const splitPath = window.location.pathname.split("/")
  const currentPath = splitPath[splitPath.length-1];
  navLinks.forEach(link => {
    if (link.getAttribute('href') === currentPath) {
      link.classList.add('active');
    } else if (currentPath === "" && link.getAttribute('href').indexOf("index.html")!==-1) {
      link.classList.add('active');
    } else {
      link.classList.remove('active');
    }
  });
}

// 格式化日期
function formatDate(timestamp) {
  const date = new Date(timestamp * 1000);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
}

// 格式化金额
function formatAmount(amount) {
  return `$${parseFloat(amount).toFixed(6)}`;
}

// 登出
function logout() {
  clearApiKey();
  updateNavigation();
  window.location.href = '/static/index.html';
}