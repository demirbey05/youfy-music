let audioElement;

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
                return response.blob().then(createObjectURLAndPlay);
            } else {
                throw new Error('Invalid content type');
            }
        })
        .catch(error => {
            console.error('Error:', error);
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

// Optional: Stop playback when the user navigates away from the page
window.addEventListener('beforeunload', () => {
    if (audioElement) {
        audioElement.pause();
    }
});
