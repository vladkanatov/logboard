let currentTab = 'packages-common';
let searchText = { 'packages-common': '', eap: '', sdk: '' };
const logDisplay = document.getElementById('log-display');
const searchInput = document.getElementById('search-input');
searchInput.value = '';

const wsocks = {
  'packages-common': new WebSocket(
    'ws://localhost:8000/logs?tab=packages-common',
  ),
  eap: null,
  sdk: null,
};

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(tab) {
  wsocks[currentTab].close();
  wsocks[currentTab].removeEventListener('close', wsocks[currentTab]);
  logDisplay.replaceChildren();
  searchText[currentTab] = searchInput.value;
  currentTab = tab;
  const logDisplay = async fetch('http://localhost:8000');
  
  wsocks[currentTab].onmessage = lineHandler;
  wsocks[currentTab].addEventListener('open', (event) => {
    console.log(`Websocket ${currentTab} opened`);
  });
  wsocks[currentTab].addEventListener('close', (event) => {
    console.log(`Websocket ${currentTab} closed`);
  });
  searchInput.value = searchText[currentTab];
}

const chgSortDirHandler = () => {
  logDisplay.classList.contains('reversed')
    ? logDisplay.classList.replace('reversed', 'stright')
    : logDisplay.classList.replace('stright', 'reversed');
};

const lineHandler = (line) => {
  const logLine = document.createElement('div');
  if (line.data.startsWith('success:')) {
    logLine.classList.add('success', 'log');
    line = line.data.slice(8).trim(); // Убираем "success:" из вывода и пробелы
  } else if (line.data.startsWith('error:')) {
    logLine.classList.add('error', 'log');
    line = line.data.slice(6).trim(); // Убираем "error:" из вывода и пробелы
  } else if (line.data.startsWith('info:')) {
    logLine.classList.add('info', 'log');
    line = line.data.slice(5).trim(); // Убираем "info:" из вывода и пробелы
  } else {
    line = line.data;
    logLine.classList.add('date-divider', 'log');
  }

  // Обработка Markdown-ссылок в формате [текст](ссылка)
  const formattedLine = line.replace(
    /\[([^\]]+)\]\((https?:\/\/[^\s]+)\)/g,
    '<a href="$2" target="_blank">$1</a>',
  );
  logLine.innerHTML = formattedLine; // Вставляем обработанную строку как HTML

  logDisplay.appendChild(logLine);
};

wsocks[currentTab].onmessage = lineHandler;

const searching = (event) => {
  // searchText[currentTab] = searchInput.value;
  let logs = document.getElementsByClassName('log');
  if (searchInput.value.length < 2) {
    for (let log of logs) log.style.display = 'block';
  } else {
    for (let log of logs) {
      log.textContent.toLowerCase().includes(event.target.value.toLowerCase())
        ? (log.style.display = 'block')
        : (log.style.display = 'none');
    }
  }
};

searchInput.addEventListener('input', searching);
