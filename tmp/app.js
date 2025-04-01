// Базовый URL API
const API_URL = 'http://localhost:8080';

// Функции для работы с токеном
const setToken = (token) => localStorage.setItem('token', token);
const getToken = () => localStorage.getItem('token');
const removeToken = () => localStorage.removeItem('token');

// Базовая функция для API запросов
async function apiRequest(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json',
    };
    const token = getToken();
    // if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    // }

    const response = await fetch(`${API_URL}${endpoint}`, {
        method,
        headers,
        body: body ? JSON.stringify(body) : null
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
}

// Функции аутентификации
async function register(username, email, password) {
    const response = await apiRequest('/auth/register', 'POST', {
        username,
        email,
        password
    });
    setToken(response.token);
    return response;
}

async function login(email, password) {
    const response = await apiRequest('/auth/login', 'POST', {
        email,
        password
    });
    setToken(response.token);
    return response;
}

// Функции для работы с чатами
async function createChat(name) {
    return await apiRequest('/chats/create', 'POST', { name });
}

async function getChats() {
    return await apiRequest('/chats/get', 'POST');
}

// Функции для работы с сообщениями
async function sendMessage(chatId, content) {
    return await apiRequest('/messages/create', 'POST', {
        chat_id: chatId,
        content
    });
}

async function getMessages(chatId, limit = 50, offset = 0) {
    return await apiRequest('/messages/get', 'POST', { chat_id: chatId, limit: limit, offset: offset });
}

// Обработчики форм
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        await register(
            document.getElementById('regUsername').value,
            document.getElementById('regEmail').value,
            document.getElementById('regPassword').value
        );
        showChat();
        loadChats();
        updateUserInfo();
    } catch (error) {
        alert('Ошибка при регистрации: ' + error.message);
    }
});

document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        await login(
            document.getElementById('loginEmail').value,
            document.getElementById('loginPassword').value
        );
        showChat();
        loadChats();
        updateUserInfo();
    } catch (error) {
        alert('Ошибка при входе: ' + error.message);
    }
});

document.getElementById('messageForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const messageInput = document.getElementById('messageInput');
    const content = messageInput.value;
    const activeChatId = document.querySelector('.chat-item.active')?.dataset.chatId;

    if (!activeChatId) {
        alert('Выберите чат');
        return;
    }

    try {
        await sendMessage(activeChatId, content);
        messageInput.value = '';
        await loadMessages(activeChatId);
    } catch (error) {
        alert('Ошибка при отправке сообщения: ' + error.message);
    }
});

document.getElementById('newChatBtn').addEventListener('click', async () => {
    const chatName = prompt('Введите название чата:');
    if (chatName) {
        try {
            await createChat(chatName);
            await loadChats();
        } catch (error) {
            alert('Ошибка при создании чата: ' + error.message);
        }
    }
});

// Функции UI
function updateUserInfo() {
    const userId = getUserId();
    const userInfoElement = document.getElementById('userInfo');
    const userIdElement = document.getElementById('userId');
    
    if (userId) {
        userIdElement.textContent = userId;
        userInfoElement.style.display = 'block';
    } else {
        userInfoElement.style.display = 'none';
    }
}

function showChat() {
    document.getElementById('authForms').style.display = 'none';
    document.getElementById('chatContainer').style.display = 'block';
    updateUserInfo();
}

function showAuth() {
    document.getElementById('authForms').style.display = 'flex';
    document.getElementById('chatContainer').style.display = 'none';
    document.getElementById('userInfo').style.display = 'none';
}

async function loadChats() {
    try {
        const response = await getChats();
        const chatList = document.getElementById('chatList');
        const chatsHtml = response.chats.map(chat => `
            <div class="chat-item" data-chat-id="${chat.id}" onclick="selectChat('${chat.id}')">
                ${chat.title}
            </div>
        `).join('');
        chatList.innerHTML = `
            <div style="padding: 15px;">
                <button id="newChatBtn">Новый чат</button>
            </div>
            ${chatsHtml}
        `;
    } catch (error) {
        console.error('Ошибка при загрузке чатов:', error);
    }
}

// Функция для получения ID пользователя из токена
function getUserId() {
    const token = getToken();
    if (!token) return null;
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.user_id;
    } catch (error) {
        console.error('Ошибка при получении ID пользователя:', error);
        return null;
    }
}

async function loadMessages(chatId) {
    try {
        const response = await getMessages(chatId);
        const messagesContainer = document.getElementById('messagesContainer');
        
        if (!response.messages || !Array.isArray(response.messages)) {
            messagesContainer.innerHTML = '<div class="no-messages">Нет сообщений</div>';
            return;
        }

        const currentUserId = getUserId();
        const messagesHtml = response.messages.map(message => {
            if (!message || typeof message.content !== 'string') {
                return '';
            }
            const isOwnMessage = currentUserId && message.user_id === currentUserId;
            return `
                <div class="message ${isOwnMessage ? 'sent' : 'received'}">
                    ${message.content}
                </div>
            `;
        }).filter(html => html).join('');

        messagesContainer.innerHTML = messagesHtml || '<div class="no-messages">Нет сообщений</div>';
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    } catch (error) {
        console.error('Ошибка при загрузке сообщений:', error);
        document.getElementById('messagesContainer').innerHTML = 
            '<div class="error-message">Ошибка при загрузке сообщений</div>';
    }
}

function selectChat(chatId) {
    document.querySelectorAll('.chat-item').forEach(item => {
        item.classList.remove('active');
    });
    const selectedChat = document.querySelector(`[data-chat-id="${chatId}"]`);
    if (selectedChat) {
        selectedChat.classList.add('active');
        loadMessages(chatId);
    }
}

// Инициализация
if (getToken()) {
    showChat();
    loadChats();
} else {
    showAuth();
}