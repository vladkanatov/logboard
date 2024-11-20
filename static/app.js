let currentTab = 'packages-common';
let searchText = { 'packages-common': '', eap: '', sdk: '' };
const logDisplay = document.getElementById('log-display');
const searchInput = document.getElementById('search-input');
searchInput.value = 'sasd';

const ws = new WebSocket('ws://localhost:8000/logs?tab=packages-common');

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(tab) {
  currentTab = tab;
  searchInput.value = searchText[currentTab];
}
searchInput.addEventListener('input', () => {
  searchText[currentTab] = searchInput.value;
  console.log(searchText, currentTab);
});

const chgSortDirHandler = () => {
  logDisplay.classList.contains('reversed')
    ? logDisplay.classList.replace('reversed', 'stright')
    : logDisplay.classList.replace('stright', 'reversed');
};

ws.addEventListener('open', (event) => {
  console.log('Websocket connection opened');
});
ws.addEventListener('close', (event) => {
  console.log('Websocket connection closed');
});
ws.onmessage = function (line) {
  const logLine = document.createElement('div');
  if (line.data.startsWith('success:')) {
    logLine.classList.add('success', 'log');
    line = line.slice(8).trim(); // Убираем "success:" из вывода и пробелы
  } else if (line.data.startsWith('error:')) {
    logLine.classList.add('error', 'log');
    line = line.data.slice(6).trim(); // Убираем "error:" из вывода и пробелы
  } else if (line.data.startsWith('info:')) {
    logLine.classList.add('info', 'log');
    line = line.data.slice(5).trim(); // Убираем "info:" из вывода и пробелы
  } else {
    line = line.data;
    logLine.classList.add('date-divider');
  }

  // Обработка Markdown-ссылок в формате [текст](ссылка)
  const formattedLine = line.replace(
    /\[([^\]]+)\]\((https?:\/\/[^\s]+)\)/g,
    '<a href="$2" target="_blank">$1</a>',
  );
  logLine.innerHTML = formattedLine; // Вставляем обработанную строку как HTML

  logDisplay.appendChild(logLine);
};
