// Global variables
let currentUser = null;
let currentFriend = null;
let ws = null;
let friends = [];

// Theme Management
function initTheme() {
    const savedTheme = localStorage.getItem('chatTheme') || 'dark';
    if (savedTheme === 'light') {
        document.body.classList.add('light-mode');
        updateThemeIcon('light');
    } else {
        updateThemeIcon('dark');
    }
}

function toggleTheme() {
    const body = document.body;
    const isLight = body.classList.contains('light-mode');
    
    if (isLight) {
        body.classList.remove('light-mode');
        localStorage.setItem('chatTheme', 'dark');
        updateThemeIcon('dark');
    } else {
        body.classList.add('light-mode');
        localStorage.setItem('chatTheme', 'light');
        updateThemeIcon('light');
    }
}

function updateThemeIcon(theme) {
    const themeIcon = document.querySelector('.theme-icon');
    if (themeIcon) {
        themeIcon.textContent = theme === 'light' ? '🌙' : '☀️';
    }
}

// Check if user is already logged in
window.onload = function() {
    initTheme();
    const savedUser = localStorage.getItem('chatUser');
    if (savedUser) {
        currentUser = JSON.parse(savedUser);
        showChat();
        loadFriends();
        loadFriendRequests();
        connectWebSocket();
    }
};

// Show/Hide forms
function showRegister() {
    document.getElementById('login-form').style.display = 'none';
    document.getElementById('register-form').style.display = 'block';
    clearErrors();
}

function showLogin() {
    document.getElementById('register-form').style.display = 'none';
    document.getElementById('login-form').style.display = 'block';
    clearErrors();
}

function clearErrors() {
    document.getElementById('login-error').textContent = '';
    document.getElementById('register-error').textContent = '';
}

// Register function
async function register() {
    const username = document.getElementById('register-username').value.trim();
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;
    const errorDiv = document.getElementById('register-error');

    errorDiv.textContent = '';

    if (!username || !password || !confirmPassword) {
        errorDiv.textContent = 'All fields are required';
        return;
    }

    if (password !== confirmPassword) {
        errorDiv.textContent = 'Passwords do not match';
        return;
    }

    if (password.length < 4) {
        errorDiv.textContent = 'Password must be at least 4 characters';
        return;
    }

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (data.success) {
            errorDiv.className = 'success-message';
            errorDiv.textContent = 'Registration successful! Please login.';
            setTimeout(() => {
                showLogin();
                document.getElementById('login-username').value = username;
            }, 1500);
        } else {
            errorDiv.className = 'error-message';
            errorDiv.textContent = data.message;
        }
    } catch (error) {
        errorDiv.className = 'error-message';
        errorDiv.textContent = 'Network error. Please try again.';
        console.error('Register error:', error);
    }
}

// Login function
async function login() {
    const username = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;
    const errorDiv = document.getElementById('login-error');

    errorDiv.textContent = '';

    if (!username || !password) {
        errorDiv.textContent = 'Username and password are required';
        return;
    }

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (data.success) {
            currentUser = data.data;
            localStorage.setItem('chatUser', JSON.stringify(currentUser));
            showChat();
            loadFriends();
            loadFriendRequests();
            connectWebSocket();
        } else {
            errorDiv.textContent = data.message;
        }
    } catch (error) {
        errorDiv.textContent = 'Network error. Please try again.';
        console.error('Login error:', error);
    }
}

// Logout function
function logout() {
    currentUser = null;
    currentFriend = null;
    localStorage.removeItem('chatUser');
    if (ws) {
        ws.close();
        ws = null;
    }
    document.getElementById('chat-container').style.display = 'none';
    document.getElementById('auth-container').style.display = 'flex';
    document.getElementById('messages').innerHTML = '';
    document.getElementById('message-input').value = '';
}

// Show chat interface
function showChat() {
    document.getElementById('auth-container').style.display = 'none';
    document.getElementById('chat-container').style.display = 'flex';
    document.getElementById('current-username').textContent = currentUser.username;
}

// Load friends list
async function loadFriends() {
    try {
        console.log('Loading friends for user:', currentUser.user_id);
        const response = await fetch(`/api/friends?user_id=${currentUser.user_id}`);
        const data = await response.json();
        
        console.log('Friends API response:', data);

        if (data.success && data.data) {
            friends = data.data;
            console.log('Friends loaded:', friends.length, 'friends');
            displayFriends();
        }
    } catch (error) {
        console.error('Error loading friends:', error);
    }
}

// Display friends in sidebar
function displayFriends() {
    const friendsList = document.getElementById('friends-list');
    
    if (friends.length === 0) {
        friendsList.innerHTML = `
            <div class="empty-state">
                <p style="font-size: 16px; margin-bottom: 10px;">No friends yet!</p>
                <p style="font-size: 14px; color: #666;">
                    Click the <strong>➕</strong> button above to add friends.<br><br>
                    After sending a request, wait for them to accept it.
                </p>
            </div>
        `;
        return;
    }

    friendsList.innerHTML = '';
    friends.forEach(friend => {
        const friendItem = document.createElement('div');
        friendItem.className = 'friend-item';
        if (currentFriend && currentFriend.id === friend.id) {
            friendItem.classList.add('active');
        }
        
        friendItem.innerHTML = `
            <div class="friend-avatar">
                <svg viewBox="0 0 24 24" width="28" height="28">
                    <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                </svg>
            </div>
            <div class="friend-info">
                <span class="friend-name">${escapeHtml(friend.username)}</span>
            </div>
            ${friend.unread_count > 0 ? `<span class="unread-badge">${friend.unread_count}</span>` : ''}
        `;
        
        friendItem.onclick = () => selectFriend(friend);
        friendsList.appendChild(friendItem);
    });
}

