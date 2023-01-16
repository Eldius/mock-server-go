
(document.onload = function () {

    function fillRoutesFile() {
        var headers = new Headers();
        headers.set('Accept', 'application/yaml');
        var cfg = {
            method: 'GET',
            headers: headers,
        };

        fetch('/route', cfg).then(response => response.text())
            .then(response => {
                document.querySelector("#routes-config-file").innerHTML = response;

                let link = document.createElement("a");
                link.classList = ["btn btn-primary"]
                link.href = "data:application/yaml;base64," + btoa(response);
                link.download = "routes.yaml";
                link.innerText = "Download config file here";
        
                document.querySelector("#config-file-card").appendChild(link);
            }).catch(function (err) {
                // There was an error
                console.warn('Failed to fetch routes def file.', err);
            });
    }

    function fillRoutesList() {
        fetch('/route').then(response => response.json())
            .then(response => {
                response.routes.forEach(element => {
                    var li = document.createElement("li");
                    li.classList = ["list-group-item"];
                    var text = document.createTextNode(`[${element.method}] ${element.path}`);
                    li.appendChild(text);
        
                    var element = document.querySelector("#my-routes-list");
                    element.appendChild(li);
                });
            }).catch(function (err) {
                // There was an error
                console.warn('Failed to fetch routes def.', err);
            });
    }

    function fillRequestList() {
        fetch('/request').then(response => response.json())
            .then(response => {
                response.forEach(el => {
                    let tr = document.createElement("tr");
                    let colID = document.createElement("th"); // ID
                    colID.scope = ["row"];
                    colID.innerHTML = el.id;
                    tr.appendChild(colID);

                    let colReqId = tr.appendChild(document.createElement("td")); // ReqID
                    colReqId.innerHTML = el.reqId;
                    tr.appendChild(colReqId);

                    let colReqDate = tr.appendChild(document.createElement("td")); // ReqDate
                    colReqDate.innerHTML = el.requestDate;
                    tr.appendChild(colReqDate);

                    let colMothod = document.createElement("th"); // Method
                    colMothod.innerHTML = el.request.method;
                    tr.appendChild(colMothod);

                    let colPath = document.createElement("th"); // Path
                    colPath.innerHTML = el.request.path;
                    tr.appendChild(colPath);

                    let colRequest = document.createElement("th"); // Request
                    colRequest.innerHTML = el.request.body;
                    tr.appendChild(colRequest);

                    let colRequestHeaders = document.createElement("th"); // Request Headers
                    let divReqHeaders = document.createElement("div");
                    colRequestHeaders.appendChild(divReqHeaders);
                    Object.entries(el.request.headers != null ? el.request.headers : []).forEach(entry => {
                        const [key, values] = entry;
                        let p = document.createElement("p");
                        p.innerHTML = key + ": " + values;
                        divReqHeaders.appendChild(p);
                    });
                    tr.appendChild(colRequestHeaders);

                    let colResponse = document.createElement("th"); // Response
                    colResponse.innerHTML = el.response.body;
                    tr.appendChild(colResponse);

                    let colResponseHeaders = document.createElement("th"); // Response Headers
                    let divResHeaders = document.createElement("div");
                    colResponseHeaders.appendChild(divResHeaders);
                    Object.entries(el.response.headers != null ? el.response.headers : []).forEach(entry => {
                        const [key, values] = entry;
                        let p = document.createElement("p");
                        p.innerHTML = key + ": " + values;
                        divResHeaders.appendChild(p);
                    });
                    tr.appendChild(colResponseHeaders);

                    let colStatusCode = document.createElement("th"); // Code
                    colStatusCode.innerHTML = el.response.code;
                    tr.appendChild(colStatusCode);

                    document.querySelector("#requests-tbody").appendChild(tr);
                });
        }).catch(function (err) {
            // There was an error
            console.warn('Failed to fetch requests.', err);
        });
    }

    fillRoutesFile();
    fillRoutesList();
    fillRequestList();
})();
