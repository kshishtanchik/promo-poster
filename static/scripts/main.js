function register(event) {
  const formElement = document.getElementById("registrationForm");
  const formData = new FormData(formElement);
  const registerData = Object.fromEntries(formData);
  console.log(registerData);
  fetch("event", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(registerData), 
  })
    .then((res) => res.json())
    .then((data) => console.log("Успешно:", data))
    .catch((error) => console.error("Ошибка:", error));
  debugger;
}

function solution(answer) {
  console.log(answer);
  fetch("event", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ answer: answer }), // Ваши данные для отправки
  })
    .then((res) => res.json())
    .then((data) => console.log("Успешно:", data))
    .catch((error) => console.error("Ошибка:", error));
}
