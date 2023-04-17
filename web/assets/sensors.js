function connectWebsocket() {
    const host = document.location.host;
    const socket = new WebSocket("ws://" + host + "/ws");
    socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);
        const card = document.getElementById(data.sensor);
        if(!card)
            return;

        const body = card.getElementsByClassName("card-body");
        if(!body)
            return;

        body[0].innerHTML = data.status;
    });
}

document.addEventListener("DOMContentLoaded", connectWebsocket);
