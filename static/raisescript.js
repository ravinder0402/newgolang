
document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('issueForm');
    const message = document.getElementById('message');
   
    form.addEventListener('submit', function (event) {
      event.preventDefault();
      
      const formData = new FormData(form);
   
      fetch('/issue/request', {
        method: 'POST',
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
          message.innerHTML = `<p>issue has been registered successfully</p>`;
          form.reset();
        }
      })
      .catch(error => {
        console.error('Error:', error);
        message.innerHTML = '<p id="error">An error occurred. Please try again later.</p>';
      });
    });
  });
