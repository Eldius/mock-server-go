
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
        console.log(xmlHttp.responseText);
        JSON.parse(xmlHttp.responseText).routes.forEach(element => {
            var li = document.createElement("li");
            li.classList = ["list-group-item"];
            var text = document.createTextNode(`[${element.method}] ${element.path}`);
            li.appendChild(text);

            var element = document.querySelector("#my-routes-list");
            element.appendChild(li);
        });
    }

    fillRoutesFile();
    fillRoutesList();
})();
