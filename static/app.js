let currentTab = 'packages-common';

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(tab) {
  currentTab = tab;
  fetchLogs();
}

let isStright = true;

const chgSortDirHandler = () => {
  chgSortDir();
  fetchLogs();
};

const reverseSwitch = () => {
  let sortStright = true;
  return () => {
    sortStright = !sortStright;
    isStright = sortStright;
    return sortStright;
  };
};
const chgSortDir = reverseSwitch();

// Функция для получения логов выбранной вкладки
async function fetchLogs() {
  const response = await fetch(`/logs?tab=${currentTab}`);
  const data = await response.text();
  displayLogs(data);
}
async function fetchQuery() {
  const response = await fetch(`/logs`, 'POST');
  const data = await response.text();
  displayLogs(data);
}

// Функция для отображения логов с цветовым выделением и обработкой Markdown-ссылок
function displayLogs(data) {
  let dataSorted;
  const logDisplay = document.getElementById('logDisplay');
  logDisplay.innerHTML = ''; // Очищаем старые логи
  console.log(isStright);
  dataSorted = isStright ? data.split('\n') : data.split('\n').reverse();
  dataSorted.forEach((line) => {
    const logLine = document.createElement('div');
    // Проверяем статус по ключевым словам и добавляем соответствующий класс
    if (line.startsWith('success:')) {
      logLine.classList.add('success', 'log');
      line = line.slice(8).trim(); // Убираем "success:" из вывода и пробелы
    } else if (line.startsWith('error:')) {
      logLine.classList.add('error', 'log');
      line = line.slice(6).trim(); // Убираем "error:" из вывода и пробелы
    } else if (line.startsWith('info:')) {
      logLine.classList.add('info', 'log');
      line = line.slice(5).trim(); // Убираем "info:" из вывода и пробелы
    } else if (line.includes('---')) {
      // Новый день
      logLine.classList.add('date-divider');
    }

    // Обработка Markdown-ссылок в формате [текст](ссылка)
    const formattedLine = line.replace(
      /\[([^\]]+)\]\((https?:\/\/[^\s]+)\)/g,
      '<a href="$2" target="_blank">$1</a>',
    );
    logLine.innerHTML = formattedLine; // Вставляем обработанную строку как HTML

    logDisplay.appendChild(logLine);
  });
}

// Обновление логов каждую секунду
setInterval(fetchLogs, 1000);

// Начальная загрузка логов для первой вкладки
fetchLogs();
