document.addEventListener('DOMContentLoaded',() => {
    const btn1 = document.getElementById("logout")
    const btn2 = document.getElementById("execute")

    btn1.onclick = () => {
        logoutSesion()
    }

    btn2.onclick = () => {
        executeComand()
    }

    async function executeComand(){
        let area1 = document.getElementById('floatingTextarea').value
        if (area1 === "") {
          alert("por favor ingrese algun comando")
        } else  {
          ExecuteCommand(area1)
        }
    }

    async function logoutSesion(){
        console.log("helo")
        const envio = await fetch('http://127.0.0.1:8000/logout',{
        method : "GET"
        })
        const response = await envio.json()
        let arrRes = response.Message
        alert(arrRes)
        window.location.replace("login.html")
    }

    async function ExecuteCommand(content) {
        const envio = await fetch('http://127.0.0.1:8000/individualComand',{
         method : "POST",
         body : JSON.stringify({"cmd":content}),
        })
        document.getElementById("graph").textContent = ""
        const response = await envio.json()
        let arrRes = response.Message
        let typeC = response.typeC
        if (typeC === "rep") {
            console.log(arrRes)
            d3.select("#graph")
            .graphviz()
                .dot(arrRes)
                .render();
        } else {
            alert(Message)
        }
    }
    
})

