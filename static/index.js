
(document.onload = function () {

    function fillRoutesFile() {
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", "/route", false);
        xmlHttp.setRequestHeader("Accept", "application/yaml")
        xmlHttp.send(null);
        document.querySelector("#routes-config-file").innerHTML = xmlHttp.responseText;
    }

    function fillRoutesList() {
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", "/route", false);
        xmlHttp.send(null);
        JSON.parse(xmlHttp.response).routes.forEach(element => {
            var li = document.createElement("li");
            li.classList = ["list-group-item"];
            var text = document.createTextNode(`[${element.method}] ${element.path}`);
            li.appendChild(text);

            var element = document.querySelector("#my-routes-list");
            element.appendChild(li);
        });
    }

    function fillRequestList() {
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", "/request", false);
        xmlHttp.send(null);
        JSON.parse(xmlHttp.responseText).forEach(el => {
            let tr = document.createElement("tr");
            let colID = document.createElement("th"); // ID
            colID.scope = ["row"];
            colID.innerHTML = el.id;
            tr.appendChild(colID);

            let colReqId = tr.appendChild(document.createElement("td")); // ReqID
            colReqId.innerHTML = el.reqId;
            tr.appendChild(colReqId);

            let colMothod = document.createElement("th"); // Method
            colMothod.innerHTML = el.request.method;
            tr.appendChild(colMothod);

            let colPath = document.createElement("th"); // Path
            colPath.innerHTML = el.request.path;
            tr.appendChild(colPath);

            let colRequest = document.createElement("th"); // Request
            colRequest.innerHTML = el.request.body;
            tr.appendChild(colRequest);

            let colResponse = document.createElement("th"); // Response
            colResponse.innerHTML = el.response.body;
            tr.appendChild(colResponse);

            let colStatusCode = document.createElement("th"); // Code
            colStatusCode.innerHTML = el.response.code;
            tr.appendChild(colStatusCode);

            document.querySelector("#requests-tbody").appendChild(tr);
        });
    }

    function ajaxRequest() {
        var url = '/request'
        $.get(url).then(function (res) {
            params.success(res);
        })
    }

    fillRoutesFile();
    fillRoutesList();
    fillRequestList();
})();
