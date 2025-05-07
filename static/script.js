document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('searchInput');
    const submitButton = document.getElementById('submitButton');
    const statusArea = document.getElementById('statusArea');
    const resultsList = document.getElementById('resultsList');

    const loadingGifContainer = document.getElementById('loadingGifContainer');
    const loadingGifImage = document.getElementById('loadingGif');
    const gifUrls = [
        "https://c.tenor.com/aAzr2xbOK-EAAAAd/tenor.gif",
        "https://media.tenor.com/FF5STjI711cAAAAj/anthony-fantano-ant-cube.gif",
        "https://c.tenor.com/Y1cFztN5CZwAAAAd/tenor.gif",
        "https://i.giphy.com/JlpguxhuxwHiqKXoW7.webp",
        "https://c.tenor.com/Dd3ePv9oUeQAAAAd/tenor.gif",
        "https://c.tenor.com/cmNQuKadXqUAAAAd/tenor.gif",
        "https://c.tenor.com/v6VivIG_8BEAAAAd/tenor.gif",
        "https://c.tenor.com/hmQKFElzIMQAAAAd/tenor.gif",
        "https://c.tenor.com/S4c83WORXVEAAAAd/tenor.gif"
    ];

    function showLoadingGif() {
        if (gifUrls.length > 0) {
            const randomIndex = Math.floor(Math.random() * gifUrls.length);
            loadingGifImage.src = gifUrls[randomIndex];
            loadingGifContainer.style.display = 'block';
        }
    }

    function hideLoadingGif() {
        loadingGifContainer.style.display = 'none';
        loadingGifImage.src = ''; // Clear src to stop animation and loading
    }

    let currentJobId = null;
    let pollingInterval = null;

    submitButton.addEventListener('click', async () => {
        const query = searchInput.value.trim();
        if (!query) {
            updateStatus('Please enter a search query.', 'status-failed');
            return;
        }

        submitButton.disabled = true;
        updateStatus('Submitting your request...', 'status-processing');
        resultsList.innerHTML = ''; // Clear previous results
        hideLoadingGif(); // Ensure GIF is hidden at the start of a new submission

        try {
            const response = await fetch('/api/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query }),
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({ error: 'Failed to submit request. Server returned an error.' }));
                throw new Error(errorData.error || `Server error: ${response.status}`);
            }

            const data = await response.json();
            currentJobId = data.job_id;
            updateStatus(`Request submitted! Job ID: ${currentJobId}. Queued...`, 'status-queued');
            // GIF will be shown by handleJobStatus if status becomes 'processing'
            startPolling();

        } catch (error) {
            console.error('Submission error:', error);
            updateStatus(`Error: ${error.message}`, 'status-failed');
            submitButton.disabled = false;
            hideLoadingGif(); // Hide GIF on submission error
        }
    });

    function startPolling() {
        if (pollingInterval) {
            clearInterval(pollingInterval);
        }

        pollingInterval = setInterval(async () => {
            if (!currentJobId) {
                clearInterval(pollingInterval);
                return;
            }

            try {
                const response = await fetch(`/api/status/${currentJobId}`);
                if (!response.ok) {
                    if (response.status === 404) {
                         updateStatus(`Job ${currentJobId} not found. It might have expired or an error occurred.`, 'status-failed');
                         stopPollingAndEnableSubmit(); // This will also hide GIF
                         return;
                    }
                    throw new Error(`Server error: ${response.status}`);
                }

                const job = await response.json();
                handleJobStatus(job);

            } catch (error) {
                console.error('Polling error:', error);
                updateStatus(`Error polling status: ${error.message}`, 'status-failed');
                hideLoadingGif(); // Hide GIF on polling error
                // Consider stopping polling on persistent errors
                // stopPollingAndEnableSubmit(); 
            }
        }, 3000); // Poll every 3 seconds
    }

    function handleJobStatus(job) {
        switch (job.status) {
            case 'queued':
                hideLoadingGif(); // Not processing yet
                updateStatus(`Your request is queued... (Job ID: ${job.job_id})`, 'status-queued');
                resultsList.innerHTML = ''; // Keep results area clear
                break;
            case 'processing':
                updateStatus(`Processing your request for "${job.query}"...`, 'status-processing');
                resultsList.innerHTML = ''; // Clear results area to make space for GIF
                showLoadingGif();
                break;
            case 'completed':
                hideLoadingGif(); // Hide GIF before showing results
                updateStatus(`Successfully found vibes for "${job.query}"!`, 'status-completed');
                displayResults(job.result);
                stopPollingAndEnableSubmit();
                break;
            case 'failed':
                hideLoadingGif(); // Hide GIF on failure
                updateStatus(`Failed to process "${job.query}": ${job.error}`, 'status-failed');
                stopPollingAndEnableSubmit();
                break;
            default:
                hideLoadingGif(); // Hide GIF for any unknown status
                updateStatus(`Unknown status: ${job.status}`, 'status-failed');
        }
    }
    
    function stopPollingAndEnableSubmit() {
        if (pollingInterval) {
            clearInterval(pollingInterval);
            pollingInterval = null;
        }
        currentJobId = null;
        submitButton.disabled = false;
        hideLoadingGif(); // Ensure GIF is hidden when polling stops
    }

    function updateStatus(message, className) {
        statusArea.innerHTML = `<p class="${className || ''}">${message}</p>`;
    }

    function displayResults(tracks) {
        resultsList.innerHTML = ''; // Clear previous results or loading messages
        if (!tracks || tracks.length === 0) {
            const li = document.createElement('li');
            li.textContent = 'No tracks found matching your vibe. Try a different query!';
            resultsList.appendChild(li);
            return;
        }

        tracks.forEach(track => {
            const li = document.createElement('li');
            li.innerHTML = `
                <strong>${escapeHtml(track.title)}</strong>
                <span>${escapeHtml(track.artist)}</span><br>
                <a href="${escapeHtml(track.link)}" target="_blank">Listen on Spotify</a>
            `;
            resultsList.appendChild(li);
        });
    }

	function escapeHtml(unsafe) {
		if (unsafe === null || typeof unsafe === 'undefined') return '';
		return unsafe
			 .toString()
			 .replace(/&/g, "&amp;")
			 .replace(/</g, "&lt;")
			 .replace(/>/g, "&gt;")
			 .replace(/"/g, "&quot;")
			 .replace(/'/g, "&#039;");
	}

    searchInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            submitButton.click();
        }
    });
});