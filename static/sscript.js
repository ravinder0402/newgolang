function searchBooks() {
    var title = document.getElementById('title').value;
    var author = document.getElementById('author').value;
    var publisher = document.getElementById('publisher').value;
 
    var requestData = {
        title: title,
        author: author,
        publisher: publisher
      
    };
 
    fetch('/search/book', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            var resultsContainer = document.getElementById('searchResults');
            resultsContainer.innerHTML = 'ISBN: '+ data.isbn+' LibID: '+ data.lib_id+ ' Title: ' + data.title + ', Author: ' + data.authors+ ', Publisher: ' + data.publisher+' Version: '
            + data.version+' TotalCopies: '+ data.total_copies+' AvailableCopies: '+ data.available_copies ;
            //window.location.href = './searchbook.html';
        }

       

    })
    .catch((error) => {
        console.error('Error:', error);
    });
}
 
