// ========== ГЛОБАЛЬНЫЕ ПЕРЕМЕННЫЕ =============================================================
let testData = null;
let currentQuestionIndex = 0;
let userAnswers = {};
let userData = null;

// ========== DOM ЭЛЕМЕНТЫ =================================================================================
const welcomeScreen = document.getElementById("welcomeScreen");
const welcomeHint = document.getElementById("welcomeHint");
const testScreen = document.getElementById("testScreen");
const resultPage = document.getElementById("resultPage");
const nextBtn = document.getElementById("nextBtn");
const backBtn = document.getElementById("backBtn");
const newTestBtn = document.getElementById("newTestBtn");
const questionText = document.getElementById("questionText");
const optionsContainer = document.getElementById("optionsContainer");
const progressIndicator = document.getElementById("progressIndicator");

// ========== ЗАГРУЗКА ТЕСТА ============================================================
async function loadTest() {
  try {
    const response = await fetch("http://localhost:8080/tests/1");

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    testData = await response.json();
    console.log("Тест загружен:", testData);

    userAnswers = {};
    testData.questions.forEach((q) => {
      userAnswers[q.id] = null;
    });

    currentQuestionIndex = 0;
    renderCurrentQuestion();
  } catch (error) {
    console.error("Ошибка загрузки теста:", error);
    alert(
      "Не удалось загрузить тест. Проверьте, запущен ли сервер на http://localhost:8080",
    );
  }
}

// ========== ОБРАБОТЧИК КЛИКА НА ПРИВЕТСТВЕННОМ ЭКРАНЕ ==================================
function handleWelcomeScreenClick() {
  // Запускаем тест (откроется окно регистрации)
  startTestFlow();
}

// Инициализация обработчиков приветственного экрана
function initWelcomeScreen() {
  welcomeScreen.addEventListener("click", handleWelcomeScreenClick);
  welcomeScreen.style.cursor = "pointer";
}

