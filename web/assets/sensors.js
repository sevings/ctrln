function connectWebsocket() {
    const socket = new WebSocket("ws://" + document.location.host + "/ws");
    socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);
        if(data.sensor === "enable") {
            switchGroupElement(data.group, data.status);
            return;
        }

        const card = document.getElementById(data.sensor);
        if(!card)
            return;

        card.classList.toggle("bg-danger", !!data.critical);

        const group = card.closest(".col");
        const isSafe = !group.querySelector(".card.bg-danger");
        const button = group.getElementsByTagName("button")[0];
        button.disabled = isSafe;

        const body = card.getElementsByClassName("card-body")[0];
        body.innerHTML = data.status;
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
        form.action = "/" + group + "/" + sensor + "/name";
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

function initDisableButtons() {
    const switchGroup = function (event) {
        const button = event.currentTarget;
        const card = button.closest(".card");
        const status = card.dataset.status === "on" ? "off" : "on";
        const group = card.closest(".col").id;

        const url = "/" + group + "/" + status;
        fetch(url, {
            method: "post"
        });

        switchGroupElement(group, status);
    };

    const groupButtons = document.getElementsByClassName("switch-group");
    for(let i = 0; i < groupButtons.length; i++) {
        groupButtons[i].addEventListener("click", switchGroup);
    }
}

document.addEventListener("DOMContentLoaded", initDisableButtons);

function switchGroupElement(group, status) {
    const col = document.getElementById(group);
    const button = col.getElementsByTagName("button")[0];
    const card = button.closest(".card");
    card.dataset.status = status;

    const body = card.getElementsByClassName("card-body")[0];
    if(status === "on") {
        body.innerHTML = "enabled";
        button.innerHTML = "Turn off";
    } else if(status === "off") {
        body.innerHTML = "disabled";
        button.innerHTML = "Turn on";
    }
}
