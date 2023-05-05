document.addEventListener("DOMContentLoaded",()=>{
    const btnEnviar = document.getElementById('enviarBtn')
    btnEnviar.onclick = () => {
        login()
    }
    async function login(){
        const username = document.getElementById('userTxt').value
        const password = document.getElementById('pwdTxt').value
        const idparticion = document.getElementById('idTxt').value
        let content = {username,password,idparticion}
        console.log(username)
        if (username != "" && password != "" && idparticion != "") {
            const envio = await fetch('http://3.14.134.83/login',{
            method : "POST",
            body : JSON.stringify(content),
            })
            const response = await envio.json()
            let arrRes = response.Message
            console.log(arrRes)
            if (arrRes === "EL") {
                alert("Error al realizar el login")
            } else if (arrRes === "ui") {
                alert("Usuario Inexistente")
            } else if (arrRes === "pi"){
                alert("Password incorrecta")
            }else if (arrRes === "SA"){
                alert("Sesion Activa")
            }else if (arrRes === "ok"){
                window.location.replace("main.html")
            }
        } else {
            alert("Todos los espacios son obligatorios")
        }
    }
})