// Select a friend to chat with
function selectFriend(friend) {
    currentFriend = friend;
    displayFriends();
    
    document.getElementById('no-chat-selected').style.display = 'none';
    document.getElementById('chat-area').style.display = 'flex';
    document.getElementById('chat-friend-name').textContent = friend.username;
    
    loadMessages();
}

// Load messages with a friend
async function loadMessages() {
    if (!currentFriend) return;

    try {
        const response = await fetch(`/api/messages/${currentFriend.id}?user_id=${currentUser.user_id}`);
        const data = await response.json();

        if (data.success && data.data) {
            const messagesDiv = document.getElementById('messages');
            messagesDiv.innerHTML = '';
            data.data.forEach(message => {
                displayMessage(message);
            });
            scrollToBottom();
            
            // Refresh friends list to update unread count
            loadFriends();
        }
    } catch (error) {
        console.error('Error loading messages:', error);
    }
}

// Send message
async function sendMessage() {
    if (!currentFriend) return;

    const messageInput = document.getElementById('message-input');
    const message = messageInput.value.trim();

    if (!message) return;

    try {
        const response = await fetch('/api/messages', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                sender_id: currentUser.user_id,
                recipient_id: currentFriend.id,
                message: message,
            }),
        });

        const data = await response.json();

        if (data.success) {
            messageInput.value = '';
            displayMessage(data.data);
            scrollToBottom();
        } else {
            alert('Error sending message: ' + data.message);
        }
    } catch (error) {
        console.error('Error sending message:', error);
        alert('Network error. Please try again.');
    }
}

// Handle Enter key press
function handleKeyPress(event) {
    if (event.key === 'Enter') {
        sendMessage();
    }
}

// Display message
function displayMessage(message) {
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';
    
    if (currentUser && message.sender_id === currentUser.user_id) {
        messageDiv.classList.add('sent');
    }

    // Extract time directly from database timestamp string (format: "2025-10-15 16:16:25")
    // This avoids timezone conversion issues
    let time = '';
    if (message.created_at) {
        const timeStr = message.created_at.toString();
        const timeMatch = timeStr.match(/(\d{2}):(\d{2}):/);
        if (timeMatch) {
            time = `${timeMatch[1]}:${timeMatch[2]}`;
        }
    }

    messageDiv.innerHTML = `
        <div class="message-header">
            <span class="message-sender">${escapeHtml(message.sender_name)}</span>
            <span class="message-time">${time}</span>
        </div>
        <div class="message-content">${escapeHtml(message.message)}</div>
    `;

    messagesDiv.appendChild(messageDiv);
}

// Show add friend modal
function showAddFriend() {
    document.getElementById('add-friend-modal').style.display = 'block';
    document.getElementById('search-username').value = '';
    document.getElementById('search-results').innerHTML = '';
    document.getElementById('add-friend-message').textContent = '';
}

// Close add friend modal
function closeAddFriend() {
    document.getElementById('add-friend-modal').style.display = 'none';
}

// Search users
let searchTimeout;
async function searchUsers() {
    clearTimeout(searchTimeout);
    const query = document.getElementById('search-username').value.trim();
    const resultsDiv = document.getElementById('search-results');

    if (!query) {
        resultsDiv.innerHTML = '';
        return;
    }

    searchTimeout = setTimeout(async () => {
        try {
            const response = await fetch(`/api/friends/search?q=${encodeURIComponent(query)}&user_id=${currentUser.user_id}`);
            const data = await response.json();

            if (data.success && data.data) {
                resultsDiv.innerHTML = '';
                
                if (data.data.length === 0) {
                    resultsDiv.innerHTML = '<p style="text-align:center; color: #999;">No users found</p>';
                } else {
                    data.data.forEach(user => {
                        const resultItem = document.createElement('div');
                        resultItem.className = 'search-result-item';
                        resultItem.innerHTML = `
                            <span>${escapeHtml(user.username)}</span>
                            <button onclick="sendFriendRequest('${escapeHtml(user.username)}', this)">Add Friend</button>
                        `;
                        resultsDiv.appendChild(resultItem);
                    });
                }
            }
        } catch (error) {
            console.error('Error searching users:', error);
        }
    }, 300);
}

