function updateDateTime() {
    var dt = new Date();
    var datetimeSpan = document.getElementById('datetime');
    datetimeSpan.textContent = dt.toLocaleString();
}

setInterval(updateDateTime, 1000);
updateDateTime(); 
