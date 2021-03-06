
(document.onload = function () {

    const RoutesHandling = {
        data() {
            return {
                routes: []
            }
        },
        methods: {
            refresh() {
                fetch('/route')
                    .then((resp) => resp.json())
                    .then((data) => {
                        console.log(data);
                        console.log(JSON.stringify(data.routes));
                        this.routes = data.routes;
                    })
                    .catch(function (err) {
                        console.log(err);
                    });
            }
        }
    };
    const RoutesFileHandling = {
        data() {
            return {
                fileContent: "",
                link: {
                    href: "javascript:void(0);",
                    text: "Download config file here",
                    isDownloadDisabled: true
                }
            }
        },
        methods: {
            refreshFile() {
                fetch('/route', {
                    headers: {
                        Accept: 'application/yaml'
                    }
                })
                    .then((resp) => resp.text())
                    .then((data) => {
                        this.fileContent = data;
                        this.link.href = "data:application/yaml;base64," + btoa(data);
                        this.link.isDownloadDisabled = false;
                    })
                    .catch(function (err) {
                        console.log(err);
                    });
            }
        }
    };
    const RequestsHandling = {
        data() {
            return {
                requests: [],
            }
        },
        methods: {
            refreshTable() {
                fetch('/request')
                    .then((resp) => resp.json())
                    .then((data) => {
                        this.requests = data;
                    })
                    .catch(function (err) {
                        console.log(err);
                    });
            }
        }
    };

    Vue.createApp(RoutesHandling).mount('#routes');
    Vue.createApp(RoutesFileHandling).mount('#routes_file');
    Vue.createApp(RequestsHandling).mount('#requestsTable');

})();
