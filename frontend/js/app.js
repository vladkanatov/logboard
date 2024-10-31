const ws = new WebSocket("ws://192.168.100.179:8080/ws");

ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  updateTab(message);
};

function updateTab(data) {
  const tab = document.getElementById(data.tab);
  const statusClass = data.status === 'success' ? 'success' : data.status === 'error' ? 'error' : 'info';
  tab.innerHTML += `<p class="${statusClass}">${data.data}</p>`;
}

function showTab(tabName) {
  document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
  document.getElementById(tabName).classList.add('active');
}
