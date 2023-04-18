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

function initRenameModal() {
    const renameModal = document.getElementById("renameModal");
    const form = renameModal.getElementsByTagName("form")[0];
    const input = form.getElementsByTagName("input")[0];

    renameModal.addEventListener("show.bs.modal", (event) => {
        const link = event.relatedTarget;
        const card = link.closest(".card");
        if(!card)
            return;

        const col = card.closest(".col");
        if(!col)
            return;

        const name = link.innerHTML;
        const sensor = card.id;
        const group = col.id;

        form.dataset.sensor = sensor;
        form.action = "/" + group + "/" + sensor;
        input.value = name;
    });

    const renameSensor = function () {
        const data = new FormData(form);
        fetch(form.action, {
            method: "post",
            body: data,
        });

        const card = document.getElementById(form.dataset.sensor);
        const footer = card.getElementsByClassName("card-footer")[0];
        const link = footer.getElementsByTagName("a")[0];
        link.innerHTML = input.value;
    };

    const saveButton = document.getElementById("saveName");
    saveButton.addEventListener("click", renameSensor);

    input.addEventListener("keypress", (event) => {
        if(event.key === "Enter") {
            event.preventDefault();
            renameSensor();
        }
    });
}

document.addEventListener("DOMContentLoaded", initRenameModal);

function loadSensorNames() {
    fetch("/sensors").then((resp) => {
        return resp.json();
    }).then((names) => {
        for(const group in names) {
            for(const sensor in names[group]) {
                const card = document.getElementById(sensor);
                if(!card)
                    continue;

                const footer = card.getElementsByClassName("card-footer")[0];
                const link = footer.getElementsByTagName("a")[0];
                link.innerHTML = names[group][sensor];
            }
        }
    });
}

document.addEventListener("DOMContentLoaded", loadSensorNames);
