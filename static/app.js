let currentTab = 'packages-common';
const sockets = {}; // Храним активные WebSocket подключения
const logs = {}; // Сохраняем логи для каждой вкладки

// Устанавливает текущую вкладку и активирует соответствующий WebSocket
function setCurrentTab(tab) {
    // Если вкладка уже активна, не очищаем логи, только переключаем
    if (currentTab === tab) {
        console.log(`Switching to the same tab: ${tab}`);
        displayLogs(tab);  // Показываем логи для этой вкладки
        return;
    }

    // Очищаем логи, если вкладка меняется
    clearLogs();

    currentTab = tab;

    // Если WebSocket для вкладки уже существует, ничего не делаем
    if (sockets[currentTab]) {
        console.log(`Switching to existing WebSocket for tab: ${currentTab}`);
        displayLogs(currentTab);  // Показываем логи для этой вкладки
        return;
    }

    // Если WebSocket ещё не создан для этой вкладки, создаём его
    createWebSocket(currentTab);
}

// Создание или использование WebSocket подключения
function createWebSocket(tab) {
    const socket = new WebSocket(`ws://${window.location.host}/logs?tab=${tab}`);
    sockets[tab] = socket; // Сохраняем WebSocket в объект

    socket.onopen = () => {
        console.log(`WebSocket connection established for tab: ${tab}`);
    };

    socket.onmessage = (event) => {
        console.log(`Message received on tab ${tab}: ${event.data}`);
        // Проверяем, активна ли вкладка, чтобы обновлять логи
        if (currentTab === tab) {
            // Сохраняем логи в памяти для текущей вкладки
            if (!logs[tab]) {
                logs[tab] = [];
            }
            logs[tab].push(event.data);
            displayLogs(tab);
        }
    };

    socket.onclose = () => {
        console.log(`WebSocket for tab ${tab} closed.`);
        delete sockets[tab]; // Удаляем WebSocket из кеша
    };

    socket.onerror = (error) => {
        console.error(`WebSocket error on tab ${tab}:`, error);
    };
}

// Функция для отображения логов с цветовым выделением и обработкой Markdown-ссылок
function displayLogs(tab) {
    const logDisplay = document.getElementById('logDisplay');
    logDisplay.innerHTML = ''; // Очищаем старые логи перед добавлением новых

    if (logs[tab]) {
        logs[tab].forEach(data => {
            data.split('\n').forEach(line => {
                const logLine = document.createElement('div');

                // Проверяем статус по ключевым словам и добавляем соответствующий класс
                if (line.startsWith('success:')) {
                    logLine.classList.add('success');
                    line = line.slice(8).trim(); // Убираем "success:" из вывода и пробелы
                } else if (line.startsWith('error:')) {
                    logLine.classList.add('error');
                    line = line.slice(6).trim(); // Убираем "error:" из вывода и пробелы
                } else if (line.startsWith('info:')) {
                    logLine.classList.add('info');
                    line = line.slice(5).trim(); // Убираем "info:" из вывода и пробелы
                } else if (line.includes('---')) {  // Новый день
                    logLine.classList.add('date-divider');
                }

                // Обработка Markdown-ссылок в формате [текст](ссылка)
                const formattedLine = line.replace(/\[([^\]]+)\]\((https?:\/\/[^\s]+)\)/g, '<a href="$2" target="_blank">$1</a>');
                logLine.innerHTML = formattedLine; // Вставляем обработанную строку как HTML

                logDisplay.appendChild(logLine);  // Добавляем новую строку в конец
            });
        });
    }

    // Прокручиваем до самого низа
    logDisplay.scrollTop = logDisplay.scrollHeight;
}

// Функция для очистки логов
function clearLogs() {
    const logDisplay = document.getElementById('logDisplay');
    logDisplay.innerHTML = ''; // Очищаем текущие логи
}

// Инициализация WebSocket для первой вкладки
createWebSocket(currentTab);
