$(document).ready(function () {
    console.log("ready");
    $('#myTabs').tabs();
    connectWebsocket();
    loadStatus();
});

function loadStatus() {
    console.log("loading...");
    $.get("/list", function (data, status) {
        body = $('#tbodyStatus')
        body.empty();

        console.log(status)
        console.log(data)

        json = data;

        json.map((i) => {
            body.prepend(createRow(i.currentBlock, new Date(i.currentDate * 1000), ''));
        });
    });
}


function createRow(rowName, value) {
    row =
        '<tr>' +
        '<td class="left row">' + rowName + '</td>' +
        '<td class="right">' + value + '</td>' +
        '</tr>'

    return row
}

function openIndex(id) {
    console.log(id)
    //window.location.replace("/create.html/" + id);

}


const websocketUrl = 'wss://localhost:8080/ws';

const connectWebsocket = () => {
    console.log("trying to connect to websocket...")
    const ws = new WebSocket(websocketUrl);
    ws_status = $('#ws');
    ws_update = $('#ws_update');
    ws.onopen = (evt) => {
        console.log("websocket opened");
        ws_status.text("ws open");
    };

    ws.onclose = (evt) => {
        ws_status.text("ws closed");
        console.log("weboscket closed");
        const id = setInterval(() => {
            clearInterval(id)
            connectWebsocket();
            console.log("interval...")
        }, 3000);
    };

    ws.onmessage = function (message) {
        ws_update.text('last update: ' + new (Date));
        const json = JSON.parse(message.data);
        console.log("websocket message: " + json.messageType);
        switch (json.messageType) {
            case "syncStatus":
                const i = json.payload;
                body.prepend(createRow(i.currentBlock, new Date(i.currentDate * 1000), ''));
                console.log("websocket message payload: " + json.payload);
        }

    };

    ws.onerror = function (evt) {
        ws_status.text("ws error");
        console.log("websocket error: " + evt);
    };
}