// ========== МОДАЛЬНОЕ ОКНО РЕГИСТРАЦИИ ==========================================================
function showRegistrationModal() {
  return new Promise((resolve) => {
    const modalDiv = document.createElement("div");
    modalDiv.className = "modal";
    modalDiv.innerHTML = `
            <div class="modal-content">
                <p class="registration-prompt">Пожалуйста, представься, зайчик</p>                
                <input type="text" id="regName" placeholder="Имя *" autocomplete="off">
                <input type="password" id="regPassword" placeholder="Пароль *" autocomplete="off">
                <input type="date" id="regBirth" placeholder="Дата рождения">
                <input type="email" id="regEmail" placeholder="E-mail">
                <div id="modalError" class="error-message"></div>
                <div class="modal-buttons">
                    <button class="btn-secondary" id="cancelReg">Отмена</button>
                    <button class="btn-primary" id="confirmReg">Продолжить</button>
                </div>
            </div>
        `;
    document.body.appendChild(modalDiv);

    const nameInput = modalDiv.querySelector("#regName");
    const passwordInput = modalDiv.querySelector("#regPassword");
    const birthInput = modalDiv.querySelector("#regBirth");
    const emailInput = modalDiv.querySelector("#regEmail");
    const errorDiv = modalDiv.querySelector("#modalError");
    const confirmBtn = modalDiv.querySelector("#confirmReg");
    const cancelBtn = modalDiv.querySelector("#cancelReg");

    const cleanup = () => modalDiv.remove();

    function escapeHtml(str) {
      const div = document.createElement("div");
      div.textContent = str;
      return div.innerHTML;
    }

    const validate = () => {
      const name = nameInput.value.trim();
      const password = passwordInput.value;
      const birth = birthInput.value;
      const email = emailInput.value.trim();

      if (!name) {
        errorDiv.innerText = "Зайчик забыл ввести имя";
        return false;
      }
      if (!password) {
        errorDiv.innerText = "Зайчик забыл ввести пароль";
        return false;
      }
      if (!birth) {
        errorDiv.innerText = "Зайчик забыл ввести дату рождения";
        return false;
      }
      if (!email) {
        errorDiv.innerText = "Зайчик забыл ввести email";
        return false;
      }
      return true;
    };

    // ─── Вспомогательные запросы к API ─────────────────
    async function checkFullMatch(name, password, birth, email) {
      try {
        const res = await fetch("http://localhost:8080/check-full-match", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name, password, birth, email }),
        });
        if (res.ok) return await res.json();
      } catch (e) {
        console.error("checkFullMatch error:", e);
      }
      return { found: false };
    }

    async function checkEmailExists(email) {
      try {
        const res = await fetch("http://localhost:8080/check-email-exists", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email }),
        });
        if (res.ok) return await res.json();
      } catch (e) {
        console.error("checkEmailExists error:", e);
      }
      return { found: false };
    }

    async function checkPasswordOwner(password, email) {
      try {
        const res = await fetch("http://localhost:8080/check-password-owner", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ password, email }),
        });
        if (res.ok) return await res.json();
      } catch (e) {
        console.error("checkPasswordOwner error:", e);
      }
      return { found: false };
    }

    async function checkEmailAndPassword(email, password) {
      try {
        const res = await fetch("http://localhost:8080/check-email-password", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password }),
        });
        if (res.ok) return await res.json();
      } catch (e) {
        console.error("checkEmailAndPassword error:", e);
      }
      return { found: false };
    }

    // ─── Главный обработчик ─────────────────────────
    confirmBtn.onclick = async () => {
      if (!validate()) return;

      const user = {
        name: nameInput.value.trim(),
        password: passwordInput.value,
        birth: birthInput.value || null,
        email: emailInput.value.trim(),
      };

      // 1. Проверка полного совпадения
      const fullMatch = await checkFullMatch(
        user.name,
        user.password,
        user.birth,
        user.email,
      );
      if (fullMatch.found && fullMatch.all_match) {
        userData = {
          id: fullMatch.user.id,
          name: fullMatch.user.name,
          password: user.password,
          birth: fullMatch.user.birth,
          email: fullMatch.user.email,
        };
        cleanup();
        resolve(true);
        return;
      }

      // 2. Проверка совпадения email
      const emailCheck = await checkEmailExists(user.email);
      if (emailCheck.found) {
        // Проверяем, совпадает ли пароль у этого email
        const emailPwdCheck = await checkEmailAndPassword(
          user.email,
          user.password,
        );
        if (emailPwdCheck.found) {
          // Email и пароль совпадают, но не все данные → "другие данные"
          errorDiv.innerText =
            " Странно, кажется в прошлый раз были другие данные...";
        } else {
          errorDiv.innerText = "Зайчик, такой логин уже используется";
        }
        return;
        if (owner.found) {
          // У кого-то другого такой пароль? Но email занят, а пароль совпадает с чьим-то?
          // Логичнее: если email занят, и пароль совпадает с паролем владельца email,
          // тогда это "другие данные". Иначе просто "логин занят".
          // Нужно узнать пароль владельца email. Давай сделаем ещё один эндпоинт или используем
          // тот факт, что fullMatch не сработал. Значит, при занятом email пароль не совпал,
          // потому что fullMatch проверяет email+password+... => нет полного совпадения.
          // Следовательно, это просто занятый email.
          errorDiv.innerText = "Зайчик, такой логин уже используется";
        } else {
          errorDiv.innerText = "Зайчик, такой логин уже используется";
        }
        return;
      }

      // 3. Проверка совпадения пароля
      const passwordCheck = await checkPasswordOwner(user.password, user.email);
      if (passwordCheck.found) {
        errorDiv.innerHTML = `Такой пароль уже используется пользователем <strong>${escapeHtml(passwordCheck.suggested_email)}</strong> с именем <strong>${escapeHtml(passwordCheck.suggested_name)}</strong>`;
        return;
      }

      // 4. Если ничего не занято – регистрируем
      try {
        const response = await fetch("http://localhost:8080/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(user),
        });

        if (response.ok) {
          const savedUser = await response.json();
          userData = {
            id: savedUser.id,
            name: user.name,
            password: user.password,
            birth: user.birth,
            email: user.email,
          };
          cleanup();
          resolve(true);
        } else {
          const error = await response.json();
          const errorText = error.message || error.error || "";
          errorDiv.innerText =
            errorText || "Зайчик, произошла ошибка регистрации";
        }
      } catch (error) {
        errorDiv.innerText = "❌ Ошибка соединения с сервером";
        console.error("Registration error:", error);
      }
    };

    cancelBtn.onclick = () => {
      cleanup();
      resolve(false);
    };
  });
}

