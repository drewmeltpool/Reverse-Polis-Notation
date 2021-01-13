const url = "127.0.0.1:3002"

var inputForm = document.getElementById("inputForm")

inputForm.addEventListener("submit", (e)=>{

    //prevent auto submission
    e.preventDefault()

    const formdata = new FormData(inputForm)
    fetch(url,{
        method:"POST",
        body:formdata,
    }).then(
        response => response.text()
    ).then(
        (data) => {
            const result = document.getElementById("serverMessageBox")
            result.innerHTML = ""
            //result.textContent = data
            const jsonList = JSON.parse(data).Items

            let list = document.createElement("ul")
            list.className = "list"
            result.appendChild(list)
            
            for (let i = 0; i < jsonList.length; i++) {
                let item = document.createElement("li");
                item.className = "list__item"

                for(let name in jsonList[i]){
                    let type = document.createElement("p");
                    let value = document.createElement("span");
                    type.className = "text"
                    value.className = "value"
                    value.innerHTML = jsonList[i][name]
                    type.innerHTML = name + ": "
                    type.appendChild(value)
                    item.appendChild(type)
                }

                list.append(item)
            }

            result.appendChild(list)
        }
    ).catch(
        error => console.error(error)
    )




})