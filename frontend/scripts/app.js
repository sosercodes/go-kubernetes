// app.js

const getMessage = async (wordCount) => {
    const response = await fetch(`http://localhost:80/api/message`, {
        method: 'GET',
        mode: 'cors', // cors, no-cors, *cors, same-origin);
    });
    if (response.status === 200) {
        const data = await response.json();
        return data;
    } else {
        throw new Error('Unable to get message')
    }
}
getMessage().then((data) => {
    document.querySelector('#title').innerHTML = data.title;
    document.querySelector('#body').innerHTML = data.body;
}).catch((err) => {
    console.log(`Error: ${err}`);
})