// ========== ГЛОБАЛЬНЫЕ ПЕРЕМЕННЫЕ ==========
let testData = null;           // Данные теста с бэкенда
let currentQuestionIndex = 0;
let userAnswers = {};          // question_id -> answer_index (для отправки на бэк)
let userData = null;           // { id, name, password, birth, email }

// ========== DOM ЭЛЕМЕНТЫ ==========
const welcomeScreen = document.getElementById('welcomeScreen');
const testScreen = document.getElementById('testScreen');
const resultPage = document.getElementById('resultPage');
const startBtn = document.getElementById('startBtn');
const nextBtn = document.getElementById('nextBtn');
const backBtn = document.getElementById('backBtn');
const newTestBtn = document.getElementById('newTestBtn');
const questionText = document.getElementById('questionText');
const optionsContainer = document.getElementById('optionsContainer');
const progressIndicator = document.getElementById('progressIndicator');

// ========== ЗАГРУЗКА ТЕСТА С БЭКЕНДА ==========
async function loadTest() {
    try {
       // Загружаем тест (ID = 1, потом можно будет выбирать)
        const response = await fetch("http://localhost:8080/tests/1");
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        testData = await response.json();
        console.log("Тест загружен:", testData);
        
        // Инициализируем userAnswers для всех вопросов
        userAnswers = {};
        testData.questions.forEach(q => {
            userAnswers[q.id] = null;
        });
        
        // Показываем первый вопрос
        currentQuestionIndex = 0;
        renderCurrentQuestion();
        
    } catch (error) {
        console.error("Ошибка загрузки теста:", error);
        alert("Не удалось загрузить тест. Проверьте, запущен ли сервер на http://localhost:8080");
    }
}

// ========== МОДАЛЬНОЕ ОКНО РЕГИСТРАЦИИ ==========
function showRegistrationModal() {
    return new Promise((resolve) => {
        const modalDiv = document.createElement('div');
        modalDiv.className = 'modal';
        modalDiv.innerHTML = `
            <div class="modal-content">
                <h2>📝 Регистрация</h2>
                <p style="color: #D9CBC2; margin-bottom: 10px;">Пожалуйста, представьтесь</p>
                <input type="text" id="regName" placeholder="Имя *" autocomplete="off">
                <input type="password" id="regPassword" placeholder="Пароль *" autocomplete="off">
                <input type="date" id="regBirth" placeholder="Дата рождения">
                <input type="email" id="regEmail" placeholder="E-mail (необязательно)">
                <div id="modalError" class="error-message"></div>
                <div class="modal-buttons">
                    <button class="btn-secondary" id="cancelReg">Отмена</button>
                    <button class="btn-primary" id="confirmReg">Продолжить</button>
                </div>
            </div>
        `;
        document.body.appendChild(modalDiv);

        const nameInput = modalDiv.querySelector('#regName');
        const passwordInput = modalDiv.querySelector('#regPassword');
        const birthInput = modalDiv.querySelector('#regBirth');
        const emailInput = modalDiv.querySelector('#regEmail');
        const errorDiv = modalDiv.querySelector('#modalError');
        const confirmBtn = modalDiv.querySelector('#confirmReg');
        const cancelBtn = modalDiv.querySelector('#cancelReg');

        const cleanup = () => modalDiv.remove();

        const validate = () => {
            const name = nameInput.value.trim();
            const password = passwordInput.value;
            
            if (!name) {
                errorDiv.innerText = 'Пожалуйста, введите имя';
                return false;
            }

            if (!password) {
                errorDiv.innerText = 'Пожалуйста, введите пароль';
                return false;
            }

            if (password.length < 6) {
                errorDiv.innerText = 'Пароль должен быть не менее 6 символов';
                return false;
            }

            return true;
        };

        confirmBtn.onclick = async () => {
            if (validate()) {
                const user = {
                    name: nameInput.value.trim(),
                    password: passwordInput.value,
                    birth: birthInput.value || null,
                    email: emailInput.value.trim() || null
                };
                
                // Отправляем данные пользователя на сервер в /register
                try {
                    const response = await fetch("http://localhost:8080/register", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(user)
                    });
                    
                    if (response.ok) {
                        const savedUser = await response.json();
                        userData = {
                            id: savedUser.id,
                            name: user.name,
                            password: user.password,
                            birth: user.birth,
                            email: user.email
                        };
                        cleanup();
                        resolve(true);
                    } else {
                        const error = await response.json();
                        errorDiv.innerText = error.message || 'Ошибка регистрации';
                    }
                } catch (error) {
                    errorDiv.innerText = 'Ошибка соединения с сервером';
                    console.error('Registration error:', error);
                }
            }
        };

        cancelBtn.onclick = () => {
            cleanup();
            resolve(false);
        };
    });
}

