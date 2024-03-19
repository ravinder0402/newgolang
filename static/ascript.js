function addBook() {
    var isbn = document.getElementById('isbn').value;
    var lib_id = parseInt(document.getElementById('lib_id').value);
    var title = document.getElementById('title').value;
    var authors = document.getElementById('authors').value;
    var publisher = document.getElementById('publisher').value;
    var version = parseInt(document.getElementById('version').value);
    var total_copies = parseInt(document.getElementById('total_copies').value);
    var available_copies = parseInt(document.getElementById('available_copies').value);
    var adminEmail = document.getElementById('email').value;
 
    var requestData = {
        book: {
            ISBN: isbn,
            Lib_id: lib_id,
            Title: title,
            Authors: authors,
            Publisher: publisher,
            Version: version,
            TotalCopies: total_copies,
            AvailableCopies: available_copies
        },
        email: adminEmail

    };
 
    fetch('/add-book', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        
        body: JSON.stringify(requestData),
    })
    .then(response => response.json())
    .then(data => {
        var messageContainer = document.getElementById('message');
        if (data.error) {
            console.log(requestData);
            messageContainer.textContent = data.error;
        } else {
            messageContainer.textContent = 'Book added successfully!';
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}