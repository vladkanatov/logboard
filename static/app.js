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
fetchAllLogs();

wsocks[currentTab].onopen = (event) => {
  console.log(`Websocket ${currentTab} opened`);
};
wsocks[currentTab].onclose = (event) => {
  console.log(`Websocket ${currentTab} closed`);
};

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(tab) {
  wsocks[currentTab].close();
  logDisplay.replaceChildren();
  searchText[currentTab] = searchInput.value;
  currentTab = tab;
  if (wsocks[currentTab] === null)
    wsocks[currentTab] = new WebSocket(
      `ws://localhost:8000/logs?tab=${currentTab}`,
    );
  fetchAllLogs();
  wsocks[currentTab].onmessage = (event) => lineHandler(event.data);
  wsocks[currentTab].addEventListener('open', (event) => {
    console.log(`Websocket ${currentTab} opened`);
  });
  wsocks[currentTab].addEventListener('close', (event) => {
    console.log(`Websocket ${currentTab} closed`);
  });
  searchInput.value = searchText[currentTab];
}

async function fetchAllLogs() {
  const response = await fetch(
    `http://localhost:8000/all_logs?tab=${currentTab}`,
  );
  const data = await response.text();
  data.split('\n').forEach((el) => lineHandler(el));
}

const chgSortDirHandler = () => {
  logDisplay.classList.contains('reversed')
    ? logDisplay.classList.replace('reversed', 'stright')
    : logDisplay.classList.replace('stright', 'reversed');
  logDisplay.scrollTop = logDisplay.scrollHeight;
};

const lineHandler = (line) => {
  const logLine = document.createElement('div');
  if (line.startsWith('success:')) {
    logLine.classList.add('success', 'log');
    line = line.slice(8).trim(); // Убираем "success:" из вывода и пробелы
  } else if (line.startsWith('error:')) {
    logLine.classList.add('error', 'log');
    line = line.slice(6).trim(); // Убираем "error:" из вывода и пробелы
  } else if (line.startsWith('info:')) {
    logLine.classList.add('info', 'log');
    line = line.slice(5).trim(); // Убираем "info:" из вывода и пробелы
  } else {
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