// ========== ОТРИСОВКА ТЕКУЩЕГО ВОПРОСА ==========
function renderCurrentQuestion() {
    if (!testData || !testData.questions) return;
    
    const question = testData.questions[currentQuestionIndex];
    questionText.innerText = question.text;
    progressIndicator.innerText = `Вопрос ${currentQuestionIndex + 1} / ${testData.questions.length}`;
    
    optionsContainer.innerHTML = '';
    
    question.answers.forEach((answer, idx) => {
        const isSelected = (userAnswers[question.id] === answer.id);
        const div = document.createElement('div');
        div.className = `option-item ${isSelected ? 'selected' : ''}`;
        
        const radio = document.createElement('input');
        radio.type = 'radio';
        radio.name = 'question_${question.id}';
        radio.value = answer.id;
        radio.checked = isSelected;
        radio.className = 'option-radio';
        radio.id = `q_${question.id}_${answer.id}`;
        
        const label = document.createElement('label');
        label.className = 'option-label';
        label.htmlFor = `q_${question.id}_${answer.id}`;
        label.innerText = answer.text;
        
        div.appendChild(radio);
        div.appendChild(label);
        
        div.addEventListener('click', (e) => {
            if (e.target.tagName !== 'INPUT') {
                radio.checked = true;
            }
            userAnswers[question.id] = answer.id;
            document.querySelectorAll('.option-item').forEach(item => item.classList.remove('selected'));
            div.classList.add('selected');
        });
        
        radio.addEventListener('change', () => {
            userAnswers[question.id] = answer.id;
            document.querySelectorAll('.option-item').forEach(item => item.classList.remove('selected'));
            div.classList.add('selected');
        });
        
        optionsContainer.appendChild(div);
    });
    
    backBtn.disabled = (currentQuestionIndex === 0);
}

// ========== ОТПРАВКА РЕЗУЛЬТАТОВ НА СЕРВЕР ==========
async function submitTest() {
    // Проверяем, что на все вопросы ответили
    const allAnswered = testData.questions.every(q => userAnswers[q.id] !== null);
    if (!allAnswered) {
        alert('Пожалуйста, ответьте на все вопросы');
        return false;
    }
    
    // Показываем экран загрузки (если он есть)
    const loadingScreen = document.getElementById('loadingScreen');
    if (loadingScreen) {
        document.getElementById('app').style.display = 'none';
        loadingScreen.style.display = 'block';
    }
    
    try {
        const response = await fetch("http://localhost:8080/tests/1/submit", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                user_id: userData?.id,
                answers: userAnswers
            })
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        
        if (loadingScreen) {
            loadingScreen.style.display = 'none';
        }
        
        showResultPage(result);
        return true;
        
    } catch (error) {
        console.error("Ошибка отправки результатов:", error);
        
        if (loadingScreen) {
            loadingScreen.style.display = 'none';
            document.getElementById('app').style.display = 'block';
        }
        
        alert("Ошибка при отправке результатов. Попробуйте ещё раз.");
        return false;
    }
}

// ========== ПОКАЗ СТРАНИЦЫ РЕЗУЛЬТАТОВ ==========
function showResultPage(result) {
    document.getElementById('app').style.display = 'none';
    resultPage.style.display = 'block';
    
    const userInfoDiv = document.getElementById('resultUserInfo');
    if (userInfoDiv) {
        userInfoDiv.innerHTML = `
            <strong>👤 ${userData?.name || 'Гость'}</strong><br>
            📅 Дата рождения: ${userData?.birth || 'Не указана'}<br>
            📧 Email: ${userData?.email || 'Не указан'}
        `;
    }
    
    const resultMessageDiv = document.getElementById('resultMessage');
    if (resultMessageDiv) {
        resultMessageDiv.innerHTML = `
            <strong>${result.title || 'Результат'}</strong><br><br>
            ${result.description || 'Спасибо за прохождение теста!'}
        `;
    }
    
    const statsContainer = document.getElementById('statisticsContainer');
    if (statsContainer && testData) {
        statsContainer.innerHTML = '<h3 style="margin-bottom: 15px; color: #E0C58F;">📈 Детальная статистика</h3>';
        
        for (let i = 0; i < testData.questions.length; i++) {
            const question = testData.questions[i];
            const answerId = userAnswers[question.id];
            const answer = question.answers.find(a => a.id === answerId);
            const answerText = answer ? answer.text : 'Не отвечено';
            
            const statDiv = document.createElement('div');
            statDiv.className = 'stat-item';
            statDiv.innerHTML = `
                <div class="stat-question">Вопрос ${i + 1}: ${question.text}</div>
                <div class="stat-answer">✓ Ваш ответ: ${answerText}</div>
            `;
            statsContainer.appendChild(statDiv);
        }
    }
    
    console.log('Результаты теста:', result);
}

// ========== ПЕРЕХОД К СЛЕДУЮЩЕМУ ВОПРОСУ ==========
function goToNext() {
    if (!testData) return;
    
    const currentQ = testData.questions[currentQuestionIndex];
    
    if (userAnswers[currentQ.id] === null) {
        alert('Пожалуйста, выберите вариант ответа');
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

// ========== ПОЛНЫЙ СБРОС И ЗАНОВО ==========
function resetAndStartOver() {
    currentQuestionIndex = 0;
    userAnswers = {};
    userData = null;
    
    if (testData) {
        testData.questions.forEach(q => {
            userAnswers[q.id] = null;
        });
    }
    
    resultPage.style.display = 'none';
    document.getElementById('app').style.display = 'block';
    
    welcomeScreen.classList.remove('hidden');
    testScreen.classList.add('hidden');
}

// ========== ЗАПУСК ТЕСТА ПОСЛЕ РЕГИСТРАЦИИ ==========
async function startTestFlow() {
    const registered = await showRegistrationModal();
    if (!registered) {
        return;
    }
    
    welcomeScreen.classList.add('hidden');
    testScreen.classList.remove('hidden');
    
    await loadTest();
}

// ========== НАВЕШИВАЕМ ОБРАБОТЧИКИ ==========
startBtn.onclick = startTestFlow;
nextBtn.onclick = goToNext;
backBtn.onclick = goToPrev;

if (newTestBtn) {
    newTestBtn.onclick = resetAndStartOver;
}

// ========== ИНИЦИАЛИЗАЦИЯ ==========
welcomeScreen.classList.remove('hidden');
testScreen.classList.add('hidden');
resultPage.style.display = 'none';