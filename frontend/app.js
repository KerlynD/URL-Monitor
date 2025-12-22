const API_BASE_URL = 'http://localhost:8080';
const REFRESH_INTERVAL = 30000; 

const addMonitorForm = document.getElementById('addMonitorForm');
const urlInput = document.getElementById('urlInput');
const intervalInput = document.getElementById('intervalInput');
const monitorList = document.getElementById('monitorList');

document.addEventListener('DOMContentLoaded', () => {
    loadMonitors();
    
    // Auto-refresh monitors
    setInterval(loadMonitors, REFRESH_INTERVAL);
    
    // Form submission
    addMonitorForm.addEventListener('submit', handleAddMonitor);
});

// Load all monitors from API
async function loadMonitors() {
    try {
        const response = await fetch(`${API_BASE_URL}/monitor`);
        
        if (!response.ok) {
            throw new Error('Failed to fetch monitors');
        }
        
        const monitors = await response.json();
        displayMonitors(monitors);
    } catch (error) {
        console.error('Error loading monitors:', error);
        monitorList.innerHTML = `
            <div class="error-message">
                Failed to load monitors. Make sure the backend is running on ${API_BASE_URL}
            </div>
        `;
    }
}

// Display monitors in the UI
function displayMonitors(monitors) {
    if (!monitors || monitors.length === 0) {
        monitorList.innerHTML = `
            <div class="empty-state">
                <p>📭 No monitors yet</p>
                <p>Add your first URL to start monitoring!</p>
            </div>
        `;
        return;
    }
    
    monitorList.innerHTML = monitors.map(monitor => {
        const status = getStatus(monitor.last_result);
        const responseTime = getResponseTime(monitor.last_result);
        const lastChecked = getLastChecked(monitor.last_result);
        
        return `
            <div class="monitor-card ${status}">
                <div class="monitor-header">
                    <div class="monitor-url">${escapeHtml(monitor.url)}</div>
                    <span class="status-badge ${status}">${status}</span>
                </div>
                
                <div class="monitor-details">
                    <div class="detail-item">
                        <span class="detail-label">Response Time</span>
                        <span class="detail-value">${responseTime}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Check Interval</span>
                        <span class="detail-value">${monitor.check_interval}s</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Last Checked</span>
                        <span class="detail-value">${lastChecked}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Status Code</span>
                        <span class="detail-value">${monitor.last_result?.status_code || 'N/A'}</span>
                    </div>
                </div>
                
                ${monitor.last_result?.error ? `
                    <div class="error-message">
                        ${escapeHtml(monitor.last_result.error)}
                    </div>
                ` : ''}
                
                <div class="monitor-actions">
                    <button class="btn btn-check" onclick="triggerCheck('${monitor.id}')">
                        Check Now
                    </button>
                </div>
            </div>
        `;
    }).join('');
}

// Handle adding a new monitor
async function handleAddMonitor(e) {
    e.preventDefault();
    
    const url = urlInput.value.trim();
    const checkInterval = parseInt(intervalInput.value);
    
    if (!url || !checkInterval) {
        alert('Please fill in all fields');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/monitor`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                url: url,
                check_interval: checkInterval
            })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to add monitor');
        }
        
        urlInput.value = '';
        intervalInput.value = '60';
        
        showMessage('Monitor added successfully!', 'success');
        
        loadMonitors();
    } catch (error) {
        console.error('Error adding monitor:', error);
        showMessage(error.message, 'error');
    }
}

// Trigger a manual check
async function triggerCheck(monitorId) {
    try {
        const response = await fetch(`${API_BASE_URL}/monitor/${monitorId}/check`, {
            method: 'POST'
        });
        
        if (!response.ok) {
            throw new Error('Failed to trigger check');
        }
        
        showMessage('Check completed!', 'success');
        
        setTimeout(loadMonitors, 1000);
    } catch (error) {
        console.error('Error triggering check:', error);
        showMessage('Failed to trigger check', 'error');
    }
}

// Helper: Get status from result
function getStatus(result) {
    if (!result) return 'unknown';
    return result.is_up ? 'up' : 'down';
}

// Helper: Get formatted response time
function getResponseTime(result) {
    if (!result || !result.response_time) return 'N/A';
    
    const ms = result.response_time / 1000000;
    
    if (ms < 1000) {
        return `${Math.round(ms)}ms`;
    } else {
        return `${(ms / 1000).toFixed(2)}s`;
    }
}

// Helper: Get formatted last checked time
function getLastChecked(result) {
    if (!result || !result.timestamp) return 'Never';
    
    const date = new Date(result.timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
}

// Helper: Show temporary message
function showMessage(message, type) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `${type}-message`;
    messageDiv.textContent = message;
    messageDiv.style.position = 'fixed';
    messageDiv.style.top = '20px';
    messageDiv.style.right = '20px';
    messageDiv.style.zIndex = '1000';
    messageDiv.style.minWidth = '250px';
    messageDiv.style.boxShadow = '0 4px 12px rgba(0,0,0,0.15)';
    
    document.body.appendChild(messageDiv);
    
    setTimeout(() => {
        messageDiv.style.opacity = '0';
        messageDiv.style.transition = 'opacity 0.3s';
        setTimeout(() => messageDiv.remove(), 300);
    }, 3000);
}

// Helper: Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

