let audioElement;
let loadingCircle = document.getElementById('loading-circle');
let playButton = document.getElementById('playButton');
let stopButton = document.getElementById('stopButton');

function showLoading() {
    loadingCircle.style.display = 'block';
}

function hideLoading() {
    loadingCircle.style.display = 'none';
}

function enableButtons() {
    playButton.disabled = false;
    stopButton.disabled = false;
}

function disableButtons() {
    playButton.disabled = true;
    stopButton.disabled = true;
}

function postData() {
    const inputData = document.getElementById('inputData').value;

    // Show loading indicator while waiting for response
    showLoading();

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
                return response.blob().then(createObjectURLAndPlay);
            } else {
                throw new Error('Invalid content type');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        })
        .finally(() => {
            // Hide loading indicator and enable playback controls
            hideLoading();
            enableButtons();
        });
}

function createObjectURLAndPlay(blob) {
    const objectURL = URL.createObjectURL(blob);

    if (audioElement) {
        audioElement.pause();
        audioElement.src = objectURL;
    } else {
        audioElement = new Audio(objectURL);
    }

    audioElement.play();
}

function playAudio() {
    if (audioElement) {
        audioElement.play();
    }
}

function stopAudio() {
    if (audioElement) {
        audioElement.pause();
    }
}
