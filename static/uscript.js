
document.addEventListener('DOMContentLoaded', function () {
  const form = document.getElementById('updateBookForm');
  const message = document.getElementById('message');
 
  form.addEventListener('submit', function (event) {
    event.preventDefault();
    
    const formData = new FormData(form);
    const updatedDetails = {
      title: formData.get('title'),
      authors: formData.get('authors'),
      publisher: formData.get('publisher'),
      version: parseInt(formData.get('version'))
    };
 
    const requestData = {
      isbn: formData.get('isbn'),
      updated_details: updatedDetails
    };
 
    fetch('/update-book', {
      method: 'PATCH',
      body: JSON.stringify(requestData),
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        message.innerHTML = `<p id="error">${data.error}</p>`;
      } else {
        message.innerHTML = `<p>Book details updated successfully:</p><pre>${JSON.stringify(data, null, 2)}</pre>`;
        form.reset();
      }
    })
    .catch(error => {
      console.error('Error:', error);
      message.innerHTML = '<p id="error">An error occurred. Please try again later.</p>';
    });
  });
});