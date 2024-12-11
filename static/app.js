let currentTab = 'packages-common';
let searchText = { 'packages-common': '', eap: '', sdk: '' };
let tabs = ['packages-common', 'eap', 'sdk'];
let isAutoScrollEnabled = true;
const logContainer = document.querySelector('.cont');
const logDisplay = document.getElementById('log-display');
const searchInput = document.getElementById('search-input');
const showTimeInput = document.getElementById('show-date');
const wsocks = {};

tabs.forEach((el) => {
  wsocks[el] = new WebSocket(`/logs?tab=${el}`);
  wsocks[el].onmessage = (line) => lineHandler(line.data);
  wsocks[el].addEventListener('open', (event) => {
    console.log(`Websocket ${el} opened`);
  });
  wsocks[el].addEventListener('close', (event) => {
    console.log(`Websocket ${el} closed`);
  });
});

logContainer.addEventListener('scroll', () => {
  const isAtBottom =
    logContainer.scrollHeight - logContainer.scrollTop === logContainer.clientHeight;

  isAutoScrollEnabled = isAtBottom;
});

fetchAllLogs();

//фильтрация
const searching = () => {
  let logs = document.getElementsByClassName('log');
  for (let log of logs) {
    const textToSearch = log.textContent;
    textToSearch.toLowerCase().includes(searchInput.value.toLowerCase())
      ? (log.style.display = 'block')
      : (log.style.display = 'none');
  }
};

// Устанавливает текущую вкладку и загружает данные для нее
function setCurrentTab(obj) {
  document
    .getElementsByClassName('active-tab')[0]
    .classList.remove('active-tab');
  obj.classList.add('active-tab');
  let tab = obj.textContent.toLowerCase();
  searchText[currentTab] = searchInput.value;
  currentTab = tab;
  fetchAllLogs();
  searchInput.value = searchText[currentTab];
}

async function fetchAllLogs() {
  logDisplay.innerHTML = '';
  const response = await fetch(
    `/all_logs?tab=${currentTab}`,
    { cache: 'no-cache' },
  );
  const data = await response.text();
  data.split('\n').forEach((el) => lineHandler(el));
  searching();
  
  // Прокрутка вниз после загрузки всех логов
  logContainer.scrollTop = logContainer.scrollHeight;
}

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
  } else if (line === '') {
    return;
  }

  const markdownLinkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
  line = line.replace(markdownLinkRegex, '<a href="$2" target="_blank">$1</a>');


  logLine.innerHTML = line
  logDisplay.appendChild(logLine);

  // Если автоскролл включен, прокручиваем вниз
  if (isAutoScrollEnabled) {
    logContainer.scrollTop = logContainer.scrollHeight;
  }
};


wsocks[currentTab].onmessage = (event) => lineHandler(event.data);

searchInput.addEventListener('input', searching);
