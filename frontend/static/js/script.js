console.log("script.js is successfully connected!");
// This is a placeholder function. In a real app, you would implement the actual logic.
function chooseToGift(childId) {
    console.log("Gift chosen for child with ID:", childId);
    // Here you would likely make an AJAX call to your backend
}

// Add event listeners to gift buttons
document.addEventListener("DOMContentLoaded", function() {
    var giftButtons = document.querySelectorAll(".gift-button");
    giftButtons.forEach(function(button) {
        button.addEventListener("click", function() {
            var childId = this.closest('.child-wish').getAttribute('data-child-id');
            chooseToGift(childId);
        });
    });
});
// Example JavaScript for client-side interactivity
document.addEventListener("DOMContentLoaded", function() {
    var updateButton = document.querySelector("button[type='submit']");
    updateButton.addEventListener("click", function(event) {
        // You can add validation or other interactivity here
        console.log("Wishes updated!");
    });
});

document.addEventListener('DOMContentLoaded', function() {
    var form = document.querySelector('form');
    form.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent the default form submission

        var wishes = document.getElementById('wishes').value;
        var data = {
            wishes: wishes
        };

        var jsonData = JSON.stringify(data);

        console.log(jsonData); // For debugging

        // Here you can send jsonData to your server
        // Example using fetch API
        fetch('/update-wishes', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData
        })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            // Handle success here (e.g., showing a success message)
        })
        .catch((error) => {
            console.error('Error:', error);
            // Handle errors here
        });
    });
});

document.addEventListener('DOMContentLoaded', (event) => {
    document.getElementById('js-check').innerText = "script.js is successfully connected!";
});
