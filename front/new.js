let testData = null;
let currentQuestion = 0;
let userAnswers = {}; // question_id -> answer_id

async function loadTest() {
  const response = await fetch("http://localhost:8080/tests/1");
  testData = await response.json();

  console.log("Тест:", testData);

  showQuestion();
}

document.getElementById("startBtn").addEventListener("click", () => {
  document.getElementById("welcomeScreen").classList.add("hidden");
  document.getElementById("testScreen").classList.remove("hidden");

  loadTest(); // вот тут подключение к бэку
});

function showQuestion() {
  const question = testData.questions[currentQuestion];

  document.getElementById("questionText").textContent = question.text;

  const container = document.getElementById("optionsContainer");
  container.innerHTML = "";

  question.answers.forEach(answer => {
    const btn = document.createElement("button");
    btn.textContent = answer.text;

    btn.onclick = () => {
      userAnswers[question.id] = answer.id;
    };

    container.appendChild(btn);
  });

  document.getElementById("progressIndicator").textContent =
    `Вопрос ${currentQuestion + 1} / ${testData.questions.length}`;
}

document.getElementById("nextBtn").addEventListener("click", () => {
  currentQuestion++;

  if (currentQuestion >= testData.questions.length) {
    submitTest(); // отправка на сервер
    return;
  }

  showQuestion();
});

async function submitTest() {
  const response = await fetch("http://localhost:8080/tests/1/submit", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      answers: userAnswers
    })
  });

  const result = await response.json();

  showResult(result);
}

function showResult(result) {
  document.getElementById("app").style.display = "none";
  document.getElementById("resultPage").style.display = "block";

  document.getElementById("resultMessage").textContent =
    result.title + " — " + result.description;
}

