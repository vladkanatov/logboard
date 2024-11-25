let searchText = { 'packages-common': '', eap: '', sdk: '' };
let currentTab = 'all';
let tabs = ['all'];
const logDisplay = document.getElementById('log-display');
const searchInput = document.getElementById('search-input');
const showTimeInput = document.getElementById('show-date');
const newTabNameInput = document.getElementById('new-tab-name');
const tabContainer = document.querySelector('.tabs');
const wsocks = {};

function setupWebSocket(tabName) {
  const ws = new WebSocket(`/logs?tab=${tabName}`);

  ws.onmessage = (line) => lineHandler(line.data);

  ws.addEventListener('open', () => {
    console.log(`WebSocket for ${tabName} opened`);
  });

  ws.addEventListener('close', () => {
    console.log(`WebSocket for ${tabName} closed`);
  });

  wsocks[tabName] = ws;
}

tabs.forEach(setupWebSocket);

initializeTabs();
fetchAllLogs();

function initializeTabs() {
  tabContainer.innerHTML = '';
  tabs.forEach((tab) => addTabButton(tab));
  const firstTabButton = document.querySelector('.tab-button');
  if (firstTabButton) setCurrentTab(firstTabButton); // Устанавливаем активной первую вкладку, если она есть
}

//фильтрация
const searching = () => {
  let logs = document.getElementsByClassName('log');
  for (let log of logs) {
    const textToSearch = showTimeInput.checked
      ? log.textContent.slice(20)
      : log.textContent;
    textToSearch.toLowerCase().includes(searchInput.value.toLowerCase())
      ? (log.style.display = 'block')
      : (log.style.display = 'none');
  }
};

function addNewTab() {
  newTabNameInput.style.display = 'inline';
  newTabNameInput.focus();

  newTabNameInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      const tabName = newTabNameInput.value.trim();
      if (tabName && !tabs.includes(tabName)) {
        tabs.push(tabName);
        addTabButton(tabName);
        setupWebSocket(tabName)
        setCurrentTab(
          document.querySelector(`.tab-button:last-child`) // Установим активной только что созданную вкладку
        );
      }
      newTabNameInput.value = '';
      newTabNameInput.style.display = 'none';
    }
  });
}

// Добавление кнопки для новой вкладки
function addTabButton(tabName) {
  const button = document.createElement('button');
  button.textContent = tabName;
  button.className = 'tab-button';
  button.onclick = () => setCurrentTab(button);

  // Двойной клик для изменения названия
  button.ondblclick = () => renameTab(button);

  tabContainer.appendChild(button);
}


// Установить текущую вкладку
function setCurrentTab(button) {
  if (!button) return; // Защита от вызова с undefined
  document
    .querySelectorAll('.tab-button')
    .forEach((btn) => btn.classList.remove('active-tab'));
  button.classList.add('active-tab');
  currentTab = button.textContent;
  fetchAllLogs();
}

// Переименовать вкладку
function renameTab(button) {
  const originalName = button.textContent;
  const input = document.createElement('input');
  input.type = 'text';
  input.value = originalName;
  input.className = 'tab-rename-input';

  // Сохранить размеры кнопки
  const buttonWidth = button.offsetWidth;
  input.style.width = `${buttonWidth}px`;

  button.textContent = ''; // Очистка кнопки для вставки ввода
  button.appendChild(input);
  input.focus();

  input.addEventListener('keydown', async (e) => {
    if (e.key === 'Enter') {
      const newName = input.value.trim();

      if (newName && !tabs.includes(newName)) {
        try {
          // Отправить запрос на сервер
          await renameTabOnBackend(originalName, newName);

          if (wsocks[originalName]) {
            // Закрываем старое WebSocket-соединение
            wsocks[originalName].close();
            delete wsocks[originalName];
          }
          
          // Создаем новое WebSocket-соединение для нового имени
          setupWebSocket(newName);
          
          // Обновить данные на клиенте
          tabs[tabs.indexOf(originalName)] = newName;
          button.textContent = newName;
          if (currentTab === originalName) {
            currentTab = newName
            fetchAllLogs()
          };
        } catch (error) {
          console.error('Ошибка переименования на сервере:', error);
          button.textContent = originalName; // Вернуть оригинальное имя при ошибке
        }
      } else {
        button.textContent = originalName; // Вернуть оригинальное имя при конфликте
      }
    }
  });

  input.addEventListener('blur', () => {
    button.textContent = originalName; // Возвращаем оригинальное имя, если ввод не завершён
  });
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
}

async function renameTabOnBackend(oldName, newName) {
  // const response = await fetch('/rename-tab', {
  //   method: 'POST',
  //   headers: {
  //     'Content-Type': 'application/json',
  //   },
  //   body: JSON.stringify({ oldName, newName }),
  // });

  // if (!response.ok) {
  //   throw new Error('Ошибка на сервере');
  // }
  console.log("Cool!")
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
  } else if (line === '') {
  } else {
  }

  logLine.innerHTML = showTimeInput.checked ? line : line.slice(20); // обработка даты
  // Вставляем обработанную строку как HTML
  logDisplay.appendChild(logLine);
};

wsocks[currentTab].onmessage = (event) => lineHandler(event.data);

searchInput.addEventListener('input', searching);
