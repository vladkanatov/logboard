let currentTab = 'packages-common';

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(tab) {
    currentTab = tab;
    fetchLogs();
}

// Функция для получения логов выбранной вкладки
async function fetchLogs() {
    const response = await fetch(`/logs?tab=${currentTab}`);
    const data = await response.text();
    displayLogs(data);
}

// Функция для отображения логов с цветовым выделением
function displayLogs(data) {
    const logDisplay = document.getElementById('logDisplay');
    logDisplay.innerHTML = ''; // Очищаем старые логи

    data.split('\n').forEach(line => {
        const logLine = document.createElement('div');

        // Проверяем статус по ключевым словам и добавляем соответствующий класс
        if (line.startsWith('success:')) {
            logLine.classList.add('success');
            logLine.textContent = line.slice(8); // Убираем "success:" из вывода
        } else if (line.startsWith('error:')) {
            logLine.classList.add('error');
            logLine.textContent = line.slice(6); // Убираем "error:" из вывода
        } else if (line.startsWith('info:')) {
            logLine.classList.add('info');
            logLine.textContent = line.slice(5); // Убираем "info:" из вывода
        } else if (line.includes('---')) {  // Новый день
            logLine.classList.add('date-divider');
            logLine.textContent = line;
        }

        const formattedLine = line.replace(/\[([^\]]+)\]\((https?:\/\/[^\s]+)\)/g, '<a href="$2" target="_blank">$1</a>');
        logLine.innerHTML = formattedLine; // Вставляем обработанную строку как HTML

        logDisplay.appendChild(logLine);
    });
}

// Обновление логов каждую секунду
setInterval(fetchLogs, 1000);

// Начальная загрузка логов для первой вкладки
fetchLogs();