// Send friend request
async function sendFriendRequest(username, button) {
    const messageDiv = document.getElementById('add-friend-message');
    
    try {
        const response = await fetch('/api/friends/request', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                user_id: currentUser.user_id,
                username: username,
            }),
        });

        const data = await response.json();

        if (data.success) {
            messageDiv.className = 'modal-message success';
            messageDiv.innerHTML = `
                ✅ Friend request sent to <strong>${username}</strong>!<br>
                <small>They need to accept your request before you can chat.</small>
            `;
            button.disabled = true;
            button.textContent = 'Requested';
        } else {
            messageDiv.className = 'modal-message error';
            messageDiv.textContent = data.message;
        }
    } catch (error) {
        messageDiv.className = 'modal-message error';
        messageDiv.textContent = 'Network error. Please try again.';
        console.error('Error sending friend request:', error);
    }
}

// Load friend requests
async function loadFriendRequests() {
    try {
        const response = await fetch(`/api/friends/requests?user_id=${currentUser.user_id}`);
        const data = await response.json();

        if (data.success && data.data) {
            const count = data.data.length;
            if (count > 0) {
                document.getElementById('requests-badge').style.display = 'block';
                document.getElementById('request-count').textContent = count;
            } else {
                document.getElementById('requests-badge').style.display = 'none';
            }
        }
    } catch (error) {
        console.error('Error loading friend requests:', error);
    }
}

// Show friend requests modal
async function showFriendRequests() {
    document.getElementById('requests-modal').style.display = 'block';
    
    try {
        const response = await fetch(`/api/friends/requests?user_id=${currentUser.user_id}`);
        const data = await response.json();

        const requestsList = document.getElementById('requests-list');
        requestsList.innerHTML = '';

        if (data.success && data.data && data.data.length > 0) {
            data.data.forEach(request => {
                const requestItem = document.createElement('div');
                requestItem.className = 'request-item';
                requestItem.id = `request-${request.id}`;
                requestItem.innerHTML = `
                    <div class="request-info">
                        <span>${escapeHtml(request.username)}</span>
                        <small>${new Date(request.created_at).toLocaleDateString()}</small>
                    </div>
                    <div class="request-actions">
                        <button class="btn-accept" onclick="acceptRequest(${request.id})">Accept</button>
                        <button class="btn-reject" onclick="rejectRequest(${request.id})">Reject</button>
                    </div>
                `;
                requestsList.appendChild(requestItem);
            });
        } else {
            requestsList.innerHTML = '<p style="text-align:center; color: #999; padding: 20px;">No pending requests</p>';
        }
    } catch (error) {
        console.error('Error loading friend requests:', error);
    }
}

// Close requests modal
function closeRequests() {
    document.getElementById('requests-modal').style.display = 'none';
}

// Accept friend request
async function acceptRequest(requestId) {
    try {
        const response = await fetch(`/api/friends/accept/${requestId}`, {
            method: 'POST',
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById(`request-${requestId}`).remove();
            loadFriendRequests();
            loadFriends();
            
            const requestsList = document.getElementById('requests-list');
            if (requestsList.children.length === 0) {
                requestsList.innerHTML = '<p style="text-align:center; color: #999; padding: 20px;">No pending requests</p>';
            }
        } else {
            alert('Error accepting request: ' + data.message);
        }
    } catch (error) {
        console.error('Error accepting request:', error);
    }
}

// Reject friend request
async function rejectRequest(requestId) {
    try {
        const response = await fetch(`/api/friends/reject/${requestId}`, {
            method: 'POST',
        });

        const data = await response.json();

        if (data.success) {
            document.getElementById(`request-${requestId}`).remove();
            loadFriendRequests();
            
            const requestsList = document.getElementById('requests-list');
            if (requestsList.children.length === 0) {
                requestsList.innerHTML = '<p style="text-align:center; color: #999; padding: 20px;">No pending requests</p>';
            }
        } else {
            alert('Error rejecting request: ' + data.message);
        }
    } catch (error) {
        console.error('Error rejecting request:', error);
    }
}

// Connect WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/${currentUser.user_id}`;

    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        console.log('WebSocket connected');
    };

    ws.onmessage = function(event) {
        const wsMessage = JSON.parse(event.data);
        
        if (wsMessage.type === 'message') {
            const message = wsMessage.data;
            
            // If message is from current chat friend, display it
            if (currentFriend && message.sender_id === currentFriend.id) {
                displayMessage(message);
                scrollToBottom();
            }
            
            // Refresh friends list to update unread count
            loadFriends();
        } else if (wsMessage.type === 'friend_request') {
            loadFriendRequests();
        } else if (wsMessage.type === 'friend_accepted') {
            loadFriends();
        }
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };

    ws.onclose = function() {
        console.log('WebSocket disconnected');
        // Attempt to reconnect after 3 seconds
        if (currentUser) {
            setTimeout(connectWebSocket, 3000);
        }
    };
}

// Scroll to bottom
function scrollToBottom() {
    const messagesDiv = document.getElementById('messages');
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, m => map[m]);
}

// Close modals when clicking outside
window.onclick = function(event) {
    const addFriendModal = document.getElementById('add-friend-modal');
    const requestsModal = document.getElementById('requests-modal');
    
    if (event.target === addFriendModal) {
        addFriendModal.style.display = 'none';
    }
    if (event.target === requestsModal) {
        requestsModal.style.display = 'none';
    }
};
