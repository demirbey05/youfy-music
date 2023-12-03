let audioContext = new (window.AudioContext || window.webkitAudioContext)();
let audioBuffer;
let sourceNode;

function postData() {
    const inputData = document.getElementById('inputData').value;

    fetch('/submit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ data: inputData }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const contentType = response.headers.get('content-type');

        if (contentType && contentType.includes('audio/mp3')) {
            return response.body.getReader().read().then(processAudioStream);
        } else {
            throw new Error('Invalid content type');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function processAudioStream({ value, done }) {
    new Audio("data:audio/mp3;base64," + btoa(String.fromCharCode.apply(null, value))).play()
}

function playAudio() {
    // Trigger the postData function to fetch and play the audio
    postData();
}