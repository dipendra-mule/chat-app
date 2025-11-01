class ChatApp {
  constructor() {
    this.ws = null;
    this.username = '';
    this.room = '';
    this.connected = false;
    this.typingTimer = null;
    this.users = new Set();

    this.initializeEventListeners();
  }

  initializeEventListeners() {
    const messageInput = document.getElementById('messageInput');

    messageInput.addEventListener('keypress', e => {
      if (e.key === 'Enter') {
        this.sendMessage();
      } else {
        this.sendTyping();
      }
    });

    // Handle page visibility change
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.sendActivity(false);
      } else {
        this.sendActivity(true);
      }
    });
  }

  connect() {
    this.username = document.getElementById('usernameInput').value.trim();
    this.room = document.getElementById('roomInput').value.trim() || 'general';

    if (!this.username) {
      alert('Please enter a username');
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?username=${encodeURIComponent(this.username)}&room=${encodeURIComponent(
      this.room
    )}`;

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      this.connected = true;
      this.updateConnectionStatus(true);
      this.showChat();
      this.addSystemMessage('Connected to chat room');
    };

    this.ws.onmessage = event => {
      const messages = event.data.split('\n');
      messages.forEach(message => {
        if (message.trim()) {
          this.handleMessage(JSON.parse(message));
        }
      });
    };

    this.ws.onclose = () => {
      this.connected = false;
      this.updateConnectionStatus(false);
      this.addSystemMessage('Disconnected from chat');
    };

    this.ws.onerror = error => {
      console.error('WebSocket error:', error);
      this.addSystemMessage('Connection error occurred');
    };
  }

  handleMessage(data) {
    switch (data.type) {
      case 'chat':
        this.addChatMessage(data);
        break;
      case 'join':
        this.addSystemMessage(`${data.username} joined the room`);
        this.users.add(data.username);
        this.updateUsersList();
        break;
      case 'leave':
        this.addSystemMessage(`${data.username} left the room`);
        this.users.delete(data.username);
        this.updateUsersList();
        break;
      case 'typing':
        this.showTypingIndicator(data.username);
        break;
      case 'error':
        this.addSystemMessage(`Error: ${data.content}`, true);
        break;
    }
  }

  addChatMessage(data) {
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');

    const isOwnMessage = data.username === this.username;
    messageDiv.className = `message ${isOwnMessage ? 'own' : 'other'}`;

    const time = new Date(data.timestamp).toLocaleTimeString();
    messageDiv.innerHTML = `
            <div class="message-header">
                ${!isOwnMessage ? data.username : 'You'} • ${time}
            </div>
            <div class="message-content">${this.escapeHtml(data.content)}</div>
        `;

    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }

  addSystemMessage(content, isError = false) {
    const messagesDiv = document.getElementById('messages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message system ${isError ? 'error' : ''}`;
    messageDiv.textContent = content;

    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }

  sendMessage() {
    const input = document.getElementById('messageInput');
    const content = input.value.trim();

    if (!content || !this.connected) return;

    const message = {
      type: 'chat',
      content: content,
      timestamp: new Date().toISOString()
    };

    this.ws.send(JSON.stringify(message));
    input.value = '';
    this.hideTypingIndicator();
  }

  sendTyping() {
    if (!this.connected) return;

    clearTimeout(this.typingTimer);

    const message = {
      type: 'typing',
      timestamp: new Date().toISOString()
    };

    this.ws.send(JSON.stringify(message));

    this.typingTimer = setTimeout(() => {
      this.hideTypingIndicator();
    }, 3000);
  }

  showTypingIndicator(username) {
    const indicator = document.getElementById('typingIndicator');
    indicator.textContent = `${username} is typing...`;

    clearTimeout(this.typingTimer);
    this.typingTimer = setTimeout(() => {
      this.hideTypingIndicator();
    }, 3000);
  }

  hideTypingIndicator() {
    document.getElementById('typingIndicator').textContent = '';
  }

  updateConnectionStatus(connected) {
    const status = document.getElementById('status');
    status.textContent = connected ? 'Connected' : 'Disconnected';
    status.className = connected ? 'connected' : 'disconnected';
  }

  showChat() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('chatContainer').style.display = 'flex';
    document.getElementById('currentRoom').textContent = this.room;

    document.getElementById('messageInput').focus();
  }

  updateUsersList() {
    const usersList = document.getElementById('usersList');
    usersList.innerHTML = '';

    this.users.forEach(user => {
      const li = document.createElement('li');
      li.textContent = user;
      usersList.appendChild(li);
    });
  }

  escapeHtml(unsafe) {
    return unsafe.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#039;');
  }

  sendActivity(active) {
    if (this.connected) {
      const message = {
        type: active ? 'active' : 'inactive',
        timestamp: new Date().toISOString()
      };
      this.ws.send(JSON.stringify(message));
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Global functions for HTML onclick handlers
const chatApp = new ChatApp();

function connect() {
  chatApp.connect();
}

function sendMessage() {
  chatApp.sendMessage();
}

// Load room stats
async function loadRoomStats() {
  try {
    const response = await fetch('/stats');
    const stats = await response.json();
    document.getElementById('roomStats').innerHTML = Object.entries(stats)
      .map(([room, count]) => `<div>${room}: ${count} users</div>`)
      .join('');
  } catch (error) {
    console.error('Failed to load room stats:', error);
  }
}

// Load stats every 30 seconds
setInterval(loadRoomStats, 30000);
loadRoomStats();
