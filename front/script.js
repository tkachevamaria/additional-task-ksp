// ========== ДАННЫЕ ТЕСТА (3 вопроса, по 3 варианта ответа) ==========
const QUESTIONS = [
    {
        text: "Как вы предпочитаете проводить свободное время?",
        options: ["Читать книгу или смотреть фильм", "Встречаться с друзьями", "Заниматься спортом или хобби"]
    },
    {
        text: "Какой стиль работы вам ближе?",
        options: ["Планомерно и по правилам", "Креативно и свободно", "В команде, с поддержкой"]
    },
    {
        text: "Что для вас важнее в жизни?",
        options: ["Стабильность и покой", "Новые впечатления", "Гармония с собой и миром"]
    }
];

// ========== ГЛОБАЛЬНЫЕ ПЕРЕМЕННЫЕ ==========
let currentQuestionIndex = 0;
let userAnswers = new Array(QUESTIONS.length).fill(null);
let userData = null; // { name, birth, email }

// DOM элементы
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

// ========== МОДАЛЬНОЕ ОКНО РЕГИСТРАЦИИ ==========
function showRegistrationModal() {
    return new Promise((resolve) => {
        const modalDiv = document.createElement('div');
        modalDiv.className = 'modal';
        modalDiv.innerHTML = `
            <div class="modal-content">
                <h2>📝 Регистрация</h2>
                <p style="color: #112250; margin-bottom: 10px;">Пожалуйста, представьтесь</p>
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

            return true;
        };

        confirmBtn.onclick = () => {
            if (validate()) {
                userData = {
                    name: nameInput.value.trim(),
                    password: passwordInput.value,
                    birth: birthInput.value || 'Не указана',
                    email: emailInput.value.trim() || 'Не указан'
                };
                cleanup();
                resolve(true);
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
    const q = QUESTIONS[currentQuestionIndex];
    questionText.innerText = q.text;
    progressIndicator.innerText = `Вопрос ${currentQuestionIndex + 1} / ${QUESTIONS.length}`;
    
    optionsContainer.innerHTML = '';
    q.options.forEach((opt, idx) => {
        const isSelected = (userAnswers[currentQuestionIndex] === idx);
        const div = document.createElement('div');
        div.className = `option-item ${isSelected ? 'selected' : ''}`;
        
        const radio = document.createElement('input');
        radio.type = 'radio';
        radio.name = 'question';
        radio.value = idx;
        radio.checked = isSelected;
        radio.className = 'option-radio';
        
        const label = document.createElement('label');
        label.className = 'option-label';
        label.innerText = opt;
        
        div.appendChild(radio);
        div.appendChild(label);
        
        div.addEventListener('click', (e) => {
            if (e.target.tagName !== 'INPUT') {
                radio.checked = true;
            }
            userAnswers[currentQuestionIndex] = idx;
            document.querySelectorAll('.option-item').forEach(item => item.classList.remove('selected'));
            div.classList.add('selected');
        });
        
        radio.addEventListener('change', () => {
            userAnswers[currentQuestionIndex] = idx;
            document.querySelectorAll('.option-item').forEach(item => item.classList.remove('selected'));
            div.classList.add('selected');
        });
        
        optionsContainer.appendChild(div);
    });
    
    backBtn.disabled = (currentQuestionIndex === 0);
}

// ========== РАСЧЁТ СТАТИСТИКИ И РЕЗУЛЬТАТА ==========
function calculateResults() {
    // Подсчёт, какой вариант ответа чаще выбирался
    const answerCounts = [0, 0, 0];
    for (let i = 0; i < userAnswers.length; i++) {
        if (userAnswers[i] !== null) {
            answerCounts[userAnswers[i]]++;
        }
    }
    
    // Определяем доминирующий тип
    let dominantType = '';
    let maxCount = Math.max(...answerCounts);
    let dominantIndex = answerCounts.indexOf(maxCount);
    
    if (dominantIndex === 0) dominantType = '🏛️ Рациональный тип';
    else if (dominantIndex === 1) dominantType = '🎨 Креативный тип';
    else dominantType = '🌿 Гармоничный тип';
    
    // Текстовое описание результата
    let resultDescription = '';
    if (dominantIndex === 0) {
        resultDescription = 'Вы цените порядок, стабильность и логику. Вам нравится всё планировать и следовать правилам.';
    } else if (dominantIndex === 1) {
        resultDescription = 'Вы творческая личность, любите свободу и новые впечатления. Вам важно самовыражение.';
    } else {
        resultDescription = 'Вы стремитесь к балансу, цените гармонию и хорошие отношения с окружающими.';
    }
    
    return {
        dominantType: dominantType,
        description: resultDescription,
        answerCounts: answerCounts,
        totalQuestions: QUESTIONS.length
    };
}

// ========== ПОКАЗ СТРАНИЦЫ РЕЗУЛЬТАТОВ ==========
function showResultPage() {
    const results = calculateResults();
    
    // Скрываем карточку теста, показываем страницу результатов
    document.getElementById('app').style.display = 'none';
    resultPage.style.display = 'block';
    
    // Отображаем информацию о пользователе
    const userInfoDiv = document.getElementById('resultUserInfo');
    userInfoDiv.innerHTML = `
        <strong>👤 ${userData.name}</strong><br>
        📅 Дата рождения: ${userData.birth}<br>
        📧 Email: ${userData.email}
    `;
    
    // Отображаем основной результат
    const resultMessageDiv = document.getElementById('resultMessage');
    resultMessageDiv.innerHTML = `
        <strong>${results.dominantType}</strong><br><br>
        ${results.description}
    `;
    
    // Строим статистику по каждому вопросу
    const statsContainer = document.getElementById('statisticsContainer');
    statsContainer.innerHTML = '<h3 style="margin-bottom: 15px; color: #112250;">📈 Детальная статистика</h3>';
    
    for (let i = 0; i < QUESTIONS.length; i++) {
        const answerIndex = userAnswers[i];
        const question = QUESTIONS[i];
        const answerText = answerIndex !== null ? question.options[answerIndex] : 'Не отвечено';
        
        const statDiv = document.createElement('div');
        statDiv.className = 'stat-item';
        statDiv.innerHTML = `
            <div class="stat-question">Вопрос ${i + 1}: ${question.text}</div>
            <div class="stat-answer">✓ Ваш ответ: ${answerText}</div>
        `;
        statsContainer.appendChild(statDiv);
    }
    
    // Добавляем общую статистику по типам
    const summaryDiv = document.createElement('div');
    summaryDiv.className = 'stat-item';
    summaryDiv.style.marginTop = '20px';
    summaryDiv.style.background = '#E0C58F40';
    summaryDiv.innerHTML = `
        <div class="stat-question">📊 Общая статистика ответов</div>
        <div class="stat-answer">🏛️ Рациональные ответы: ${results.answerCounts[0]}</div>
        <div class="stat-answer">🎨 Креативные ответы: ${results.answerCounts[1]}</div>
        <div class="stat-answer">🌿 Гармоничные ответы: ${results.answerCounts[2]}</div>
    `;
    statsContainer.appendChild(summaryDiv);
    
    // Здесь можно отправить данные на сервер (для будущего бэка)
    console.log('Данные для отправки на сервер:', {
        user: userData,
        answers: userAnswers,
        results: results
    });
}

// ========== ПЕРЕХОД К СЛЕДУЮЩЕМУ ВОПРОСУ ==========
function goToNext() {
    if (userAnswers[currentQuestionIndex] === null) {
        alert('Пожалуйста, выберите вариант ответа');
        return;
    }
    
    if (currentQuestionIndex + 1 < QUESTIONS.length) {
        currentQuestionIndex++;
        renderCurrentQuestion();
    } else {
        // Последний вопрос - показываем результаты
        showResultPage();
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
    userAnswers = new Array(QUESTIONS.length).fill(null);
    userData = null;
    
    // Прячем страницу результатов и карточку теста
    resultPage.style.display = 'none';
    document.getElementById('app').style.display = 'block';
    
    // Показываем экран приветствия, прячем тест
    welcomeScreen.classList.remove('hidden');
    testScreen.classList.add('hidden');
}

// ========== ЗАПУСК ТЕСТА ПОСЛЕ РЕГИСТРАЦИИ ==========
async function startTestFlow() {
    const registered = await showRegistrationModal();
    if (!registered) {
        return;
    }
    
    // Прячем приветствие, показываем тест
    welcomeScreen.classList.add('hidden');
    testScreen.classList.remove('hidden');
    
    // Сбрасываем ответы на случай повторного прохождения
    currentQuestionIndex = 0;
    userAnswers = new Array(QUESTIONS.length).fill(null);
    renderCurrentQuestion();
}

// ========== НАВЕШИВАЕМ ОБРАБОТЧИКИ ==========
startBtn.onclick = startTestFlow;
nextBtn.onclick = goToNext;
backBtn.onclick = goToPrev;

if (newTestBtn) {
    newTestBtn.onclick = resetAndStartOver;
}

// Начальное состояние
welcomeScreen.classList.remove('hidden');
testScreen.classList.add('hidden');
resultPage.style.display = 'none';