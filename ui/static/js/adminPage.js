const form = document.getElementById('mailForm')

form.addEventListener('submit', (e) => {
    e.preventDefault()
    const {text} = Object.fromEntries(new FormData(form).entries())
    const options = {
        method: 'POST',
        body: JSON.stringify({text})
    }
    fetch('http://localhost:8080/send-mail', options)
})