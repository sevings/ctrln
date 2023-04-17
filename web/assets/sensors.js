function connectWebsocket() {
    const socket = new WebSocket("ws://" + document.location.host + "/ws");
    socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);
        const card = document.getElementById(data.sensor);
        if(!card)
            return;

        card.classList.toggle("bg-danger", !!data.critical);

        const body = card.getElementsByClassName("card-body");
        if(!body)
            return;

        body[0].innerHTML = data.status;
    });
}

document.addEventListener("DOMContentLoaded", connectWebsocket);
