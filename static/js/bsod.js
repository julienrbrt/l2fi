// bsod.js
// Provides a triggerBSOD function for consistent error handling via POST

function triggerBSOD(errorType, errorMessage) {
  const form = document.createElement("form");
  form.method = "POST";
  form.action = "/bsod";
  form.style.display = "none";

  const typeInput = document.createElement("input");
  typeInput.name = "type";
  typeInput.value = errorType;
  form.appendChild(typeInput);

  const messageInput = document.createElement("input");
  messageInput.name = "message";
  messageInput.value = errorMessage;
  form.appendChild(messageInput);

  document.body.appendChild(form);
  form.submit();
}
