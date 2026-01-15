// Основные функции JavaScript

// Валидация в реальном времени
function validateField(field, type) {
    const value = field.value.trim();
    const errorElement = document.getElementById(`${field.name}-error`);
    
    if (!value) {
        if (errorElement) {
            errorElement.innerHTML = '<span class="text-red-400">This field is required</span>';
            errorElement.classList.remove('hidden');
        }
        return false;
    }
    
    // Валидация по типу
    let isValid = true;
    let message = '';
    
    switch(type) {
        case 'username':
            if (value.length < 3) {
                isValid = false;
                message = 'Username must be at least 3 characters';
            } else if (value.length > 16) {
                isValid = false;
                message = 'Username cannot exceed 16 characters';
            } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
                isValid = false;
                message = 'Only letters, numbers and underscores allowed';
            }
            break;
            
        case 'email':
            if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
                isValid = false;
                message = 'Please enter a valid email address';
            }
            break;
            
        case 'password':
            if (value.length < 4) {
                isValid = false;
                message = 'Password must be at least 4 characters';
            } else if (value.length > 16) {
                isValid = false;
                message = 'Password cannot exceed 16 characters';
            }
            break;
    }
    
    if (errorElement) {
        if (!isValid) {
            errorElement.innerHTML = `<span class="text-red-400">${message}</span>`;
            errorElement.classList.remove('hidden');
        } else {
            errorElement.innerHTML = '<span class="text-green-400">✓ Valid</span>';
            errorElement.classList.remove('hidden');
            setTimeout(() => {
                errorElement.classList.add('hidden');
            }, 2000);
        }
    }
    
    return isValid;
}

// Копирование realmlist
function copyRealmlist() {
    const realmlist = document.getElementById('realmlist-text').textContent;
    navigator.clipboard.writeText(realmlist).then(() => {
        showNotification('Realmlist copied to clipboard!', 'success');
    });
}

// Уведомления
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `fixed top-4 right-4 p-4 rounded-lg shadow-lg z-50 transform translate-x-full transition-transform duration-300 ${
        type === 'success' ? 'bg-green-900 text-green-100' :
        type === 'error' ? 'bg-red-900 text-red-100' :
        'bg-blue-900 text-blue-100'
    }`;
    notification.innerHTML = `
        <div class="flex items-center">
            <i class="fas fa-${type === 'success' ? 'check' : 'exclamation'}-circle mr-3"></i>
            <span>${message}</span>
        </div>
    `;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.classList.remove('translate-x-full');
    }, 10);
    
    setTimeout(() => {
        notification.classList.add('translate-x-full');
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Тема (светлая/темная)
function toggleTheme() {
    const html = document.documentElement;
    const currentTheme = html.classList.contains('dark') ? 'dark' : 'light';
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    html.classList.remove(currentTheme);
    html.classList.add(newTheme);
    localStorage.setItem('theme', newTheme);
    
    showNotification(`Switched to ${newTheme} theme`);
}

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    // Восстановление темы
    const savedTheme = localStorage.getItem('theme') || 'dark';
    document.documentElement.classList.add(savedTheme);
    
    // Автоматическая валидация полей
    const inputs = document.querySelectorAll('input[data-validate]');
    inputs.forEach(input => {
        input.addEventListener('blur', () => {
            validateField(input, input.dataset.validate);
        });
        
        input.addEventListener('input', () => {
            const errorElement = document.getElementById(`${input.name}-error`);
            if (errorElement) {
                errorElement.classList.add('hidden');
            }
        });
    });
    
    // Обработка формы регистрации
    const form = document.getElementById('registration-form');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            // Валидация всех полей
            let isValid = true;
            inputs.forEach(input => {
                if (!validateField(input, input.dataset.validate)) {
                    isValid = false;
                }
            });
            
            if (!isValid) {
                showNotification('Please fix the errors in the form', 'error');
                return;
            }
            
            // Проверка совпадения паролей
            const password = form.querySelector('[name="password"]').value;
            const confirm = form.querySelector('[name="confirm_password"]').value;
            
            if (password !== confirm) {
                showNotification('Passwords do not match!', 'error');
                return;
            }
            
            // Отправка формы через HTMX
            htmx.trigger(form, 'submit');
        });
    }
    
    // Периодическое обновление статистики
    setInterval(() => {
        htmx.trigger('#server-stats', 'update');
    }, 30000);
});