/// ========== ОТПРАВКА РЕЗУЛЬТАТОВ ==================================================================
async function submitTest() {
  const allAnswered = testData.questions.every(
    (q) => userAnswers[q.id] !== null,
  );
  if (!allAnswered) {
    alert("Пожалуйста, ответьте на все вопросы");
    return false;
  }

  try {
    const response = await fetch("http://localhost:8080/tests/1/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        user_id: userData?.id,
        birth_date: userData?.birth,  // Дата рождения из регистрации
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    showResultPage(result);
    return true;
  } catch (error) {
    console.error("Ошибка отправки результатов:", error);
    alert("Ошибка при отправке результатов. Попробуйте ещё раз.");
    return false;
  }
}

// ========== ПОКАЗ РЕЗУЛЬТАТОВ ===================================================================
function showResultPage(data) {
  welcomeScreen.style.display = "none";
  testScreen.style.display = "none";
  resultPage.style.display = "block";

  const userInfoDiv = document.getElementById("resultUserInfo");
  if (userInfoDiv) {
    userInfoDiv.innerHTML = `
      <strong>👤 ${userData?.name || "Гость"}</strong><br>
      Дата рождения: ${userData?.birth || "Не указана"}<br>
      Знак зодиака: ${data.zodiac_sign || "Не определён"}<br>
      Email: ${userData?.email || "Не указан"}
    `;
  }

  const resultMessageDiv = document.getElementById("resultMessage");
  if (resultMessageDiv) {
    resultMessageDiv.innerHTML = `
      <strong>${data.result.title || "Результат"}</strong><br><br>
      ${data.result.description || "Спасибо за прохождение теста!"}
    `;
  }

  // Картинка результата (если есть)
  const resultImage = document.getElementById("resultImage");
  if (resultImage) {
    resultImage.src = `/static/images/results/${data.result.id}.png`;
    resultImage.style.display = "block";
  }

}
// ========== НАВИГАЦИЯ ===========================================================================
function goToNext() {
  if (!testData) return;

  const currentQ = testData.questions[currentQuestionIndex];

  if (userAnswers[currentQ.id] === null) {
    alert("Пожалуйста, выберите вариант ответа");
    return;
  }

  if (currentQuestionIndex + 1 < testData.questions.length) {
    currentQuestionIndex++;
    renderCurrentQuestion();
  } else {
    submitTest();
  }
}

function goToPrev() {
  if (currentQuestionIndex > 0) {
    currentQuestionIndex--;
    renderCurrentQuestion();
  }
}

// ========== СБРОС ===========================================================================
function resetAndStartOver() {
  currentQuestionIndex = 0;
  userAnswers = {};
  userData = null;

  if (testData) {
    testData.questions.forEach((q) => {
      userAnswers[q.id] = null;
    });
  }

  resultPage.style.display = "none";
  testScreen.style.display = "block";

  // Анимация, чтобы вопросики красиво выплывали
  testScreen.style.transition = "all 0.5s ease-out";
  testScreen.style.opacity = "0";
  testScreen.style.transform = "scale(0.95)";

  setTimeout(() => {
    testScreen.style.opacity = "1";
    testScreen.style.transform = "scale(1)";
  }, 50);

  // Перерисовываем первый вопрос
  renderCurrentQuestion();

  // Очищаем статистику
  const statsContainer = document.getElementById("statisticsContainer");
  if (statsContainer) {
    statsContainer.innerHTML = "";
  }
}

// ========== ЗАПУСК ТЕСТА =============================================================================
async function startTestFlow() {
  const registered = await showRegistrationModal();
  if (!registered) {
    return;
  }

  // Плавное переключение
  welcomeScreen.style.opacity = "0";
  welcomeScreen.style.transform = "scale(0.95)";

  setTimeout(() => {
    welcomeScreen.style.display = "none";
    testScreen.style.display = "block";

    testScreen.style.opacity = "0";
    testScreen.style.transform = "scale(0.95)";
    testScreen.style.transition = "all 0.5s ease-out";

    setTimeout(() => {
      testScreen.style.opacity = "1";
      testScreen.style.transform = "scale(1)";
    }, 50);
  }, 300);

  await loadTest();
}

// ========== ОБРАБОТЧИКИ ==================================================================
nextBtn.onclick = goToNext;
backBtn.onclick = goToPrev;

if (newTestBtn) {
  newTestBtn.onclick = resetAndStartOver;
}

// ========== ИНИЦИАЛИЗАЦИЯ ==================================
welcomeScreen.style.display = "block";
testScreen.style.display = "none";
resultPage.style.display = "none";

initWelcomeScreen();
