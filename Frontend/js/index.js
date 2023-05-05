document.addEventListener('DOMContentLoaded',()=>{
  const btn1 = document.getElementById('examinFile')
  const btn2 = document.getElementById('executeBtn1')
  btn1.onclick = () => {
    ReadFile()
  }
  btn2.onclick = () => {
    SendComands()
  }
  async function ReadFile(){
    let archivo = document.getElementById('executeFile').files[0]
    const reader = new FileReader()
    reader.addEventListener("load",(event) => {
      document.getElementById("floatingTextarea").textContent = ""
      document.getElementById('floatingTextarea').textContent = event.target.result
    })
    reader.readAsText(archivo,"UTF-8")
  }
  async function SendComands() {
    let area1 = document.getElementById('floatingTextarea').value
    if (area1 === "") {
      let archivo = document.getElementById('executeFile').files[0]
      const reader = new FileReader()
      reader.addEventListener("load",(event) => {
        document.getElementById('floatingTextarea').textContent = event.target.result
        area1 = event.target.result
        ExecuteCommand(area1)
      })
      reader.readAsText(archivo,"UTF-8")
    } else  {
      ExecuteCommand(area1)
    }
  }

  async function ExecuteCommand(content) {
     const envio = await fetch('http://3.14.134.83/execute',{
      method : "POST",
      body : JSON.stringify({"cmd":content}),
     })
    document.getElementById("responseContent").textContent = ""
    const response = await envio.json()
    let arrRes = response.Message
    let cont = ""
    arrRes.forEach(element => {
      cont += element + "\n"
    });
    document.getElementById("responseContent").textContent = cont
  }
})
