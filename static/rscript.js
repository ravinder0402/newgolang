document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('removeBookForm');
    const message = document.getElementById('message');
   
    form.addEventListener('submit', function (event) {
      event.preventDefault();
      
      const formData = new FormData(form);
   
      fetch('/remove-book', {
        method: 'DELETE',
        body: JSON.stringify(Object.fromEntries(formData)),
        headers: {
          'Content-Type': 'application/json'
        }
      })
      .then(response => response.json())
      .then(data => {
        if (data.error) {
          message.innerHTML = `<p id="error">${data.error}</p>`;
        } else {
          message.innerHTML = `<p>${data.message}</p>`;
          form.reset();
        }
      })
      .catch(error => {
        console.error('Error:', error);
        message.innerHTML = '<p id="error">An error occurred. Please try again later.</p>';
      });
    });
  });