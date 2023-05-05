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
        const envio = await fetch('http://3.14.134.83/logout',{
        method : "GET"
        })
        const response = await envio.json()
        let arrRes = response.Message
        alert(arrRes)
        window.location.replace("login.html")
    }

    async function ExecuteCommand(content) {
        const envio = await fetch('http://3.14.134.83/individualComand',{
         method : "POST",
         body : JSON.stringify({"cmd":content}),
        })
        document.getElementById("graph").textContent = ""
        const response = await envio.json()
        console.log(response)
        let arrRes = response.Message
        let typeC = response.TypeC
        if (typeC === "rep") {
            console.log(arrRes)
            if (arrRes === "No Login") {
                alert("Debe estar logeado")    
            } else {
                d3.select("#graph")
                .graphviz()
                    .dot(arrRes)
                    .render();
            }
        } else {
            alert(arrRes)
        }
    }
    
})

